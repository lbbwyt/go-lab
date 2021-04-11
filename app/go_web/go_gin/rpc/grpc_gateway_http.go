package rpc

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"go-lab/script/protobuf"
	"google.golang.org/grpc"
	"net/http"
)

//启动grpc的http服务
func RunGrpcHttpServer(port string, endpoint string) error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := protobuf.RegisterOrgServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	if err != nil {
		return err
	}

	return http.ListenAndServe(fmt.Sprintf(":%s", port), mux)
}
