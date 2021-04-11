package tencent_yun_client

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

const (
	pageSize  = 20
	domainUrl = "https://cns.api.qcloud.com/v2/index.php"
)

type TencentYunClient struct {
	auth *Auth
}

func NewTencentYunClient(accessId, secretKey string) *TencentYunClient {
	auth := NewAuth(accessId, secretKey)
	return &TencentYunClient{
		auth: auth,
	}
}

// 获取域名
func (t *TencentYunClient) GetDomain(offset int, pageSize int) (*ExportDomain, error) {

	var (
		action    = "DomainList"
		host      = "cns.api.qcloud.com/v2/index.php"
		method    = "GET"
		signature = ""
		nonce     = t.auth.GetNonce()
		timestamp = time.Now().Unix()
	)

	params := map[string]string{
		"Action":          action,
		"Timestamp":       fmt.Sprintf("%d", timestamp),
		"Nonce":           fmt.Sprintf("%d", nonce),
		"SecretId":        t.auth.SecretId,
		"SignatureMethod": t.auth.SignatureMethod,
		"offset":          fmt.Sprintf("%d", offset),
		"length":          fmt.Sprintf("%d", pageSize),
	}
	//对参数排序： 字典需升序排序
	paramsStr := t.auth.SortParam(params)
	//拼接签名原文字符串：请求方法 + 请求主机 +请求路径 + ? + 请求字符串
	singStr := fmt.Sprintf("%s%s?%s", method, host, paramsStr)
	signature = url.QueryEscape(t.auth.Sign(singStr, t.auth.SecretKey, SHA256))

	u, _ := url.Parse(fmt.Sprintf("%s?Action=%s&SecretId=%s&Timestamp=%d&Nonce=%d&Signature=%s&SignatureMethod=%s&offset=%d&length=%d",
		domainUrl, action, t.auth.SecretId, timestamp, nonce, signature, t.auth.SignatureMethod, offset, pageSize))

	resp, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("[TencentYunClient] error when get domain, code: [%d], resp: [%v]", resp.StatusCode, resp))
	}
	result := new(ExportDomain)
	if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
		return nil, errors.New("[TencentYunClient] unMarshal domain error ")
	}
	return result, nil
}

// 获取域名解析记录
func (t *TencentYunClient) GetDomainRecordList(domain string, offset, pageSize int) (*ExportDomainRecord, error) {

	var (
		action    = "RecordList"
		host      = "cns.api.qcloud.com/v2/index.php"
		method    = "GET"
		signature = ""
		nonce     = t.auth.GetNonce()
		timestamp = time.Now().Unix()
	)

	params := map[string]string{
		"Action":          action,
		"Timestamp":       fmt.Sprintf("%d", timestamp),
		"Nonce":           fmt.Sprintf("%d", nonce),
		"SecretId":        t.auth.SecretId,
		"SignatureMethod": t.auth.SignatureMethod,
		"offset":          fmt.Sprintf("%d", offset),
		"length":          fmt.Sprintf("%d", pageSize),
		"domain":          domain,
	}
	//对参数排序： 字典需升序排序
	paramsStr := t.auth.SortParam(params)
	//拼接签名原文字符串：请求方法 + 请求主机 +请求路径 + ? + 请求字符串
	singStr := fmt.Sprintf("%s%s?%s", method, host, paramsStr)
	signature = url.QueryEscape(t.auth.Sign(singStr, t.auth.SecretKey, SHA256))

	u, _ := url.Parse(fmt.Sprintf("%s?Action=%s&SecretId=%s&Timestamp=%d&Nonce=%d&Signature=%s&SignatureMethod=%s&offset=%d&length=%d&domain=%s",
		domainUrl, action, t.auth.SecretId, timestamp, nonce, signature, t.auth.SignatureMethod, offset, pageSize, domain))

	resp, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("[TencentYunClient] error when get domain record, code: [%d], resp: [%v]", resp.StatusCode, resp))
	}
	result := new(ExportDomainRecord)
	if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
		return nil, errors.New("[TencentYunClient] unMarshal domain record error ")
	}
	return result, nil

}
