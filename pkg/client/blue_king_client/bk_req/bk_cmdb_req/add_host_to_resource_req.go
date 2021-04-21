package bk_cmdb_req

import (
	"go-lab/pkg/client/blue_king_client/auth"
	"go-lab/pkg/client/blue_king_client/bk_model/bk_cmdb_model"
	"go-lab/pkg/client/blue_king_client/bk_req"
)

type AddHostToResourceReq struct {
	bk_req.BkCommonReq
	HostInfos map[string]*bk_cmdb_model.HostInfo `json:"host_info"`
}

func (a *AddHostToResourceReq) BindAuth(auth *auth.Auth) {
	a.BkUserName = auth.BkUserName
	a.BkAppSecret = auth.BkAppSecret
	a.BkAppCode = auth.BkAppCode
	a.BkToken = auth.BkToken
}
