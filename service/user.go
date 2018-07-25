package service

import (
	"fmt"
	pb "github.com/argcv/go-argcvapis/app/manul/user"
	"github.com/argcv/go-argcvapis/status/errcodes"
	"github.com/argcv/manul/client/mongo"
	"github.com/argcv/manul/model"
	"github.com/argcv/webeh/log"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"google.golang.org/grpc/metadata"
	"gopkg.in/mgo.v2/bson"
	"strings"
	"time"
)

type UserServiceImpl struct {
	env *Env
}

func (u *UserServiceImpl) findUserByName(name string) (user *model.User, err error) {
	mc := u.env.SpawnMgoCli()
	defer mc.Close()
	q := bson.M{
		"name": name,
	}
	err = mc.One(DbUserColl, q, &user)
	return
}

func (u *UserServiceImpl) findUserById(id string) (user *model.User, err error) {
	oid, err := mongo.SafeToObjectId(id)
	if err != nil {
		return nil, err
	}
	mc := u.env.SpawnMgoCli()
	defer mc.Close()
	q := bson.M{
		"_id": oid,
	}
	err = mc.One(DbUserColl, q, &user)
	return
}

func (u *UserServiceImpl) findUserByMixedId(q string) (user *model.User, err error) {
	if strings.HasPrefix(q, "$") && mongo.IsObjectIdHex(q[1:]) {
		return u.findUserById(q[1:])
	} else {
		return u.findUserByName(q)
	}
}

func (u *UserServiceImpl) findUserByIdOrName(id, name string) (user *model.User, err error) {
	if mongo.IsObjectIdHex(id) {
		return u.findUserById(id)
	} else {
		return u.findUserByName(name)
	}
}

func (u *UserServiceImpl) authUser(id, name, secret string) (user *model.User, err error) {
	user = &model.User{}
	if id != "" {
		if user, err = u.findUserById(id); err != nil {
			log.Infof("find failed: %v", err)
			return
		}
	} else {
		if user, err = u.findUserByName(name); err != nil {
			log.Infof("find failed: %v", err)
			return
		}
	}
	if verify := u.env.SecretService.verifySecret(user.Id.Hex(), secret); verify {
		return
	} else {
		return nil, errors.New("invalid_secret")
	}
}

func (u *UserServiceImpl) ParseAuthInfo(ctx context.Context, auth *pb.AuthToken) (user *model.User, err error) {
	var id string
	var name string
	var secret string

	// try extract from auth first
	if auth != nil {
		id = auth.Id
		name = auth.Name
		secret = auth.Secret
	}

	// try extract from metadata (header Authorization)
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if authl := md.Get("authorization"); len(authl) > 0 {
			auth := authl[0]
			log.Infof("auth: %v", auth)
			kv := strings.Split(auth, ":")
			if len(kv) == 3 {
				id = kv[0]
				name = kv[1]
				secret = kv[2]
			} else {
				log.Warnf("Invalid Authorization Code: %v, skip...", auth)
			}
		}
	} else {
		log.Warnf("Get Metadata FAILED!!!")
	}

	return u.authUser(id, name, secret)
}

// List users by options
func (u *UserServiceImpl) listUsers(option *model.UserServiceFilterOption, offset, size int) (ul []model.User, err error) {
	mc := u.env.SpawnMgoCli()
	defer mc.Close()
	q := bson.M{}
	var sort []string
	if option != nil {
		if len(option.UserType) > 0 {
			q["user_type"] = mongo.InQuery(option.UserType)
		}
		if len(option.Name) > 0 {
			q["name"] = option.Name
		}
		if option.Order == model.OrderTypeAsc {
			// asc
			sort = []string{"_id"}
		} else if option.Order == model.OrderTypeDesc {
			// desc
			sort = []string{"-_id"}
		}
	}
	err = mc.Search(DbUserColl, q, sort, offset, size, &ul)
	return
}

func (u *UserServiceImpl) countUsers(option *model.UserServiceFilterOption) (count int, err error) {
	mc := u.env.SpawnMgoCli()
	defer mc.Close()
	q := bson.M{}
	if option != nil {
		if len(option.UserType) > 0 {
			q["user_type"] = mongo.InQuery(option.UserType)
		}
		if len(option.Name) > 0 {
			q["name"] = option.Name
		}
	}
	return mc.Count(DbUserColl, q)
}

