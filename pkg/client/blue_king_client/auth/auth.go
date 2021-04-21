package auth

type Auth struct {
	BkAppCode   string `json:"bk_app_code"`
	BkAppSecret string `json:"bk_app_secret"`
	BkToken     string `json:"bk_token"`
	BkUserName  string `json:"bk_user_name"`
}

func NewAuthWithToken(appCode, appSecret, token string) *Auth {
	return &Auth{
		BkAppCode:   appCode,
		BkAppSecret: appSecret,
		BkToken:     token,
	}
}

func NewAuthWithWhiteList(appCode, appSecret, userName string) *Auth {
	return &Auth{
		BkAppCode:   appCode,
		BkAppSecret: appSecret,
		BkUserName:  userName,
	}
}

func GetTestAuth() *Auth {
	return NewAuthWithWhiteList("kcmdb", "b6ac1a4b-74ca-4d3b-9a27-be4cf88f8576", "admin")
}
