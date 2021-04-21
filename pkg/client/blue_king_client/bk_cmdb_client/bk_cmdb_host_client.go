package bk_cmdb_client

import (
	"errors"
	"go-lab/pkg/client/blue_king_client"
	"go-lab/pkg/client/blue_king_client/bk_model/bk_cmdb_model"
	"go-lab/pkg/client/blue_king_client/bk_req/bk_cmdb_req"
	"go-lab/pkg/client/blue_king_client/bk_res"

	"strconv"
)

//新增主机到资源池
func (c *BkCmdbClient) AddHostToResource(hostList []*bk_cmdb_model.HostInfo) (*bk_res.BkResult, error) {
	var (
		path = "/api/c/compapi/v2/cc/add_host_to_resource/"
	)
	if hostList == nil || len(hostList) == 0 {
		return nil, errors.New("invalid params")
	}
	req := new(bk_cmdb_req.AddHostToResourceReq)
	req.BindAuth(c.GetAuth())
	params := make(map[string]*bk_cmdb_model.HostInfo, 0)
	for index, v := range hostList {
		params[strconv.Itoa(index)] = v
	}
	req.HostInfos = params
	return blue_king_client.DoPost(path, req)
}
