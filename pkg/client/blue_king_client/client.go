package blue_king_client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-lab/pkg/client/blue_king_client/auth"
	"go-lab/pkg/client/blue_king_client/bk_res"
	"io/ioutil"
	"net/http"
	"time"
)

const HostName = "http://bk-paas.cn"

type BlueKingClient interface {
	GetAuth() *auth.Auth
	GetTestAuth() *auth.Auth
}

func DoPost(path string, params interface{}) (*bk_res.BkResult, error) {
	res := new(bk_res.BkResult)
	url := fmt.Sprintf("%s%s", HostName, path)
	client := &http.Client{Timeout: 30 * time.Second}
	jsonStr, _ := json.Marshal(params)
	resp, err := client.Post(url, "application/json", bytes.NewBuffer(jsonStr))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	result, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(result, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}
