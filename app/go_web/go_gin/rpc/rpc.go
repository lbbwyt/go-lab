package rpc

import (
	"fmt"
	"go-lab/app/go_web/go_gin/conf"
	"go-lab/script/protobuf"
	"google.golang.org/grpc"
	"net"
)

func Initialize() {
	var port = "9403"
	if conf.GConfig.RpcServer.Port != "" {
		port = conf.GConfig.RpcServer.Port
	}
	addr := fmt.Sprintf("%s%s", ":", port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		panic("failed to listen  " + err.Error())
	}

	grpcServer := grpc.NewServer()

	protobuf.RegisterTestServiceServer(grpcServer, NewTestSrvImpl())

	err = grpcServer.Serve(lis)
	if err != nil {
		panic(fmt.Sprintf("grpc failed to serve: %v", err))
	}
}

func InitializeGrpcGateWay() {
	var port = "9403"
	var httpPort = "9404"
	if conf.GConfig.RpcServer.Port != "" {
		port = conf.GConfig.RpcServer.Port
	}
	if conf.GConfig.RpcHttpServer.Port != "" {
		httpPort = conf.GConfig.RpcHttpServer.Port
	}
	addr := fmt.Sprintf("%s%s", ":", port)
	go RunGrpcHttpServer(httpPort, addr)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		panic("failed to listen  " + err.Error())
	}

	grpcServer := grpc.NewServer()

	protobuf.RegisterOrgServiceServer(grpcServer, NewOrgSrvImpl())
	err = grpcServer.Serve(lis)
	if err != nil {
		panic(fmt.Sprintf("grpc failed to serve: %v", err))
	}
}
