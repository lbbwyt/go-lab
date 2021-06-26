package main

import (
	"encoding/xml"
	"io/ioutil"
	"os"
	"sync"
)

//xml 解析器

type KXmlParser struct {
	ExceptionPool *sync.Pool
	fileName      string
}

func NewKXmlParser(fileName string) *KXmlParser {
	return &KXmlParser{
		ExceptionPool: &sync.Pool{New: func() interface{} {
			return &Exception{}
		}},
		fileName: fileName,
	}
}

func (k *KXmlParser) OpenFile(filePath string, v interface{}) error {
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}
	return xml.Unmarshal(data, v)
}
