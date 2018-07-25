package mongo

import (
	"github.com/argcv/webeh/log"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"math"
	"time"
)

var (
	MongoBatchSize = 2000
)

type Client struct {
	Session *mgo.Session
	Db      string
	isRoot  bool
}

/**
 *  Database is the default database name used when the Session.DB method
 *  is called with an empty name, and is also used during the initial
 *  authentication if Source is unset.
 *
 *  Username and Password inform the credentials for the initial authentication
 *  done on the database defined by the Source field. See Session.Login.
 *
 */
type Auth struct {
	Source    string
	Username  string
	Password  string
	Mechanism string
}

func (auth *Auth) UpdateDialInfo(info *mgo.DialInfo) {
	info.Source = auth.Source
	info.Username = auth.Username
	info.Password = auth.Password
	info.Mechanism = auth.Mechanism
}

/**
 *
 *  Addrs holds the addresses for the seed servers.
 *
 *  Database is the default database name used when the Session.DB method
 *  is called with an empty name, and is also used during the initial
 *  authentication if Source is unset.
 *
 *  Username and Password inform the credentials for the initial authentication
 *  done on the database defined by the Source field. See Session.Login.
 *
 *  Timeout is the amount of time to wait for a server to respond when
 *  first connecting and on follow up operations in the session. If
 *  timeout is zero, the call may block forever waiting for a connection
 *  to be established. Timeout does not affect logic in DialServer.
 */
func NewMongoClient(addrs []string, db string, timeout time.Duration, auth *Auth) (client *Client, err error) {
	info := &mgo.DialInfo{
		Addrs:    addrs,
		Timeout:  timeout,
		Database: db,
	}
	if auth != nil {
		auth.UpdateDialInfo(info)
	}
	var ses *mgo.Session
	ses, err = mgo.DialWithInfo(info)
	if err != nil {
		return nil, err
	}
	client = &Client{
		Session: ses,
		Db:      db,
		isRoot:  true,
	}
	return
}

func (client *Client) Spawn() (sub *Client) {
	if client.isRoot {
		sub = &Client{
			Session: client.Session.Copy(),
			Db:      client.Db,
			isRoot:  false,
		}
	} else {
		log.Output(log.ERROR, "Spawning from Non-Root", 2)
		sub = nil
	}
	return
}

func (client *Client) Close() {
	if client.isRoot {
		log.Output(log.ERROR, "Trying to **Close** root!!??", 2)
		client.Session.Close()
	} else {
		client.Session.Close()
	}
}

func (client *Client) UseDB(db string) *Client {
	client.Db = db
	return client
}

// Insert inserts one or more documents in the respective collection.  In
// case the session is in safe mode (see the SetSafe method) and an error
// happens while inserting the provided documents, the returned error will
// be of type *LastError.
func (client *Client) Insert(coll string, data interface{}) error {
	//data must be an address
	return client.Session.DB(client.Db).C(coll).Insert(data)
}

func (client *Client) Remove(coll string, query bson.M) error {
	_, err := client.Session.DB(client.Db).C(coll).RemoveAll(query)
	return err
}

func (client *Client) Count(coll string, query bson.M) (int, error) {
	return client.Session.DB(client.Db).C(coll).Find(query).Count()
}

func (client *Client) All(coll string, query bson.M, selector bson.M, results interface{}) error {
	//results must me an address
	return client.Session.DB(client.Db).C(coll).Find(query).Select(selector).All(results)
}

func (client *Client) AllSorted(coll string, query bson.M, sort []string, selector bson.M, results interface{}) error {
	//results must me an address
	return client.Session.DB(client.Db).C(coll).Find(query).Sort(sort...).Select(selector).All(results)
}

func (client *Client) Search(coll string, query bson.M, sort []string, offset int, size int, results interface{}) error {
	//results must me an address
	return client.Session.DB(client.Db).C(coll).Find(query).Sort(sort...).Skip(offset).Limit(size).All(results)
}

func (client *Client) One(coll string, query bson.M, result interface{}) error {
	//result must be an address
	return client.Session.DB(client.Db).C(coll).Find(query).One(result)
}

type docExistsTest struct {
	ID bson.ObjectId `json:"id" bson:"_id,omitempty"` // ID
}

func (client *Client) Exists(coll string, query bson.M) bool {
	r := &docExistsTest{}
	return client.One(coll, query, r) == nil
}

func (client *Client) Update(coll string, query bson.M, update interface{}) error {
	//update must be an address
	return client.Session.DB(client.Db).C(coll).Update(query, update)
}

// Upsert finds a single document matching the provided selector document
// and modifies it according to the update document.  If no document matching
// the selector is found, the update document is applied to the selector
// document and the result is inserted in the collection.
// If the session is in safe mode (see SetSafe) details of the executed
// operation are returned in info, or an error of type *LastError when
// some problem is detected.
//
// Relevant documentation:
//
//     http://www.mongodb.org/display/DOCS/Updating
//     http://www.mongodb.org/display/DOCS/Atomic+Operations
//
func (client *Client) Upsert(coll string, query bson.M, update interface{}) error {
	// update must be an address
	_, err := client.Session.DB(client.Db).C(coll).Upsert(query, update)
	return err
}

func (client *Client) BulkUpdate(coll string, pairs []interface{}) (err error) {
	c := client.Session.DB(client.Db).C(coll)

	batch := MongoBatchSize
	n := int(math.Ceil(float64(len(pairs)) / float64(batch)))
	for i := 0; i < n; i++ {
		start, end := i*batch, (i+1)*batch
		if end > len(pairs) {
			end = len(pairs)
		}
		log.Info(i, start, end, len(pairs), len(pairs[start:end]))

		bulk := c.Bulk()
		bulk.Unordered()
		bulk.Update(pairs[start:end]...)
		if _, err = bulk.Run(); err != nil {
			log.Info(err)
		}
	}

	return
}

func (client *Client) UpdateAll(coll string, query bson.M, update interface{}) (err error) {
	_, err = client.Session.DB(client.Db).C(coll).UpdateAll(query, update)
	return
}

func (client *Client) Iter(coll string, query bson.M, out chan bson.M) error {
	iter := client.Session.DB(client.Db).C(coll).Find(query).Iter()
	go func() {
		defer iter.Close()
		var result bson.M
		for iter.Next(&result) {
			out <- result
			result = bson.M{}
		}
		close(out)
	}()
	return nil

}

func (client *Client) NewMongoAutoCommiter(coll string) *MongoAutoCommitBuilder {
	return &MongoAutoCommitBuilder{
		client:        client,
		name:          coll,
		coll:          coll,
		numWorkers:    1,
		bulkActions:   1000,
		flushInterval: time.Duration(2) * time.Second,
		verbose:       true,
	}
}
