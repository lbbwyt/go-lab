package bk_cmdb_req

import (
	"go-lab/pkg/client/blue_king_client/auth"
	"go-lab/pkg/client/blue_king_client/bk_req"
)

type BatchUpdateInstanceReq struct {
	bk_req.BkCommonReq
	BkObjId string         `json:"bk_obj_id"`
	Update  []*UpdateField `json:"update"`
}

type UpdateField struct {
	Datas  map[string]string `json:"datas"`
	InstID int               `json:"inst_id"`
}

func (a *BatchUpdateInstanceReq) BindAuth(auth *auth.Auth) {
	a.BkUserName = auth.BkUserName
	a.BkAppSecret = auth.BkAppSecret
	a.BkAppCode = auth.BkAppCode
	a.BkToken = auth.BkToken
}
