package model

import (
	pb "github.com/argcv/go-argcvapis/app/manul/project"
	"github.com/argcv/manul/client/mongo"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Project struct {
	// Required: project id
	Id bson.ObjectId `bson:"_id,omitempty" json:"id"`
	// Required: project name
	Name string `bson:"name,omitempty" json:"name,omitempty"`
	// Optional: description
	Desc string `bson:"desc,omitempty" json:"desc,omitempty"`

	Config *ProjectConfig `bson:"config,omitempty" json:"config,omitempty"`

	// Meta is NOT in using
	//Meta *structpb.Struct `bson:"meta,omitempty" json:"meta,omitempty"`
	CreateTime  time.Time `bson:"create_time,omitempty" json:"create_time,omitempty"`
	UpdatedTime time.Time `bson:"updated_time,omitempty" json:"updated_time,omitempty"`
	CreatedBy   string    `bson:"created_by,omitempty" json:"created_by"`
}

func (p *Project) ToPbProject(rich bool) (pbProject *pb.Project) {
	pbProject = &pb.Project{
		Id:   p.Id.Hex(),
		Name: p.Name,
		Desc: p.Desc,
	}
	if rich {
		pbProject.Config = p.Config.ToPbProjectConfig(rich)
	}
	return
}

func FromPbProject(p *pb.Project) *Project {
	return &Project{
		Id:   mongo.SafeToObjectIdOrEmpty(p.Id),
		Name: p.Name,
		Desc: p.Desc,
		//Meta: p.Meta,
	}
}
