package model

import (
	pb "github.com/argcv/go-argcvapis/app/manul/user"
	"github.com/argcv/manul/client/mongo"
	"gopkg.in/mgo.v2/bson"
	"strings"
	"time"
)

type UserType pb.UserType

const (
	// could access everything
	// used to add/remove users
	// to/from projects
	UserType_ADMIN = pb.UserType_ADMIN
	// used to submit job
	// for ordinary students
	UserType_USER = pb.UserType_USER
	// special authorization
	// could read project
	// but can NOT submit job
	UserType_BOT = pb.UserType_BOT
)

type UserServiceFilterOption struct {
	UserType []pb.UserType
	Order    OrderType // 0: ignore, 1: asc, 2: desc
	Name     string    // user name
	Query    string
}

func ParseUserServiceFilterOption(filter string) (opt *UserServiceFilterOption) {
	opt = &UserServiceFilterOption{}
	tags := strings.Split(filter, "|")
	for _, tag := range tags {
		if strings.HasPrefix(tag, "t:") {
			// user type
			ut := tag[2:]
			if ut == "admin" {
				opt.UserType = append(opt.UserType, pb.UserType_ADMIN)
			} else if ut == "bot" {
				opt.UserType = append(opt.UserType, pb.UserType_BOT)
			} else if ut == "user" {
				opt.UserType = append(opt.UserType, pb.UserType_USER)
			}
		} else if strings.HasPrefix(tag, "o:") {
			// order
			ot := tag[2:]
			if ot == "asc" {
				opt.Order = OrderTypeAsc
			} else if ot == "desc" {
				opt.Order = OrderTypeDesc
			}
		} else if strings.HasPrefix(tag, "n:") {
			opt.Name = tag[2:]
		} else if strings.HasPrefix(tag, "q:") {
			opt.Query = tag[2:]
		}
	}
	return
}

type User struct {
	Id          bson.ObjectId `bson:"_id,omitempty" json:"id"`
	UserType    pb.UserType   `bson:"user_type,omitempty" json:"user_type,omitempty"`
	Name        string        `bson:"name,omitempty" json:"name,omitempty"`
	DisplayName string        `bson:"display_name,omitempty" json:"display_name,omitempty"`
	Email       string        `bson:"email,omitempty" json:"email,omitempty"`
	CreateTime  time.Time     `bson:"create_time,omitempty" json:"create_time,omitempty"`
	UpdatedTime time.Time     `bson:"updated_time,omitempty" json:"updated_time,omitempty"`
	CreatedBy   string        `bson:"created_by,omitempty" json:"created_by"`
}

func (u *User) UpdateTime() *User {
	if u.CreateTime == time.Unix(0, 0) {
		u.CreateTime = time.Now()
	}
	u.UpdatedTime = time.Now()
	return u
}

// from pb & to pb

func (u *User) ToPbUser() (pbUser *pb.User) {
	pbUser = &pb.User{
		Id:          u.Id.Hex(),
		UserType:    u.UserType,
		Name:        u.Name,
		DisplayName: u.DisplayName,
		Email:       u.Email,
	}
	return
}

func FromPbUser(u *pb.User) *User {
	return &User{
		Id:          mongo.SafeToObjectIdOrEmpty(u.Id),
		UserType:    u.UserType,
		Name:        u.Name,
		DisplayName: u.DisplayName,
		Email:       u.Email,
	}
}
