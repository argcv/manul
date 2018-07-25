package service

import (
	"context"
	"errors"
	"fmt"
	pb "github.com/argcv/go-argcvapis/app/manul/job"
	"github.com/argcv/go-argcvapis/status/errcodes"
	"github.com/argcv/manul/model"
	"github.com/argcv/webeh/log"
	"gopkg.in/mgo.v2/bson"
	"github.com/argcv/manul/client/mongo"
	"github.com/argcv/go-argcvapis/app/manul/file"
	"os"
	"path"
	"github.com/docker/docker/client"
	"github.com/docker/docker/api/types"
	"io/ioutil"
	"github.com/docker/docker/api/types/container"
	"time"
)

type JobServiceImpl struct {
	env          *Env
	nRunningJobs int
}

/* TODO:Not implemented yet
 */
func (j *JobServiceImpl) ListJobs(context.Context, *pb.ListJobsRequest) (*pb.ListJobsResponse, error) {
	return nil, errors.New("implement me")
}

/* TODO:Not implemented yet
 */
func (j *JobServiceImpl) startJob(job *model.Job) {
	go func() {
		log.Infof("starting job... id: %v, project id: %v, user id:%v", job.Id, job.ProjectId, job.UserId)
		//
		if pj, err := j.env.ProjectService.findProject(job.ProjectId.Hex(), ""); err != nil {
			msg := fmt.Sprintf("invalid project: %v", err)
			log.Errorf(msg)
			j.updateJobLog(job.Id.Hex(), msg)
			j.updateJobProgress(job.Id.Hex(), pb.JobProgress_FAILED)
			return
		} else {
			pjwkdir := j.env.SpawnProjectWorkdir().Goto(pj.Id.Hex()).Rebase()
			log.Infof("project folder: %v", pjwkdir.Path("/"))
			jbwkdir := j.env.SpawnJobWorkdir().Goto(job.Id.Hex()).Rebase()
			log.Infof("job dir: %v", jbwkdir.Path("/"))

			cfg, e := model.LoadProjectConfig(pjwkdir.Path("manul.project.yml"))

			if e != nil {
				msg := fmt.Sprintf("load configure file failed...: %v", e)
				log.Errorf(msg)
				j.updateJobLog(job.Id.Hex(), msg)
				j.updateJobProgress(job.Id.Hex(), pb.JobProgress_FAILED)
				return
			}

			containerName := job.Id.Hex()

			// init client
			ctx := context.Background()
			cli, err := client.NewClientWithOpts(client.FromEnv, client.WithVersion("1.37"), func(c *client.Client) error {
				log.Warnf("Get Version: %v", c.ClientVersion())
				return nil
			})
			cli.ClientVersion()
			if err != nil {
				msg := fmt.Sprintf("get docker client failed!!! :%v", err)
				log.Errorf(msg)
				j.updateJobLog(job.Id.Hex(), msg)
				j.updateJobProgress(job.Id.Hex(), pb.JobProgress_FAILED)
				return
			}

			image := cfg.Image

			if image == "" {
				// default name
				image = "fedora"
			}

			//reader, err := cli.ImagePull(ctx, "docker.io/library/alpine", types.ImagePullOptions{})
			reader, err := cli.ImagePull(ctx, image, types.ImagePullOptions{})
			if err != nil {
				msg := fmt.Sprintf("pull image failed...: %v", e)
				log.Errorf(msg)
				j.updateJobLog(job.Id.Hex(), msg)
				j.updateJobProgress(job.Id.Hex(), pb.JobProgress_FAILED)
				return
			}
			pullMsg, err := ioutil.ReadAll(reader)
			if err != nil {
				msg := fmt.Sprintf("get pull image resp failed...: %v", e)
				log.Errorf(msg)
				j.updateJobLog(job.Id.Hex(), msg)
				j.updateJobProgress(job.Id.Hex(), pb.JobProgress_FAILED)
				return
			}
			log.Infof("Pull Message: [%v]", string(pullMsg))

			bse := cfg.ToBashScriptsExecutor()

			bse.Id = containerName

			workingdir := "/home/job"

			volumes := []string{
				fmt.Sprintf("%s:%s:ro", pjwkdir.Path("/"), "/home/project"),
				fmt.Sprintf("%s:%s", jbwkdir.Path("/"), workingdir),
			}

			//for _, volume := range cfg.Volume {
			//	volumes = append(volumes, fmt.Sprintf("%v/%v", basedir, volume))
			//}

			resp, err := cli.ContainerCreate(ctx, &container.Config{
				Image:      cfg.Image,
				Entrypoint: []string{"bash", "-c"},
				//Cmd:        []string{"/proc/meminfo"},
				//Cmd:        []string{"ls -a /usr"},
				Env:        cfg.Env,
				Cmd:        []string{bse.EncodedScript()},
				Tty:        true,
				Hostname:   "bootcamp",
				Domainname: "local",
				WorkingDir: workingdir,
			}, &container.HostConfig{
				//Binds: []string{
				//	//fmt.Sprintf("%v/test-data/data:/tmp:ro", basedir),
				//	fmt.Sprintf("%v/test-data/data:/tmp:ro", basedir),
				//},
				Binds: volumes,
				Resources: container.Resources{
					Memory: int64(cfg.MaximumMemMb) * 1024 * 1024,
					//Memory: 1024 * 1024 * 5, // 100 MB
					//KernelMemory: 1024 * 1024 * 1, // 100 MB
					//Memory    int64 // Memory limit (in bytes)
				},
			}, nil, containerName)
			if err != nil {
				msg := fmt.Sprintf("container create failed...: %v", e)
				log.Errorf(msg)
				j.updateJobLog(job.Id.Hex(), msg)
				j.updateJobProgress(job.Id.Hex(), pb.JobProgress_FAILED)
				return
			}

			j.updateJobLog(job.Id.Hex(), fmt.Sprintf("job started at: %v", time.Now()))
			j.updateJobProgress(job.Id.Hex(), pb.JobProgress_PENDING)

			log.Infof("resp.id: %v", resp.ID)

			if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
				msg := fmt.Sprintf("container start failed...: %v", e)
				log.Errorf(msg)
				j.updateJobLog(job.Id.Hex(), msg)
				j.updateJobProgress(job.Id.Hex(), pb.JobProgress_FAILED)
				return
			}

			statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
			select {
			case err := <-errCh:
				if err != nil {
					panic(err)
				}
			case st := <-statusCh:
				log.Infof("Wait... %v, %v", st.StatusCode, st.Error)
			}

			logs, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true})
			if err != nil {
				panic(err)
			}

			logsMsg, err := ioutil.ReadAll(logs)
			if err != nil {
				panic(err)
			}
			log.Infof("Log: [%v]", string(logsMsg))

			var expire = 1 * time.Second

			err = cli.ContainerStop(ctx, resp.ID, &expire)

			if err != nil {
				log.Errorf("Error: %v", err)
			}

			err = cli.ContainerRemove(ctx, resp.ID, types.ContainerRemoveOptions{
				RemoveVolumes: true,
				Force:         true,
			})

			if err != nil {
				log.Errorf("remove failed..: %v", err)
			}

			log.Infof("Done.")

			j.updateJobLog(job.Id.Hex(), string(logsMsg))
			j.updateJobProgress(job.Id.Hex(), pb.JobProgress_OK)

			return
		}

	}()
}

