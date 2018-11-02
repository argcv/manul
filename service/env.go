package service

import (
	"context"
	"fmt"
	pbUser "github.com/argcv/go-argcvapis/app/manul/user"
	"github.com/argcv/manul/client/mongo"
	"github.com/argcv/manul/client/workdir"
	"github.com/argcv/manul/config"
	"github.com/argcv/manul/model"
	"github.com/argcv/webeh/log"
	"github.com/pkg/errors"
)

type Env struct {
	UserService    *UserServiceImpl
	SecretService  *SecretServiceImpl
	ProjectService *ProjectServiceImpl
	JobService     *JobServiceImpl

	// Client
	mc *mongo.Client
	fs *workdir.Workdir
}

const (
	DbUserColl    = "user"
	DbSecretColl  = "secret"
	DbProjectColl = "project"
	DbJobColl     = "job"

	FsProjectDir = "project"
	FsJobDir     = "job"
)

const (
	SysAssertAdmin = "_root"
)

func (env *Env) SpawnMgoCli() *mongo.Client {
	return env.mc.Spawn()
}

func (env *Env) SpawnWorkdir() *workdir.Workdir {
	return env.fs.Spawn()
}

func (env *Env) SpawnProjectWorkdir() *workdir.Workdir {
	return env.SpawnWorkdir().Goto(FsProjectDir).Rebase()
}

func (env *Env) SpawnJobWorkdir() *workdir.Workdir {
	return env.SpawnWorkdir().Goto(FsJobDir).Rebase()
}

func (env *Env) ParseAuthInfo(ctx context.Context, auth *pbUser.AuthToken) (user *model.User, err error) {
	return env.UserService.ParseAuthInfo(ctx, auth)
}

func (env *Env) DatabaseSetup() (err error) {
	log.Infof("Check Database Environment...")

	admin := SysAssertAdmin

	uopt := model.ParseUserServiceFilterOption(fmt.Sprintf("n:%s", admin))
	if ul, err := env.UserService.listUsers(uopt, 0, 100); err != nil {
		log.Errorf("Find users failed in querying...")
		return err
	} else {
		if len(ul) == 0 {
			log.Infof("admin NOT found... creating users...")
			if user, secret, err := env.UserService.createUser(
				admin,
				"System Admin",
				"admin@example.org",
				model.UserType_ADMIN,
				admin); err != nil {
				log.Errorf("create user failed: %v", err)
				return err
			} else {
				config.SetClientUserName(user.Name)
				config.SetClientUserSecret(secret.Secret)
				log.Infof("Created user: %s, %sxxxx", user.Name, secret.Secret[:1])
			}

		} else if len(ul) > 1 {
			errMsg := fmt.Sprintf("found more than 1 admin (%v)?", len(ul))
			log.Error(errMsg)
			err = errors.New(errMsg)
			return err
		} else {
			// only 1 user
			log.Infof("Admin user %s was created, id: %v", ul[0].Name, ul[0].Id)
		}
	}

	return
}

func (env *Env) init() error {
	// user service
	env.UserService = &UserServiceImpl{
		env: env,
	}

	env.SecretService = &SecretServiceImpl{
		env: env,
	}

	// job
	env.JobService = &JobServiceImpl{
		env: env,
	}

	// project service
	env.ProjectService = NewProjectServiceImpl(env)

	// init client
	if mc, err := config.InitMongoClient(); err != nil {
		log.Fatalf("Init Mongo Failed: %v", err.Error())
		return err
	} else {
		env.mc = mc
	}

	log.Infof("Check MongoDB Connection....")
	{
		mc := env.SpawnMgoCli()
		defer mc.Close()
		if err := mc.Session.Ping(); err != nil {
			log.Errorf("Mongo Session initialize failed")
			return err
		} else {
			log.Infof("MongoDB Session is initialized")
		}
	}

	env.DatabaseSetup()

	// fs
	env.fs = workdir.NewWorkdir(config.GetFsWorkdir())
	log.Infof("Base dir: %v", env.fs.Base)

	return nil
}

func NewManulGlobalEnv() (env *Env, err error) {
	env = &Env{}
	if err = env.init(); err != nil {
		log.Fatalf("Server Init Failed!!! %s", err.Error())
	}
	return
}
