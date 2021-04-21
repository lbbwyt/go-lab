package bk_cmdb_client

import (
	"encoding/json"
	"fmt"
	"go-lab/pkg/client/blue_king_client/bk_model/bk_cmdb_model"
	"go-lab/pkg/client/blue_king_client/bk_req/bk_cmdb_req"
	"testing"
)

func TestBkCmdbClient_AddHostToResource(t *testing.T) {
	clent := NewBkCmdbClient()
	hostList := make([]*bk_cmdb_model.HostInfo, 0)
	hostList = append(hostList, &bk_cmdb_model.HostInfo{
		BkHostInnerIp: "10.0.0.3",
		BkCloudID:     0,
		ImportFrom:    "3",
	})
	res, err := clent.AddHostToResource(hostList)
	if err != nil {
		panic(err)
	}
	v, _ := json.Marshal(res)
	fmt.Println(string(v))
}

func TestBkCmdbClient_GetModels(t *testing.T) {
	clent := NewBkCmdbClient()
	res, err := clent.GetModels()
	if err != nil {
		panic(err)
	}
	v, _ := json.Marshal(res)
	fmt.Println(string(v))
}

func TestBkCmdbClient_GetInstanceByObjId(t *testing.T) {
	clent := NewBkCmdbClient()
	res, err := clent.GetInstanceByObjId("host")
	if err != nil {
		panic(err)
	}
	v, _ := json.Marshal(res)
	fmt.Println(string(v))
}

func TestBkCmdbClient_UpdateHost(t *testing.T) {
	clent := NewBkCmdbClient()
	req := &bk_cmdb_req.BatchUpdateInstanceReq{
		BkObjId: "host",
		Update: []*bk_cmdb_req.UpdateField{
			{
				Datas: map[string]string{
					"bk_sn": "lbb_test",
				},
				InstID: 139,
			},
		},
	}
	res, err := clent.BatchUpdateInstance(req)
	if err != nil {
		panic(err)
	}
	v, _ := json.Marshal(res)
	fmt.Println(string(v))
}
