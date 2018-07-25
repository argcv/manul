package main

import (
	"fmt"
	"github.com/argcv/manul/model"
	"github.com/argcv/webeh/log"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
	"io/ioutil"
	"os"
	"path"
	"time"
)

func main() {
	//basedir, e := helpers.GetPathOfSelf()
	//
	//if e != nil {
	//	panic("get current path failed")
	//}
	basedir, e := os.Getwd()

	if e != nil {
		panic("get current path failed")
	}

	cfg, e := model.LoadProjectConfig(path.Join(basedir, "manul.project.yml.default"))

	if e != nil {
		panic(fmt.Sprintf("load configure file failed...: %v", e))
	}

	containerName := "greeting"
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithVersion("1.37"), func(c *client.Client) error {
		log.Warnf("Get Version: %v", c.ClientVersion())
		return nil
	})
	cli.ClientVersion()
	if err != nil {
		panic(err)
	}

	images, err := cli.ImageList(ctx, types.ImageListOptions{})

	for _, image := range images {
		log.Infof("IMG: %v %v %v", image.ID, image.RepoTags, image.RepoDigests)
	}

	//reader, err := cli.ImagePull(ctx, "docker.io/library/alpine", types.ImagePullOptions{})
	reader, err := cli.ImagePull(ctx, "fedora", types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}
	pullMsg, err := ioutil.ReadAll(reader)
	if err != nil {
		panic(err)
	}
	log.Infof("Pull Message: [%v]", string(pullMsg))

	bse := cfg.ToBashScriptsExecutor()

	bse.Id = "Job1"

	volumes := []string{}

	for _, volume := range cfg.Volume {
		volumes = append(volumes, fmt.Sprintf("%v/%v", basedir, volume))
	}

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image:      cfg.Image,
		Entrypoint: []string{"bash", "-c"},
		//Cmd:        []string{"/proc/meminfo"},
		//Cmd:        []string{"ls -a /usr"},
		Env:        cfg.Env,
		Cmd:        []string{bse.EncodedScript()},
		Tty:        true,
		Hostname:   "bootcamp",
		Domainname: "local",
	}, &container.HostConfig{
		//Binds: []string{
		//	//fmt.Sprintf("%v/test-data/data:/tmp:ro", basedir),
		//	fmt.Sprintf("%v/test-data/data:/tmp:ro", basedir),
		//},
		Binds: volumes,
		Resources: container.Resources{
			Memory: int64(cfg.MaximumMemMb) * 1024 * 1024,
			//Memory: 1024 * 1024 * 5, // 100 MB
			//KernelMemory: 1024 * 1024 * 1, // 100 MB
			//Memory    int64 // Memory limit (in bytes)
		},
	}, nil, containerName)
	if err != nil {
		panic(err)
	}

	log.Infof("resp.id: %v", resp.ID)

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			panic(err)
		}
	case st := <-statusCh:
		log.Infof("Wait... %v, %v", st.StatusCode, st.Error)
	}

	logs, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true})
	if err != nil {
		panic(err)
	}

	logsMsg, err := ioutil.ReadAll(logs)
	if err != nil {
		panic(err)
	}
	log.Infof("Log: [%v]", string(logsMsg))

	var expire = 1 * time.Second

	err = cli.ContainerStop(ctx, resp.ID, &expire)

	if err != nil {
		log.Errorf("Error: %v", err)
	}

	err = cli.ContainerRemove(ctx, resp.ID, types.ContainerRemoveOptions{
		RemoveVolumes: true,
		Force:         true,
	})

	if err != nil {
		panic(err)
	}

	log.Infof("Done.")
}
