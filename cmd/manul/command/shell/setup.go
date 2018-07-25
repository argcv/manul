package shell

import (
	"fmt"
	"github.com/argcv/manul/config"
	"github.com/spf13/viper"
	"gopkg.in/abiosoft/ishell.v2"
)

//KeyDBMongoAddrs      = "db.mongo.addrs"
//KeyDBMongoDatabase   = "db.mongo.db"
//KeyDBMongoUser       = "db.mongo.user"
//KeyDBMongoPass       = "db.mongo.pass"
//KeyDBMongoTimeoutSec = "db.mongo.timeout_sec"

func (e *Env) setupUpdateMongoDB(c *ishell.Context) {
	{
		// addrs
		var addrs []string
		addMongoAddr := func() {
			rtAddr := e.GetStringReplNonEmpty(c, "Type an address", "(e.g. localhost:27017) :")
			addrs = append(addrs, rtAddr)
			c.Printf("Current Addresses List: %v\n", addrs)
		}
		doContinue := true
		for doContinue {
			addMongoAddr()
			doContinue = e.ConfirmRepl(c, "Add more address? ", false)
		}
		c.Printf("Mongo Addresses: %v\n", addrs)
		viper.Set(config.KeyDBMongoAddrs, addrs)
	}

	// default database
	db := e.GetStringReplNonEmpty(c, "Type database for auth", "(e.g. admin ) :")

	viper.Set(config.KeyDBMongoDatabase, db)

	{
		// auth
		if e.ConfirmRepl(c, "Does this mongo has authorization? ", true) {
			viper.Set(config.KeyDBMongoPerformAuth, true)

			// yes
			source := e.GetStringRepl(c, fmt.Sprintf("Type database for auth, '%s' if it is empty", db), "(e.g. admin ) :")
			if len(source) > 0 {
				viper.Set(config.KeyDBMongoAuthDatabase, source)
			} else {
				viper.Set(config.KeyDBMongoAuthDatabase, nil)
			}
			user := e.GetStringRepl(c, "Type username for auth, defaults to admin", "username :")
			if len(user) > 0 {
				viper.Set(config.KeyDBMongoAuthUser, user)
			} else {
				viper.Set(config.KeyDBMongoAuthUser, "admin")
			}

			pass := e.GetStringReplNonEmpty(c, "Type password for auth", "password :")
			viper.Set(config.KeyDBMongoAuthPass, pass)

			mech := e.GetStringRepl(c, "Type authorization mechanism, defaults to 'MONGODB-CR'", "mech :")

			if len(mech) > 0 {
				viper.Set(config.KeyDBMongoAuthMechanism, mech)
			} else {
				viper.Set(config.KeyDBMongoAuthMechanism, nil)
			}

		} else {
			viper.Set(config.KeyDBMongoPerformAuth, false)
		}

	}

}

func (e *Env) AddSetup() {
	cmd := &ishell.Cmd{
		Name: "setup",
		Help: "setup server environment (not for client side)",
		Func: func(c *ishell.Context) {
			if e.ConfirmRepl(c, "Update Mongo Credential ?", true) {
				e.setupUpdateMongoDB(c)
			} else {
				c.Println("Skipped..")
			}

			if e.ConfirmRepl(c, "Save Config ?", true) {
				viper.WriteConfig()
				c.Printf("Saved config file to %v\n", viper.ConfigFileUsed())
			} else {
				c.Println("I did nothing and quitting....")
			}

		},
	}

	e.Sh.AddCmd(cmd)
}