func (u *UserServiceImpl) ListUsers(ctx context.Context, req *pb.ListUsersRequest) (ret *pb.ListUsersResponse, e error) {
	if _, err := u.ParseAuthInfo(ctx, req.Auth); err != nil {
		st := model.Status{
			Code:    errcodes.Code_PERMISSION_DENIED,
			Message: fmt.Sprintf("invalid auth: %v", err.Error()),
		}
		ret = &pb.ListUsersResponse{
			Success: false,
			Result: &pb.ListUsersResponse_Error{
				Error: st.ToPbStatus(),
			},
		}
		return

	} else {
		offset := req.Offset
		size := req.Size
		if offset < 0 {
			offset = 0
		}
		if size == 0 {
			size = 10
		}
		opt := model.ParseUserServiceFilterOption(req.Filter)
		ul, e1 := u.listUsers(opt, int(offset), int(size))
		cnt, e2 := u.countUsers(opt)

		pbul := []*pb.User{}

		for _, cu := range ul {
			pbul = append(pbul, cu.ToPbUser())
		}

		if e1 != nil || e2 != nil {
			st := model.Status{
				Code:    errcodes.Code_INTERNAL,
				Message: fmt.Sprintf("%v;%v", e1, e2),
			}
			ret = &pb.ListUsersResponse{
				Success: false,
				Result: &pb.ListUsersResponse_Error{
					Error: st.ToPbStatus(),
				},
			}
			return
		} else {
			ret = &pb.ListUsersResponse{
				Success: true,
				Result: &pb.ListUsersResponse_Users{
					Users: &pb.Users{
						Users:  pbul,
						Total:  int32(cnt),
						Offset: offset,
						Size:   size,
					},
				},
			}
		}
		return
	}
}

func (u *UserServiceImpl) createUser(name, displayName, email string, userType pb.UserType, createdBy string) (user *model.User, secret *model.Secret, err error) {
	mc := u.env.SpawnMgoCli()
	defer mc.Close()
	q := bson.M{}
	q["name"] = name
	uid := mongo.NewObjectId()
	if mc.Exists(DbUserColl, q) {
		user = nil
		secret = nil
		err = errors.New(fmt.Sprintf("user_%s_was_created", name))
		return
	}
	user = &model.User{
		Id:          uid,
		UserType:    userType,
		Name:        name,
		Email:       email,
		CreateTime:  time.Now(),
		UpdatedTime: time.Now(),
		CreatedBy:   createdBy,
	}
	if err = mc.Insert(DbUserColl, user); err == nil {
		log.Infof("Created user %s by %s, creating secret...", name, createdBy)
		if secret, err = u.env.SecretService.updateSecret(uid.Hex()); err != nil {
			log.Errorf("Init secret failed...%v", err)
			return
		} else {
			log.Infof("Created secret for user %s", name)
		}
	}
	return
}

func (u *UserServiceImpl) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (ret *pb.CreateUserResponse, e error) {
	if ucli, err := u.ParseAuthInfo(ctx, req.Auth); err != nil || ucli.UserType != model.UserType_ADMIN {
		st := model.Status{
			Code:    errcodes.Code_PERMISSION_DENIED,
			Message: fmt.Sprintf("%v", err),
		}
		ret = &pb.CreateUserResponse{
			Success: false,
			Result: &pb.CreateUserResponse_Error{
				Error: st.ToPbStatus(),
			},
		}
		return

	} else {
		if utar, _, err := u.createUser(req.Name, req.DisplayName, req.Email, req.UserType, ucli.Name); err != nil {
			st := model.Status{
				Code:    errcodes.Code_INTERNAL,
				Message: fmt.Sprintf("%v", err),
			}
			ret = &pb.CreateUserResponse{
				Success: false,
				Result: &pb.CreateUserResponse_Error{
					Error: st.ToPbStatus(),
				},
			}
			return
		} else {
			ret = &pb.CreateUserResponse{
				Success: true,
				Result: &pb.CreateUserResponse_User{
					User: utar.ToPbUser(),
				},
			}
		}
		return
	}
}

// check by id OR name
// if the id is exists, it will skip name
// assume user is valid
func (u *UserServiceImpl) updateUser(id, name string, user *model.User) (err error) {
	if id == "" && name == "" {
		log.Infof("Both id and name are empty")
		return errors.New("bad_request")
	}
	mc := u.env.SpawnMgoCli()
	defer mc.Close()
	q := bson.M{}
	if id != "" && mongo.IsObjectIdHex(id) {
		q["_id"] = mongo.ToObjectIdHex(id)
	} else if name != "" {
		q["name"] = name
	}
	o := bson.M{}
	if user.Name != "" {
		o["name"] = user.Name
	}
	if user.DisplayName != "" {
		o["display_name"] = user.DisplayName
	}
	if user.Email != "" {
		o["email"] = user.Email
	}
	o["updated_time"] = time.Now()
	if !mc.Exists(DbUserColl, q) {
		user = nil
		err = errors.New(fmt.Sprintf("user_%s_not_exists", name))
		return
	}
	return mc.Update(DbUserColl, q, mongo.SetOperator(o))
}

