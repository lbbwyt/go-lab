package utils

import (
	"crypto/rand"
	"encoding/base64"
	"io"
)

//获取随机token 字符串
func GetRandomToken() (string, error) {
	tmp := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, tmp); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(tmp), nil
}
