package main

import (
	"fmt"
	"strings"
	"testing"
)

func TestKXmlParser_OpenFile(t *testing.T) {
	path := "/Users/mac/libaobao/github/private/go-lab/cmd/xml_parser/exception.xlsx"
	index := strings.LastIndex(path, "/")
	fmt.Println(index)
	sub := path[0:index]
	fmt.Println(sub)
	index = strings.LastIndex(sub, "/")
	last := path[index:len(sub)]
	fmt.Println(last)
	last = strings.Trim(last, "/")
	fmt.Println(last)

}
