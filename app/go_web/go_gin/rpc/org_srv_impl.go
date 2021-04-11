package rpc

import (
	"context"
	"fmt"
	"go-lab/script/protobuf"
)

//组织服务rpc接口
type OrgSrvImpl struct {
}

func NewOrgSrvImpl() *OrgSrvImpl {
	return &OrgSrvImpl{}
}

func (this *OrgSrvImpl) GetOrgs(ctx context.Context, req *protobuf.GetOrgsReq) (*protobuf.GetOrgsRes, error) {
	res := new(protobuf.GetOrgsRes)
	fmt.Println(".......")
	return res, nil
}
