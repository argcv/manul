package service

import (
	"context"
	"fmt"
	pb "github.com/argcv/go-argcvapis/app/manul/secret"
	"github.com/argcv/manul/client/mongo"
	"github.com/argcv/manul/helpers"
	"github.com/argcv/manul/model"
	"github.com/argcv/webeh/log"
	"github.com/pkg/errors"
	"gopkg.in/mgo.v2/bson"
)

type SecretServiceImpl struct {
	env *Env
}

func (s *SecretServiceImpl) listSecret(uid string, offset, size int) (sl []model.Secret, err error) {
	mc := s.env.SpawnMgoCli()
	defer mc.Close()
	q := bson.M{}
	sort := []string{"_id"}
	if ouid, err := mongo.SafeToObjectId(uid); err == nil {
		q["user_id"] = ouid
	}
	err = mc.Search(DbSecretColl, q, sort, offset, size, &sl)
	return
}

func (s *SecretServiceImpl) findSecret(uid string) (secret *model.Secret, err error) {
	mc := s.env.SpawnMgoCli()
	defer mc.Close()
	q := bson.M{}
	if ouid, e := mongo.SafeToObjectId(uid); e == nil {
		q["user_id"] = ouid
	}
	log.Infof("query: %v", q)
	err = mc.One(DbSecretColl, q, &secret)
	return
}

func (s *SecretServiceImpl) verifySecret(uid, secret string) bool {
	if ret, err := s.findSecret(uid); err != nil {
		log.Errorf("find secret failed!!: %v", err)
		return false
	} else {
		return ret.Secret == secret
	}
}

func (s *SecretServiceImpl) updateSecret(uid string) (ret *model.Secret, err error) {
	mc := s.env.SpawnMgoCli()
	defer mc.Close()
	q := bson.M{}
	ret = &model.Secret{}
	if ouid, e := mongo.SafeToObjectId(uid); e == nil {
		q["user_id"] = ouid
		ret.UserId = ouid
	}
	ret.Secret = helpers.RandomString(32, helpers.CharsetHex)
	ret.UpdateTime()
	err = mc.Upsert(DbSecretColl, q, ret)
	return
}

func (s *SecretServiceImpl) tryUpdateSecret(uid, secret, temp_token string) (ret *model.Secret, err error) {
	if ret, err := s.findSecret(uid); err != nil {
		return nil, errors.New(fmt.Sprintf("secret not found"))
	} else if ret.Secret == secret || ret.TempToken == temp_token {
		return s.updateSecret(uid)
	} else {
		return nil, errors.New(fmt.Sprintf("invalid script"))
	}
}

/* TODO:Not implemented yet
 */
func (s *SecretServiceImpl) UpdateSecret(context.Context, *pb.UpdateSecretRequest) (*pb.UpdateSecretResponse, error) {
	return nil, errors.New("implement me")
}

/* TODO:Not implemented yet
 */
func (s *SecretServiceImpl) ForgotSecret(context.Context, *pb.ForgotSecretRequest) (*pb.ForgotSecretResponse, error) {
	return nil, errors.New("implement me")
}