func (j *JobServiceImpl) updateJobLog(id string, log string) error {
	mc := j.env.SpawnMgoCli()
	defer mc.Close()

	q := bson.M{}
	u := bson.M{}
	if jid, e := mongo.SafeToObjectId(id); e == nil {
		q["_id"] = jid
		u = mongo.SetOperator(bson.M{
			"log": log,
		})
		return mc.Update(DbJobColl, q, u)
	} else {
		return e
	}
}
func (j *JobServiceImpl) updateJobProgress(id string, progress pb.JobProgress) error {
	mc := j.env.SpawnMgoCli()
	defer mc.Close()

	q := bson.M{}
	u := bson.M{}
	if jid, e := mongo.SafeToObjectId(id); e == nil {
		q["_id"] = jid
		u = mongo.SetOperator(bson.M{
			"progress": progress,
		})
		return mc.Update(DbJobColl, q, u)
	} else {
		return e
	}
}

func (j *JobServiceImpl) updateJob(job *model.Job) error {
	mc := j.env.SpawnMgoCli()
	defer mc.Close()
	q := bson.M{}
	q["_id"] = job.Id
	return mc.Update(DbJobColl, q, job)
}

func (j *JobServiceImpl) findJob(id string) (job *model.Job, err error) {
	mc := j.env.SpawnMgoCli()
	defer mc.Close()
	q := bson.M{}
	if oid, e := mongo.SafeToObjectId(id); e != nil {
		return nil, e
	} else {
		q["_id"] = oid
	}
	log.Infof("query: %v", q)
	err = mc.One(DbJobColl, q, &job)
	return
}

