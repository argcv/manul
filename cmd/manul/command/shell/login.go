package shell

import (
	"gopkg.in/abiosoft/ishell.v2"
)

func (e *Env) AddLogin() {
	cmd := &ishell.Cmd{
		Name: "login",
		Help: "login",
		Func: func(c *ishell.Context) {
			c.Println("Login...: ", c.Args)
		},
	}

	e.Sh.AddCmd(cmd)
}
