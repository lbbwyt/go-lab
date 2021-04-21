package bk_cmdb_client

import (
	"go-lab/pkg/client/blue_king_client"
	"go-lab/pkg/client/blue_king_client/bk_req"
	"go-lab/pkg/client/blue_king_client/bk_res"
)

//获取全部模型
func (c *BkCmdbClient) GetModels() (*bk_res.BkResult, error) {
	var (
		path = "/api/c/compapi/v2/cc/search_objects/"
	)
	req := new(bk_req.BkCommonReq)
	req.BindAuth(c.GetAuth())
	return blue_king_client.DoPost(path, req)
}
