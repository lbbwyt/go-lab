package bk_req

import "go-lab/pkg/client/blue_king_client/auth"

type BkReq interface {
	BindAuth(auth *auth.Auth)
}

type BkCommonReq struct {
	BkAppCode   string `json:"bk_app_code"`
	BkAppSecret string `json:"bk_app_secret"`
	BkToken     string `json:"bk_token"`
	BkUserName  string `json:"bk_username"`
}

func (a *BkCommonReq) BindAuth(auth *auth.Auth) {
	a.BkUserName = auth.BkUserName
	a.BkAppSecret = auth.BkAppSecret
	a.BkAppCode = auth.BkAppCode
	a.BkToken = auth.BkToken
}
