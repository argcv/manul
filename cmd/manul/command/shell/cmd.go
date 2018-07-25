package shell

import (
	"github.com/spf13/cobra"
)

func NewShellCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "shell",
		Short: "A simple interactive interface",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			env := NewEnv()
			env.Run()
			return
		},
	}
	return cmd
}

func NewSetupCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "setup",
		Short: "setup server environment (not for client side)",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			env := NewEnv()
			env.Process("setup")
			return
		},
	}
	return cmd
}

func NewLoginCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "login",
		Short: "user login",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			env := NewEnv()
			env.Process(append([]string{"login"}, args...)...)
			return
		},
	}
	return cmd
}
