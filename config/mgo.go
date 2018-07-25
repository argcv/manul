package config

import (
	"github.com/argcv/manul/client/mongo"
	"time"
)

func GetDBMongoAddrs() []string {
	defaultAddrs := []string{"localhost:27017"}
	addrs := getStringSliceOrDefault(KeyDBMongoAddrs, defaultAddrs)
	if len(addrs) == 0 {
		addrs = defaultAddrs
	}
	return addrs
}

func GetDBMongoAuth() *mongo.Auth {
	if getBoolOrDefault(KeyDBMongoPerformAuth, false) {
		// with auth
		source := getStringOrDefault(KeyDBMongoAuthDatabase, "")
		user := getStringOrDefault(KeyDBMongoAuthUser, "admin")
		pass := getStringOrDefault(KeyDBMongoAuthPass, "")
		mech := getStringOrDefault(KeyDBMongoAuthMechanism, "")

		return &mongo.Auth{
			Source:    source,
			Username:  user,
			Password:  pass,
			Mechanism: mech,
		}
	}
	return nil
}

func GetDBMongoTimeout() time.Duration {
	dur := time.Duration(getInt64OrDefault(KeyDBMongoTimeoutSec, 0))
	return dur * time.Second
}

func InitMongoClient() (client *mongo.Client, err error) {
	addrs := GetDBMongoAddrs()
	auth := GetDBMongoAuth()
	db := getStringOrDefault(KeyDBMongoDatabase, "")
	if db == "" && auth != nil {
		db = auth.Source
	}
	timeout := GetDBMongoTimeout()
	return mongo.NewMongoClient(addrs, db, timeout, auth)
}
