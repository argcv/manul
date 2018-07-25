package client

import (
	"fmt"
	"github.com/argcv/go-argcvapis/app/manul/project"
	"github.com/argcv/go-argcvapis/app/manul/user"
	"github.com/argcv/manul/config"
	"github.com/argcv/webeh/log"
	"google.golang.org/grpc"
	"github.com/argcv/go-argcvapis/app/manul/job"
)

type RpcConn struct {
	C *grpc.ClientConn
}

func NewGrpcConn() (c *RpcConn) {
	rpcHost := config.GetRpcHost()
	rpcPort := config.GetRpcPort()
	serverAddr := fmt.Sprintf("%s:%d", rpcHost, rpcPort)
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}
	c = &RpcConn{}
	conn, err := grpc.Dial(serverAddr, opts...)
	c.C = conn
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	return
}

func (c *RpcConn) Close() {
	c.C.Close()
}

func (c *RpcConn) NewProjectCli() project.ProjectServiceClient {
	return project.NewProjectServiceClient(c.C)
}

func (c *RpcConn) NewUserCli() user.UserServiceClient {
	return user.NewUserServiceClient(c.C)
}

func (c *RpcConn) NewJobCli() job.JobServiceClient {
	return job.NewJobServiceClient(c.C)
}
