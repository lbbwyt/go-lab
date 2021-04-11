package rpc_client

import (
	logger "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type RpcCallClient struct {
	host string
}

func NewRpcCallClient(host string) *RpcCallClient {
	return &RpcCallClient{
		host: host,
	}
}

func (c *RpcCallClient) Dial(task func(conn *grpc.ClientConn)) {
	conn, err := grpc.Dial(c.host, grpc.WithInsecure())
	defer conn.Close()
	if err != nil {
		logger.Error("[grpc Dial] error:" + err.Error())
		return
	}
	task(conn)
}
