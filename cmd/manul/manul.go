package main

import (
	"fmt"
	"github.com/argcv/manul/cmd/manul/command/client"
	"github.com/argcv/manul/cmd/manul/command/server"
	"github.com/argcv/manul/cmd/manul/command/shell"
	"github.com/argcv/manul/version"
	configeh "github.com/argcv/webeh/config"
	"github.com/argcv/webeh/log"
	"github.com/spf13/cobra"
	"math/rand"
	"os"
	"path"
	"time"
)

var (
	rootCmd = &cobra.Command{
		Use:   "manul",
		Short: "Manul is an Auto Grader",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			log.Infof("Manul version: %s (%s) Built At: %v", version.Version, version.GitHash, version.BuildDate)
			if verbose, err := cmd.Flags().GetBool("verbose"); err == nil {
				if verbose {
					log.Verbose()
					log.Debug("verbose mode: ON")
				}
			}

			conf, _ := cmd.Flags().GetString("config")

			if e := configeh.LoadConfig(configeh.Option{
				Project:        "manul",
				Path:           conf,
				DefaultPath:    path.Join(os.Getenv("HOME"), ".manul"),
				FileMustExists: true,
			}); e != nil {
				return e
			}

			// set rand seed
			rand.Seed(time.Now().Unix())

			return nil
		},
	}
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Prints the version of manul.",
		// do not execute any persistent actions
		PersistentPreRun: func(cmd *cobra.Command, args []string) {},
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("manul version:", version.Version)
			fmt.Println("Git commit hash:", version.GitHash)
			if version.BuildDate != "" {
				fmt.Println("Build date:", version.BuildDate)
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(
		versionCmd,
		server.NewManulRpcServerCommand(),
		client.NewManulProjectCommand(),
		client.NewManulJobCommand(),
		shell.NewShellCommand(),
		shell.NewSetupCommand(),
		shell.NewLoginCommand(),
	)
	rootCmd.PersistentFlags().StringP("config", "c", "", "explicit assign a configuration file")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "log verbose")

}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Infof("%v", err)
		os.Exit(1)
	}
}
