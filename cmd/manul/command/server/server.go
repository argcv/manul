package server

import (
	"fmt"
	"github.com/argcv/go-argcvapis/app/manul/job"
	"github.com/argcv/go-argcvapis/app/manul/project"
	"github.com/argcv/go-argcvapis/app/manul/secret"
	"github.com/argcv/go-argcvapis/app/manul/user"
	"github.com/argcv/manul/config"
	"github.com/argcv/manul/service"
	"github.com/argcv/webeh/log"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"net"
	"net/http"
)

func NewManulRpcServerCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Rpc Server",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			rpcBind := config.GetRpcBind()
			rpcPort := config.GetRpcPort()
			httpBind := config.GetHttpBind()
			httpPort := config.GetHttpPort()

			log.Infof("The RPC Server starting.. %s:%d", rpcBind, rpcPort)
			log.Infof("The HTTP Server starting.. %s:%d", httpBind, httpPort)
			lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", rpcBind, rpcPort))
			if err != nil {
				log.Infof("failed to listen: %v", err)
				return
			}

			maxRecvMsgSizeMB := config.GetRpcOptionMaxRecvMsgSizeMB()
			maxSendMsgSizeMB := config.GetRpcOptionMaxSendMsgSizeMB()

			opts := []grpc.ServerOption{
				grpc.MaxRecvMsgSize(1024 * 1024 * maxRecvMsgSizeMB),
				grpc.MaxSendMsgSize(1024 * 1024 * maxSendMsgSizeMB),
			}

			log.Infof("[Option] MaxRecvMsgSize(mb): %d", maxRecvMsgSizeMB)
			log.Infof("[Option] MaxSendMsgSize(mb): %d", maxSendMsgSizeMB)

			grpcServer := grpc.NewServer(opts...)

			env, err := service.NewManulGlobalEnv()

			if err != nil {
				return err
			}

			// register servers here
			user.RegisterUserServiceServer(grpcServer, env.UserService)
			secret.RegisterSecretServiceServer(grpcServer, env.SecretService)
			project.RegisterProjectServiceServer(grpcServer, env.ProjectService)
			job.RegisterJobServiceServer(grpcServer, env.JobService)

			go func() {
				log.Infof("Sterted Rpc... %s:%d", rpcBind, rpcPort)
				grpcServer.Serve(lis)
			}()

			// http...

			listenString := fmt.Sprintf("%s:%d", httpBind, httpPort)

			ctx := context.Background()
			ctx, cancel := context.WithCancel(ctx)
			defer cancel()

			gwmux := runtime.NewServeMux()

			//err = pb.RegisterTasksHandlerFromEndpoint(ctx, gwmux, fmt.Sprintf("%s:%d", "127.0.0.1", rpcPort),
			//	[]grpc.DialOption{grpc.WithInsecure()},
			//)

			rpcProxy := fmt.Sprintf("%s:%d", "127.0.0.1", rpcPort)

			// bind endpoints start
			err = user.RegisterUserServiceHandlerFromEndpoint(ctx, gwmux, rpcProxy,
				[]grpc.DialOption{grpc.WithInsecure()},
			)

			if err != nil {
				fmt.Printf("bind endpoint for user failed: %v\n", err)
				return
			}

			err = secret.RegisterSecretServiceHandlerFromEndpoint(ctx, gwmux, rpcProxy,
				[]grpc.DialOption{grpc.WithInsecure()},
			)

			if err != nil {
				fmt.Printf("bind endpoint for secret failed: %v\n", err)
				return
			}

			err = project.RegisterProjectServiceHandlerFromEndpoint(ctx, gwmux, rpcProxy,
				[]grpc.DialOption{grpc.WithInsecure()},
			)

			if err != nil {
				fmt.Printf("bind endpoint for project failed: %v\n", err)
				return
			}

			err = job.RegisterJobServiceHandlerFromEndpoint(ctx, gwmux, rpcProxy,
				[]grpc.DialOption{grpc.WithInsecure()},
			)

			if err != nil {
				fmt.Printf("bind endpoint for job failed: %v\n", err)
				return
			}

			// bind endpoints finished

			log.Infof(fmt.Sprintf("Sterted http:... %s", listenString))

			err = http.ListenAndServe(listenString,
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					// hosting request here
					gwmux.ServeHTTP(w, r)
				}))
			if err != nil {
				panic(err)
			}

			return
		},
	}
	return cmd
}
