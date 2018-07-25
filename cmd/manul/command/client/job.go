package client

import (
	"github.com/argcv/go-argcvapis/app/manul/file"
	"github.com/argcv/go-argcvapis/app/manul/job"
	"github.com/argcv/manul/client/workdir"
	"github.com/argcv/manul/config"
	"github.com/argcv/webeh/log"
	"github.com/davecgh/go-spew/spew"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"os"
	"path"
	"time"
)

func NewManulJobSubmitCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "submit",
		Short: "Job submitting",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			conn := NewGrpcConn()
			defer conn.Close()
			jcli := conn.NewJobCli()

			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			auth := config.GetAuthInfo()
			files := &file.Files{}

			dir, _ := cmd.Flags().GetString("dir")

			dir = path.Clean(dir)

			pid, _ := cmd.Flags().GetString("project")

			base, _ := os.Getwd()

			if dir[0] != '/' {
				// relative path
				dir = path.Join(base, dir)
			}

			fs := workdir.NewWorkdir(dir)

			_, lastDir := fs.Split()

			name, _ := cmd.Flags().GetString("name")

			if name == "" {
				name = lastDir
			}

			log.Infof("user: %v", auth.Name)
			log.Infof("name: %v", name)
			log.Infof("base: %v", fs.GetCwd())
			fs.IterFiles("/", func(f *file.File) error {
				log.Infof("file: [%v], [%v]", f.Path, f.Name)
				files.Data = append(files.Data, f)
				return nil
			})

			req := &job.CreateJobRequest{
				Auth:      auth,
				ProjectId: pid,
				Files:     files,
			}
			if ret, err := jcli.CreateJob(ctx, req); err != nil {
				log.Infof("Error: %v", err)
			} else {
				log.Infof("Dump: %v", spew.Sdump(ret))
			}

			return
		},
	}
	cmd.PersistentFlags().String("project", "p", "project id")
	cmd.PersistentFlags().StringP("dir", "d", ".", "job base dir")
	return cmd
}

func NewManulJobCheckCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "check",
		Short: "Job Check",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			conn := NewGrpcConn()
			defer conn.Close()

			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			jcli := conn.NewJobCli()

			auth := config.GetAuthInfo()

			for _, jid := range args {
				req := &job.GetJobRequest{
					Auth: auth,
					Id:   jid,
				}

				if ret, err := jcli.GetJob(ctx, req); err != nil {
					log.Infof("Error: %v", err)
				} else {
					log.Infof("Dump: %v", spew.Sdump(ret))
					log.Infof("success: %v", ret.Success)
					log.Infof("msg: %v", ret.Message)
					if ret.Success {
						log.Infof(ret.GetJob().Logs)
					}
				}


			}

			return
		},
	}
	return cmd
}

func NewManulJobCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "job",
		Short: "Job Operations",
		PreRunE: func(cmd *cobra.Command, args []string) (err error) {
			return
		},
	}
	cmd.AddCommand(
		NewManulJobSubmitCommand(),
		NewManulJobCheckCommand(),
	)
	return cmd
}
