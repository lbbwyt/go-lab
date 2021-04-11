package tencent_yun_client

import (
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"math"
	"math/rand"
	"sort"
	"time"
)

const (
	ak = "AKIDqLFvKW8bS1iiXCSz1mfMTOt88NxXbugD"
	sk = "oN9b8PSEeeT5pZQEnanmWK7A4hWCdpyq"
)

const (
	SHA256 = "HmacSHA256"
	SHA1   = "HmacSHA1"
)

type Auth struct {
	SecretId        string //ak
	SecretKey       string //sk
	SignatureMethod string //HmacSHA256 和 HmacSHA1, 默认：HmacSHA256
	Region          string //地域参数，用来标识希望操作哪个地域的实例
	Token           string //临时证书所用的 Token，需要结合临时密钥一起使用。长期密钥不需要 Token。
}

func NewAuth(ak, sk string) *Auth {
	return &Auth{
		SecretId:        ak,
		SecretKey:       sk,
		SignatureMethod: SHA256,
		Token:           "",
	}
}

func (a *Auth) GetNonce() int {
	rand.Seed(time.Now().Unix())
	randNum := rand.Intn(math.MaxInt16)
	return randNum
}

//字典序升序排列参数
func (a *Auth) SortParam(params map[string]string) string {
	if len(params) == 0 {
		return ""
	}
	paramList := make([]string, 0)
	for k, _ := range params {
		paramList = append(paramList, k)
	}
	sort.Strings(paramList)

	res := ""
	for _, k := range paramList {
		v, _ := params[k]
		res = fmt.Sprintf("%s%s=%s&", res, k, v)
	}
	return res[:len(res)-1]
}

func (a *Auth) Sign(string2sign, secretKey, method string) string {
	hashed := hmac.New(sha1.New, []byte(secretKey))
	if method == SHA256 {
		hashed = hmac.New(sha256.New, []byte(secretKey))
	}
	hashed.Write([]byte(string2sign))

	return base64.StdEncoding.EncodeToString(hashed.Sum(nil))
}