func (u *UserServiceImpl) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (ret *pb.UpdateUserResponse, e error) {
	retDenied := func(err error) (ret *pb.UpdateUserResponse, e error) {
		st := model.Status{
			Code:    errcodes.Code_PERMISSION_DENIED,
			Message: fmt.Sprintf("%v", err),
		}
		ret = &pb.UpdateUserResponse{
			Success: false,
			Result: &pb.UpdateUserResponse_Error{
				Error: st.ToPbStatus(),
			},
		}
		return
	}
	if ucli, err := u.ParseAuthInfo(ctx, req.Auth); err != nil {
		return retDenied(err)
	} else if (req.Id != ucli.Id.Hex() && req.Name != ucli.Name) && ucli.UserType != model.UserType_ADMIN {
		// if NOT self, And NOT admin
		return retDenied(err)
	} else if req.Update == nil || (req.Id == "" && req.Name == "") {
		errMsg := []string{}
		if req.Id == "" {
			errMsg = append(errMsg, "id is empty")
		}
		if req.Name == "" {
			errMsg = append(errMsg, "name is empty")
		}
		if req.Update == nil {
			errMsg = append(errMsg, "update body is empty")
		}
		st := model.Status{
			Code:    errcodes.Code_INVALID_ARGUMENT,
			Message: strings.Join(errMsg, ";"),
		}
		ret = &pb.UpdateUserResponse{
			Success: false,
			Result: &pb.UpdateUserResponse_Error{
				Error: st.ToPbStatus(),
			},
		}
		return
	} else {
		uup := &model.User{
			DisplayName: req.Update.DisplayName,
		}

		updatedName := req.Name

		if ucli.UserType == model.UserType_ADMIN {
			uup.Email = req.Update.Email
			if req.Update.Name != "" {
				uup.Name = req.Update.Name
				updatedName = req.Update.Name
			}
		}

		if err := u.updateUser(req.Id, req.Name, uup); err != nil {
			st := model.Status{
				Code:    errcodes.Code_INTERNAL,
				Message: fmt.Sprintf("%v", err),
			}
			ret = &pb.UpdateUserResponse{
				Success: false,
				Result: &pb.UpdateUserResponse_Error{
					Error: st.ToPbStatus(),
				},
			}
			return
		} else {
			if utar, err := u.findUserByIdOrName(req.Id, updatedName); err != nil {
				st := model.Status{
					Code:    errcodes.Code_NOT_FOUND,
					Message: err.Error(),
				}
				ret = &pb.UpdateUserResponse{
					Success: false,
					Result: &pb.UpdateUserResponse_Error{
						Error: st.ToPbStatus(),
					},
				}
				return
			} else {
				ret = &pb.UpdateUserResponse{
					Success: true,
					Result: &pb.UpdateUserResponse_User{
						User: utar.ToPbUser(),
					},
				}
				return
			}
		}
		return
	}
}

func (u *UserServiceImpl) GetUser(ctx context.Context, req *pb.GetUserRequest) (ret *pb.GetUserResponse, e error) {
	if ucli, err := u.ParseAuthInfo(ctx, req.Auth); err != nil {
		st := model.Status{
			Code:    errcodes.Code_PERMISSION_DENIED,
			Message: fmt.Sprintf("invalid auth: %v", err.Error()),
		}
		ret = &pb.GetUserResponse{
			Success: false,
			Result: &pb.GetUserResponse_Error{
				Error: st.ToPbStatus(),
			},
		}
		return
	} else if req.Id == "$" {
		ret = &pb.GetUserResponse{
			Success: true,
			Result: &pb.GetUserResponse_User{
				User: ucli.ToPbUser(),
			},
		}
		return
	} else {
		if utar, err := u.findUserByIdOrName(req.Id, req.Name); err != nil {
			st := model.Status{
				Code:    errcodes.Code_NOT_FOUND,
				Message: err.Error(),
			}
			ret = &pb.GetUserResponse{
				Success: false,
				Result: &pb.GetUserResponse_Error{
					Error: st.ToPbStatus(),
				},
			}
			return
		} else {
			ret = &pb.GetUserResponse{
				Success: true,
				Result: &pb.GetUserResponse_User{
					User: utar.ToPbUser(),
				},
			}
			return
		}
	}
}

/* TODO:Not implemented yet
 */
func (u *UserServiceImpl) DeleteUser(context.Context, *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	return nil, errors.New("implement me")
}
