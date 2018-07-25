package helpers

import (
	"fmt"
	"os/exec"
	"strings"
)

type BashScriptStage struct {
	Name    string
	Scripts []string
}

func NewBashScriptStage(name string, scripts ...string) *BashScriptStage {
	return &BashScriptStage{
		Name:    name,
		Scripts: scripts,
	}
}

func (e *BashScriptStage) EncodedScripts(job string) []string {
	scripts := []string{
		fmt.Sprintf("echo 'manul:%s:stage:start:%s'", job, e.Name),
	}
	for _, script := range e.Scripts {
		prompt := fmt.Sprintf("echo '$ %s'", strings.Replace(script, "'", "'\\''", -1))
		scripts = append(scripts, prompt, script)
	}
	scripts = append(scripts, fmt.Sprintf("echo 'manul:%s:stage:end:%s'", job, e.Name))
	return scripts
}

func (e *BashScriptStage) AddScripts(scripts ...string) *BashScriptStage {
	e.Scripts = append(e.Scripts, scripts...)
	return e
}

func (e *BashScriptStage) SetScripts(scripts ...string) *BashScriptStage {
	e.Scripts = scripts
	return e
}

type BashScriptsExecutor struct {
	Id     string
	Env    []string
	Stages []BashScriptStage
}

func NewBashScriptsExecutor(id string, env ...string) *BashScriptsExecutor {
	return &BashScriptsExecutor{
		Id:  id,
		Env: env,
	}
}

func (e *BashScriptsExecutor) AddEnv(env ...string) *BashScriptsExecutor {
	e.Env = append(e.Env, env...)
	return e
}

func (e *BashScriptsExecutor) SetEnv(env ...string) *BashScriptsExecutor {
	e.Env = env
	return e
}

func (e *BashScriptsExecutor) AddStage(stage *BashScriptStage) *BashScriptsExecutor {
	e.Stages = append(e.Stages, *stage)
	return e
}

func (e *BashScriptsExecutor) AddScriptsInStage(name string, scripts ...string) *BashScriptsExecutor {
	stage := NewBashScriptStage(name, scripts...)
	e.Stages = append(e.Stages, *stage)
	return e
}

func (e *BashScriptsExecutor) EncodedScript() string {
	scripts := []string{
		"set -Eeo pipefail",
	}
	for _, stage := range e.Stages {
		scripts = append(scripts, stage.EncodedScripts(e.Id)...)
	}
	return strings.Join(scripts, ";")
}

func (e *BashScriptsExecutor) Perform() ([]byte, error) {
	cfg := exec.Command("bash", "-c", e.EncodedScript())
	cfg.Env = e.Env
	return cfg.CombinedOutput()
}
