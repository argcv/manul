package shell

import (
	"gopkg.in/abiosoft/ishell.v2"
	"strings"
)

func (e *Env) AddCmd(cmd *ishell.Cmd) {
	e.Sh.AddCmd(cmd)
}

func (e *Env) Run() {
	e.Sh.Run()
}

func (e *Env) Process(args ...string) {
	e.Sh.Process(args...)
}

// forward prints
func (e *Env) Println(val ...interface{}) {
	e.Sh.Println(val...)
}

func (e *Env) Print(val ...interface{}) {
	e.Sh.Print(val...)
}

func (e *Env) Printf(format string, val ...interface{}) {
	e.Sh.Printf(format, val...)
}

func (e *Env) PushPrompt(str string) {
	e.prompts = append(e.prompts, str)
	e.Sh.SetPrompt(str)
}

func (e *Env) ConfirmRepl(c *ishell.Context, q string, or bool) bool {
	if or {
		e.PushPrompt("[Y/n] :")
	} else {
		e.PushPrompt("[y/N] :")
	}
	c.Println(q)
	defer e.PopPrompt()
	rtMsg := strings.TrimSpace(c.ReadLine())
	if rtMsg == "n" || rtMsg == "N" {
		return false
	} else if rtMsg == "y" || rtMsg == "Y" {
		return true
	} else {
		return or
	}
}

func (e *Env) GetStringRepl(c *ishell.Context, q, prompt string) string {
	return e.GetStringReplFunc(c, q, prompt, func(s string) bool {
		return true
	})
}

func (e *Env) GetStringReplNonEmpty(c *ishell.Context, q, prompt string) string {
	return e.GetStringReplFunc(c, q, prompt, func(s string) bool {
		return len(strings.TrimSpace(s)) > 0
	})
}

func (e *Env) GetStringReplFunc(c *ishell.Context, q, prompt string, f func(string) bool) string {
	if len(prompt) > 0 {
		e.PushPrompt(prompt)
		defer e.PopPrompt()
	}
	c.Println(q)
	var rtMsg string
	for rtMsg = c.ReadLine(); !(f(rtMsg)); rtMsg = c.ReadLine() {
		c.Println("Invalid Input")
	}
	return rtMsg
}
