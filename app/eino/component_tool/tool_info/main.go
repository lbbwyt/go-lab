package main

import (
	"fmt"
	"github.com/cloudwego/eino/components/tool/utils"
)

func main() {
	type PostRequest struct {
		URL  string `json:"url" jsonschema:"description=The URL to make the POST request"`
		Body string `json:"body" jsonschema:"description=The body to send in the POST request"`
	}

	param, err := utils.GoStruct2ParamsOneOf[PostRequest]()
	if err != nil {
		panic(err)
	}
	fmt.Println(param)

}
