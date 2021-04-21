package bk_cmdb_client

import (
	"go-lab/pkg/client/blue_king_client/auth"
)

type BkCmdbClient struct {
}

func NewBkCmdbClient() *BkCmdbClient {
	return &BkCmdbClient{}
}

//临时
func (c *BkCmdbClient) GetAuth() *auth.Auth {
	//todo 临时返回测试使用的应用鉴权信息
	return c.GetTestAuth()
}

func (c *BkCmdbClient) GetTestAuth() *auth.Auth {
	return auth.GetTestAuth()
}