func (j *JobServiceImpl) createJob(uid, pid bson.ObjectId, files *file.Files) (job *model.Job, err error) {
	mc := j.env.SpawnMgoCli()
	defer mc.Close()

	// create job
	job = model.NewJob(uid, pid)

	// write files

	base := j.env.SpawnJobWorkdir().Goto(job.Id.Hex()).Rebase()

	// write files here
	for _, f := range files.Data {
		var perm os.FileMode = 0600
		if f.Meta != nil && f.Meta.Fields != nil {
			if cp, ok := f.Meta.Fields["perm"]; ok {
				fcp := cp.GetNumberValue()
				log.Debugf("file: [%v] => [%v] , perm: %v", f.Path, f.Name, fcp)
				if fcp > 0 && fcp <= 0777 {
					perm = os.FileMode(fcp)
				}
			}
		}
		base.WriteFile(path.Join(f.Path, f.Name), f.Data, perm)
	}

	if err = mc.Insert(DbJobColl, job); err != nil {
		return
	}
	// it is used to Start a job
	j.startJob(job)
	return
}

func (j *JobServiceImpl) CreateJob(ctx context.Context, req *pb.CreateJobRequest) (ret *pb.CreateJobResponse, e error) {
	log.Infof("Create Job...")
	if ucli, err := j.env.ParseAuthInfo(ctx, req.Auth); err != nil ||
		(ucli.UserType != model.UserType_ADMIN && ucli.UserType != model.UserType_USER) {
		st := model.Status{
			Code:    errcodes.Code_PERMISSION_DENIED,
			Message: fmt.Sprintf("invalid auth: %v", err),
		}
		ret = &pb.CreateJobResponse{
			Success: false,
			Result: &pb.CreateJobResponse_Error{
				Error: st.ToPbStatus(),
			},
		}
		return
	} else if req.Files == nil {
		st := model.Status{
			Code:    errcodes.Code_INVALID_ARGUMENT,
			Message: "missing files",
		}
		ret = &pb.CreateJobResponse{
			Success: false,
			Result: &pb.CreateJobResponse_Error{
				Error: st.ToPbStatus(),
			},
		}
		return
	} else {
		if p, err := j.env.ProjectService.findProject(req.ProjectId, ""); err != nil {
			st := model.Status{
				Code:    errcodes.Code_INVALID_ARGUMENT,
				Message: fmt.Sprintf("invalid project id: %v", err),
			}
			ret = &pb.CreateJobResponse{
				Success: false,
				Result: &pb.CreateJobResponse_Error{
					Error: st.ToPbStatus(),
				},
			}
			return
		} else if j, err := j.createJob(ucli.Id, p.Id, req.Files); err != nil {
			st := model.Status{
				Code:    errcodes.Code_INTERNAL,
				Message: fmt.Sprintf("invalid project id: %v", err),
			}
			ret = &pb.CreateJobResponse{
				Success: false,
				Result: &pb.CreateJobResponse_Error{
					Error: st.ToPbStatus(),
				},
			}
			return
		} else {
			ret = &pb.CreateJobResponse{
				Success: true,
				Result: &pb.CreateJobResponse_Job{
					Job: j.ToPbJob(),
				},
			}
			return
		}
	}
}

func (j *JobServiceImpl) GetJob(ctx context.Context, req *pb.GetJobRequest) (ret *pb.GetJobResponse, e error) {
	log.Infof("Get Job...")
	if ucli, err := j.env.ParseAuthInfo(ctx, req.Auth); err != nil ||
		(ucli.UserType != model.UserType_ADMIN && ucli.UserType != model.UserType_USER) {
		st := model.Status{
			Code:    errcodes.Code_PERMISSION_DENIED,
			Message: fmt.Sprintf("invalid auth: %v", err),
		}
		ret = &pb.GetJobResponse{
			Success: false,
			Result: &pb.GetJobResponse_Error{
				Error: st.ToPbStatus(),
			},
		}
		return
	} else if job, err := j.findJob(req.Id); err != nil {
		st := model.Status{
			Code:    errcodes.Code_NOT_FOUND,
			Message: fmt.Sprintf("job not found: %v", e),
		}
		ret = &pb.GetJobResponse{
			Success: false,
			Result: &pb.GetJobResponse_Error{
				Error: st.ToPbStatus(),
			},
		}
		return
	} else {
		ret = &pb.GetJobResponse{
			Success: true,
			Result: &pb.GetJobResponse_Job{
				Job: job.ToPbJob(),
			},
		}
		return
	}
}

/* TODO:Not implemented yet
 */
func (j *JobServiceImpl) CancelJob(context.Context, *pb.CancelJobRequest) (*pb.CancelJobResponse, error) {
	return nil, errors.New("implement me")
}
