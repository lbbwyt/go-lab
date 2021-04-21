package bk_cmdb_client

import (
	"errors"
	"go-lab/pkg/client/blue_king_client"
	"go-lab/pkg/client/blue_king_client/bk_req/bk_cmdb_req"
	"go-lab/pkg/client/blue_king_client/bk_res"
)

//获取模型下实例详情, bk中模型为object
func (c *BkCmdbClient) GetInstanceByObjId(ObjectId string) (*bk_res.BkResult, error) {
	var (
		path = "/api/c/compapi/v2/cc/search_inst_by_object/"
	)
	if ObjectId == "" {
		return nil, errors.New("invalid params")
	}
	req := new(bk_cmdb_req.SearchInstByObjectReq)
	req.BindAuth(c.GetAuth())
	req.BkObjId = ObjectId
	return blue_king_client.DoPost(path, req)
}

//主机批量更新，蓝鲸的主机实例字段信息无法通过新增接口导入， 需在导入主机之后调用该接口
func (c *BkCmdbClient) BatchUpdateInstance(req *bk_cmdb_req.BatchUpdateInstanceReq) (*bk_res.BkResult, error) {
	var (
		path = "/api/c/compapi/v2/cc/batch_update_inst/"
	)
	if req == nil || req.BkObjId == "" {
		return nil, errors.New("invalid params")
	}
	req.BindAuth(c.GetAuth())
	return blue_king_client.DoPost(path, req)
}
