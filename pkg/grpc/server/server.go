package server

import (
	"context"
	"errors"
	"fmt"
	"go-framework/internal/server"
	grpc2 "go-framework/pkg/grpc"
	"go-framework/pkg/grpc/config"
	"go-framework/pkg/grpc/middleware"
	"go-framework/pkg/registry"
	"go-framework/util/helper"
	"google.golang.org/grpc"
	"io"
	"net"
	"os"
)

const GrpcHost = "grpc://127.0.0.1"
const ServerNameSuffix = ".rpc"

var DefaultWriter io.Writer = os.Stdout

type Server struct {
	RpcServer   *grpc.Server
	config      config.Config
	middlewares []middleware.Middleware
}

func NewServer(c interface{}, svc *server.SvcContext, middlewares ...middleware.Middleware) *Server {
	var con config.Config
	err := helper.UnMarshalWithInterface(c, &con)
	if err != nil {
		panic(err)
	}

	var serverOpt []grpc.ServerOption

	for _, mid := range middlewares {
		serverOpt = append(serverOpt, grpc.ChainUnaryInterceptor(mid(svc)))
	}

	s := &Server{
		RpcServer: grpc.NewServer(serverOpt...),
		config:    con,
	}

	return s
}

func (s *Server) Run() error {
	if s.config.Server.Rpc.Mode == "etcd" {
		err := s.registryEtcd()
		if err != nil {
			return errors.New("when connecting to etcd using gRPC, etcd is throwing an error. message: " + err.Error())
		}
	}

	// 监听端口
	lis, err := net.Listen("tcp", s.config.Server.Rpc.Addr)
	if err != nil {
		return err
	}

	fmt.Fprintf(DefaultWriter, "Listening and serving GRPC on %s \n", s.config.Server.Rpc.Addr)

	// 运行grpc服务
	err = s.RpcServer.Serve(lis)

	return err
}

func (s *Server) registryEtcd() error {
	r, err := grpc2.RegistryEtcd(s.config.Etcd)
	if err != nil {
		panic(err)
	}

	var serviceName string
	if s.config.Etcd.Key != "" {
		serviceName = s.config.Etcd.Key
	} else {
		serviceName = s.config.App.Name + ServerNameSuffix
	}

	ins := registry.ServiceInstance{
		Name:      serviceName,
		Endpoints: s.getNodeEndpoints(),
	}

	return r.Register(context.Background(), &ins)
}

func (s *Server) getNodeEndpoints() []string {
	return []string{GrpcHost + s.config.Server.Rpc.Addr}
}
