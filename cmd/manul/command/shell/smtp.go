package shell

import (
	"gopkg.in/abiosoft/ishell.v2"
	"github.com/argcv/manul/client/mail"
	"github.com/argcv/webeh/log"
	"github.com/davecgh/go-spew/spew"
)

func (e *Env) AddSmtp() {
	cmd := &ishell.Cmd{
		Name: "smtp",
		Help: "smtp sender",
		Func: func(c *ishell.Context) {
			if session, err := mail.NewSMTPSession(); err != nil {
				log.Errorf("new smtp session failed: %v", err)
			} else {
				session.DefaultFrom()
				session.To("yujing5b5d@gmail.com")
				session.Subject("Hello~")
				session.HtmlBody("<h1>Some Content Here</h1><p>Hello</p>")
				spew.Dump(session)
				if e := session.Perform(); e != nil {
					log.Errorf("error: %v", e)
				}
			}
			c.Println("Sent....: ", c.Args)
		},
	}

	e.Sh.AddCmd(cmd)
}
