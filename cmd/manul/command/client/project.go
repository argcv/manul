package client

import (
	"github.com/argcv/go-argcvapis/app/manul/file"
	"github.com/argcv/go-argcvapis/app/manul/project"
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

func NewManulProjectSubmitCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "submit",
		Short: "Create a new project",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			conn := NewGrpcConn()
			defer conn.Close()
			pcli := conn.NewProjectCli()

			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			auth := config.GetAuthInfo()
			files := &file.Files{}

			dir, _ := cmd.Flags().GetString("dir")

			dir = path.Clean(dir)

			base, _ := os.Getwd()

			if dir[0] != '/' {
				// relative path
				dir = path.Join(base, dir)
			}

			fs := workdir.NewWorkdir(dir)

			_, lastDir := fs.Split()

			name, _ := cmd.Flags().GetString("name")
			desc, _ := cmd.Flags().GetString("desc")

			if name == "" {
				name = lastDir
			}

			log.Infof("user: %v", auth.Name)
			log.Infof("name: %v", name)
			log.Infof("desc: %v", desc)
			log.Infof("base: %v", fs.GetCwd())
			fs.IterFiles("/", func(f *file.File) error {
				log.Infof("file: [%v], [%v]", f.Path, f.Name)
				files.Data = append(files.Data, f)
				return nil
			})
			p := &project.Project{
				Name:  name,
				Desc:  desc,
				Files: files,
			}

			req := &project.CreateProjectRequest{
				Auth:    auth,
				Project: p,
			}
			if ret, err := pcli.CreateProject(ctx, req); err != nil {
				log.Infof("Error: %v", err)
			} else {
				log.Infof("Dump: %v", spew.Sdump(ret))
			}

			return
		},
	}
	cmd.PersistentFlags().StringP("dir", "d", ".", "project base dir")
	cmd.PersistentFlags().StringP("name", "n", "", "project name")
	cmd.PersistentFlags().String("desc", "", "project description")
	return cmd
}

func NewManulProjectCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "project",
		Short: "Project Operations",
		PreRunE: func(cmd *cobra.Command, args []string) (err error) {
			return nil
		},
	}
	cmd.AddCommand(
		NewManulProjectSubmitCommand(),
	)
	return cmd
}
