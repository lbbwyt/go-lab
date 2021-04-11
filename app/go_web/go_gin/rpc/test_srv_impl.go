package rpc

import (
	"context"
	"go-lab/script/protobuf"
)

type TestSrvImpl struct {
}

func NewTestSrvImpl() *TestSrvImpl {
	return &TestSrvImpl{}
}

func (this *TestSrvImpl) GetWebhooks(ct context.Context, req *protobuf.GetWebhooksReq) (*protobuf.GetWebhooksRes, error) {
	return nil, nil
}
