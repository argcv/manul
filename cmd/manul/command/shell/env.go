package shell

import (
	"gopkg.in/abiosoft/ishell.v2"
	"strings"
)

type Env struct {
	Sh      *ishell.Shell
	prompts []string
}

func NewEnv() (e *Env) {
	shell := ishell.New()

	e = &Env{
		Sh:      shell,
		prompts: []string{},
	}

	e.init()
	return
}

func (e *Env) PopPrompt() {
	if len(e.prompts) > 1 {
		e.prompts = e.prompts[0 : len(e.prompts)-1]
		e.Sh.SetPrompt(e.prompts[len(e.prompts)-1])
	} else {
		e.Printf("Unexpected PopPrompt...?")
	}
}

func (e *Env) init() {
	sh := e.Sh
	// display welcome info.
	e.Println("Manul Interactive.")
	e.Println("Type Help to get more information")

	e.PushPrompt("manul > ")

	sh.NotFound(func(c *ishell.Context) {
		c.Printf("Command Not Recognized : [%v]\n", strings.Join(c.Args, "], ["))
	})

	e.AddLogin()
	e.AddSetup()
}
