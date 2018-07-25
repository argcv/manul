package model

import (
	pb "github.com/argcv/go-argcvapis/app/manul/project"
	"github.com/argcv/manul/helpers"
	"github.com/argcv/webeh/log"
	"github.com/davecgh/go-spew/spew"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

var (
	ProjectConfigDefaultStages = []string{
		"install",
		"build",
		"test",
	}
)

type ProjectJobConfig struct {
	Script  string   `bson:"script,omitempty" json:"script" yaml:"script"`
	Scripts []string `bson:"scripts,omitempty" json:"scripts" yaml:"scripts"`
}

type ProjectChecklistElem struct {
	Path       string `bson:"path,omitempty" json:"path" yaml:"check_list"`
	TargetType int    `bson:"target_type,omitempty" json:"target_type" yaml:"target_type"`
}

type ProjectConfig struct {
	Image        string                      `bson:"image,omitempty" json:"image" yaml:"image"`
	Env          []string                    `bson:"env,omitempty" json:"env" yaml:"env"`
	Volume       []string                    `bson:"volume,omitempty" json:"volume" yaml:"volume"`
	Stages       []string                    `bson:"stage,omitempty" json:"stage" yaml:"stage"`
	Jobs         map[string]ProjectJobConfig `bson:"job,omitempty" json:"job" yaml:"job"`
	Checklist    []*pb.ProjectChecklistElem  `bson:"check_list,omitempty" json:"check_list" yaml:"check_list"`
	TimeoutSec   uint64                      `bson:"timeout_sec,omitempty" json:"timeout_sec" yaml:"timeout_sec"`
	MaximumCpu   uint64                      `bson:"maximum_cpu,omitempty" json:"maximum_cpu" yaml:"maximum_cpu"`
	MaximumMemMb uint64                      `bson:"maximum_mem_mb,omitempty" json:"maximum_mem_mb" yaml:"maximum_mem_mb"`
}

func LoadProjectConfig(path string) (*ProjectConfig, error) {
	if data, err := ioutil.ReadFile(path); err != nil {
		return nil, err
	} else {
		cfg := &ProjectConfig{}
		err = yaml.Unmarshal([]byte(data), cfg)
		if err != nil {
			return nil, err
		} else {
			return cfg, err
		}
	}
}

func (c *ProjectConfig) ToBashScriptsExecutor() *helpers.BashScriptsExecutor {
	e := helpers.NewBashScriptsExecutor("", c.Env...)
	log.Infof("dump: [%v]", spew.Sdump(c))
	if len(c.Stages) == 0 {
		c.Stages = ProjectConfigDefaultStages
	}
	for _, stage := range c.Stages {
		if job, ok := c.Jobs[stage]; ok {
			scripts := []string{}
			if len(job.Script) > 0 {
				scripts = append(scripts, job.Script)
			}
			scripts = append(scripts, job.Scripts...)
			e.AddScriptsInStage(stage, scripts...)
		} else {
			log.Infof("Stage: %v NOT found...", stage)
			e.AddScriptsInStage(stage)
		}
	}
	return e
}

func (c *ProjectConfig) ToPbProjectConfig(rich bool) (pbConfig *pb.ProjectConfig) {
	pbConfig = &pb.ProjectConfig{
		Image:        c.Image,
		Checklist:    c.Checklist,
		TimeoutSec:   c.TimeoutSec,
		MaximumCpu:   c.MaximumCpu,
		MaximumMemMb: c.MaximumMemMb,
	}
	return
}
