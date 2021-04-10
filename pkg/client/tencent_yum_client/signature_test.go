package tencent_yun_client

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestAuth_SortParam(t *testing.T) {
	a := NewAuth("", "")
	params := map[string]string{
		"Nonce":  "11886",
		"Action": "DescribeInstances",
		"Region": "ap-guangzhou",
	}
	fmt.Println(a.SortParam(params))
}

func TestTencentYunClient_GetDomain(t *testing.T) {
	fmt.Println("*******************获取域名**************************")
	client := NewTencentYunClient(ak, sk)
	res, err := client.GetDomain(0, 5)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	str, _ := json.Marshal(res)
	fmt.Println(string(str))

	fmt.Println("*******************获取域名解析记录**************************")
	if res.Data.Domains != nil && len(res.Data.Domains) > 0 {
		for _, item := range res.Data.Domains {
			record, err := client.GetDomainRecordList(item.Name, 0, 5)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			str, _ := json.Marshal(record)
			fmt.Println(string(str))
		}
	}

}
