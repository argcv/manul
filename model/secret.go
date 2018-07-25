package model

import (
	pb "github.com/argcv/go-argcvapis/app/manul/secret"
	"github.com/argcv/manul/client/mongo"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Secret struct {
	Id          bson.ObjectId `bson:"_id,omitempty" json:"id"`
	UserId      bson.ObjectId `bson:"user_id,omitempty" json:"user_id"`
	Secret      string        `bson:"secret,omitempty" json:"secret"`
	TempToken   string        `bson:"temp_token,omitempty" json:"temp_token"`
	CreateTime  time.Time     `bson:"create_time,omitempty" json:"create_time,omitempty"`
	UpdatedTime time.Time     `bson:"updated_time,omitempty" json:"updated_time,omitempty"`
}

func (u *Secret) UpdateTime() *Secret {
	if u.CreateTime == time.Unix(0, 0) {
		u.CreateTime = time.Now()
	}
	u.UpdatedTime = time.Now()
	return u
}

func (u *Secret) ToPbSecret() (pbSecret *pb.Secret) {
	pbSecret = &pb.Secret{
		UserId: u.UserId.Hex(),
		Secret: u.Secret,
	}
	return
}

func FromPbSecret(u *pb.Secret) *Secret {
	return &Secret{
		UserId: mongo.SafeToObjectIdOrEmpty(u.UserId),
		Secret: u.Secret,
	}
}
