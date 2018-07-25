package model

import (
	pb "github.com/argcv/go-argcvapis/app/manul/job"
	"github.com/argcv/manul/client/mongo"
	"gopkg.in/mgo.v2/bson"
)

type Job struct {
	Id        bson.ObjectId  `bson:"_id,omitempty" json:"id"`
	ProjectId bson.ObjectId  `bson:"project_id,omitempty" json:"project_id,omitempty"`
	UserId    bson.ObjectId  `bson:"user_id,omitempty" json:"user_id,omitempty"`
	Progress  pb.JobProgress `bson:"progress,omitempty" json:"progress,omitempty"`
	Score     int32          `bson:"score,omitempty" json:"score,omitempty"`
	Log       string         `bson:"log,omitempty" json:"log,omitempty"`
	Error     *Status        `bson:"error,omitempty" json:"error,omitempty"`
}

func NewJob(uid bson.ObjectId, pid bson.ObjectId) *Job {
	jid := mongo.NewObjectId()
	job := &Job{
		Id:        jid,
		UserId:    uid,
		ProjectId: pid,
		Progress:  pb.JobProgress_CREATED,
		Score:     0,
		Log:       "",
		Error:     nil,
	}
	return job
}

// from pb & to pb

func (j *Job) ToPbJob() (pbJob *pb.Job) {
	result := &pb.JobResult{
		Score: j.Score,
	}

	if j.Error != nil {
		result.Error = j.Error.ToPbStatus()
	}

	pbJob = &pb.Job{
		Id:        j.Id.Hex(),
		ProjectId: j.ProjectId.Hex(),
		UserId:    j.UserId.Hex(),
		Progress:  j.Progress,
		Result:    result,
		Logs:      j.Log,
	}
	return
}

func FromPbJob(j *pb.Job) *Job {
	var score int32 = 0

	var err *Status = nil

	if j.Result != nil {
		score = j.Result.Score
		err = FromPbStatus(j.Result.Error)
	}

	return &Job{
		Id:        mongo.SafeToObjectIdOrEmpty(j.Id),
		ProjectId: mongo.SafeToObjectIdOrEmpty(j.ProjectId),
		UserId:    mongo.SafeToObjectIdOrEmpty(j.UserId),
		Progress:  j.Progress,
		Score:     score,
		Log:       j.Logs,
		Error:     err,
	}
}
