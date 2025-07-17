package utils

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type File struct {
	FileName  string `json:"file_name"`
	Content   []byte `json:"content"`
	Base64Str string `json:"base_64_str"`
	DeviceId  string `json:"device_id"`
	At        int64  `json:"at"`
}

func ReadDirRecursive(dirPath string, dataChan chan *File, fileType string) error {

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.Contains(info.Name(), fileType) {
			f := &File{
				FileName: info.Name(),
			}
			if content, err := os.ReadFile(info.Name()); err == nil {
				f.Content = content
				f.Base64Str = base64.StdEncoding.EncodeToString(content)
				dataChan <- f
			}
		}
		return nil
	})
	return err
}

// 将文件转换为 Base64 字符串
func FileToBase64(filePath string) (string, error) {
	// 读取文件内容
	bytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("read file error: %w", err)
	}

	// 转换为 Base64
	return base64.StdEncoding.EncodeToString(bytes), nil
}

// 将文件转换为 Base64 字符串
func CompassFileToBase64(filePath string) (string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer f.Close()
	img, layout, err := image.Decode(f)
	if err != nil {
		return "", err
	}
	// 修改图片的大小
	NewBuf := bytes.Buffer{}
	switch layout {
	case "png":
		err = png.Encode(&NewBuf, img)
	case "jpeg", "jpg":
		err = jpeg.Encode(&NewBuf, img, &jpeg.Options{Quality: 50})
	default:
		return "", errors.New("该图片格式不支持压缩")
	}
	return base64.StdEncoding.EncodeToString(NewBuf.Bytes()), nil
}
