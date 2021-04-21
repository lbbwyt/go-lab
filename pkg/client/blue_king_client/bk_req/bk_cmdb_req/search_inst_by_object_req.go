package bk_cmdb_req

import (
	"go-lab/pkg/client/blue_king_client/auth"
	"go-lab/pkg/client/blue_king_client/bk_req"
)

type SearchInstByObjectReq struct {
	bk_req.BkCommonReq
	BkObjId string `json:"bk_obj_id"`
}

func (a *SearchInstByObjectReq) BindAuth(auth *auth.Auth) {
	a.BkUserName = auth.BkUserName
	a.BkAppSecret = auth.BkAppSecret
	a.BkAppCode = auth.BkAppCode
	a.BkToken = auth.BkToken
}
