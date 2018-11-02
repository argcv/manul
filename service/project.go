package service

import (
	"context"
	"errors"
	"fmt"
	pb "github.com/argcv/go-argcvapis/app/manul/project"
	"github.com/argcv/go-argcvapis/status/errcodes"
	"github.com/argcv/manul/client/mongo"
	"github.com/argcv/manul/client/workdir"
	"github.com/argcv/manul/model"
	"github.com/argcv/webeh/log"
	"gopkg.in/mgo.v2/bson"
	"os"
	"path"
	"sync"
	"time"
)

type ProjectServiceImpl struct {
	env  *Env
	muFs *sync.Mutex
}

func NewProjectServiceImpl(env *Env) *ProjectServiceImpl {
	return &ProjectServiceImpl{
		env:  env,
		muFs: &sync.Mutex{},
	}
}

func (p *ProjectServiceImpl) InitProjectWorkdir(id string) (err error) {
	p.muFs.Lock()
	defer p.muFs.Unlock()
	base := p.env.SpawnProjectWorkdir().Goto(id)

	if base.Exists("/") || base.IsDir("/") {
		return errors.New("already exists")
	} else {
		// init folder
		base.MkdirAll("/", 0700)
	}
	return
}

/* TODO:Not implemented yet
 */
func (p *ProjectServiceImpl) ListProjects(context.Context, *pb.ListProjectsRequest) (*pb.ListProjectsResponse, error) {
	return nil, errors.New("implement me")
}

func (p *ProjectServiceImpl) createProject(base *workdir.Workdir, id, name, desc, createdBy string) (proj *model.Project, err error) {
	mc := p.env.SpawnMgoCli()
	defer mc.Close()

	pid, e := mongo.SafeToObjectId(id)

	if e != nil {
		pid = mongo.NewObjectId()
	}

	projCfg, e := model.LoadProjectConfig(base.Path("manul.project.yml"))

	if e != nil {
		log.Errorf("Create failed... removing folder : %v err: %v", base.GetCwd(), base.RemoveCwd())
		return nil, e
	}

	proj = &model.Project{
		Id:          pid,
		Name:        name,
		Desc:        desc,
		CreateTime:  time.Now(),
		UpdatedTime: time.Now(),
		CreatedBy:   createdBy,
		Config:      projCfg,
	}
	err = mc.Insert(DbProjectColl, proj)
	return
}

/* TODO:Not implemented yet
 */
func (p *ProjectServiceImpl) CreateProject(ctx context.Context, req *pb.CreateProjectRequest) (ret *pb.CreateProjectResponse, e error) {
	log.Infof("Create project...")
	if ucli, err := p.env.ParseAuthInfo(ctx, req.Auth); err != nil || ucli.UserType != model.UserType_ADMIN {
		st := model.Status{
			Code:    errcodes.Code_PERMISSION_DENIED,
			Message: fmt.Sprintf("invalid auth: %v", err),
		}
		ret = &pb.CreateProjectResponse{
			Success: false,
			Result: &pb.CreateProjectResponse_Error{
				Error: st.ToPbStatus(),
			},
		}
		return
	} else {
		reqProj := req.Project
		id := mongo.NewObjectId().Hex()
		name := reqProj.Name
		desc := reqProj.Desc

		if err := p.InitProjectWorkdir(id); err != nil {
			st := model.Status{
				Code:    errcodes.Code_INTERNAL,
				Message: "unexpected folder already exists",
			}
			ret = &pb.CreateProjectResponse{
				Success: false,
				Result: &pb.CreateProjectResponse_Error{
					Error: st.ToPbStatus(),
				},
			}
			return
		} else {
			base := p.env.SpawnProjectWorkdir().Goto(id).Rebase()

			// write files here
			for _, f := range reqProj.Files.Data {
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

			if rp, err := p.createProject(base, id, name, desc, ucli.Name); err != nil {
				st := model.Status{
					Code:    errcodes.Code_INTERNAL,
					Message: fmt.Sprintf("Internal Error: %v", err),
				}
				ret = &pb.CreateProjectResponse{
					Success: false,
					Result: &pb.CreateProjectResponse_Error{
						Error: st.ToPbStatus(),
					},
				}
				return
			} else {
				ret = &pb.CreateProjectResponse{
					Success: true,
					Result: &pb.CreateProjectResponse_Project{
						Project: rp.ToPbProject(true),
					},
				}
				return
			}
		}
	}
	return
}

/* TODO:Not implemented yet
 */
func (p *ProjectServiceImpl) UpdateProject(context.Context, *pb.UpdateProjectRequest) (*pb.UpdateProjectResponse, error) {
	return nil, errors.New("implement me")
}

func (s *ProjectServiceImpl) findProject(id string, name string) (project *model.Project, err error) {
	mc := s.env.SpawnMgoCli()
	defer mc.Close()
	q := bson.M{}
	badreq := true
	if pid, e := mongo.SafeToObjectId(id); e == nil {
		q["_id"] = pid
		badreq = false
	}
	if name != "" {
		q["name"] = name
		badreq = false
	}

	if badreq {
		return nil, errors.New("missing args")
	}

	log.Infof("query: %v", q)
	err = mc.One(DbProjectColl, q, &project)
	return
}

/* TODO:Not implemented yet
 */
func (p *ProjectServiceImpl) GetProject(context.Context, *pb.GetProjectRequest) (*pb.GetProjectResponse, error) {
	return nil, errors.New("implement me")
}

/* TODO:Not implemented yet
 */
func (p *ProjectServiceImpl) GetProjectChecklist(context.Context, *pb.GetProjectChecklistRequest) (*pb.GetProjectChecklistResponse, error) {
	return nil, errors.New("implement me")
}

/* TODO:Not implemented yet
 */
func (p *ProjectServiceImpl) DeleteProject(context.Context, *pb.DeleteProjectRequest) (*pb.DeleteProjectResponse, error) {
	return nil, errors.New("implement me")
}

/* TODO:Not implemented yet
 */
func (p *ProjectServiceImpl) UpdateProjectMember(context.Context, *pb.UpdateProjectMemberRequest) (*pb.UpdateProjectMemberResponse, error) {
	return nil, errors.New("implement me")
}

/* TODO:Not implemented yet
 */
func (p *ProjectServiceImpl) ListProjectMembers(context.Context, *pb.ListProjectMembersRequest) (*pb.ListProjectMembersResponse, error) {
	return nil, errors.New("implement me")
}
