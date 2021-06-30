package main

import (
	"fmt"
	"github.com/golang/protobuf/ptypes"
	"go-lab/script/protobuf"
)

func main() {
	test := &protobuf.TestAny{
		Id:      123,
		Title:   "lbb",
		Content: "wyt",
	}
	any, err := ptypes.MarshalAny(test)
	fmt.Println(any, err) // [type.googleapis.com/rpc.TestAny]:{Id:1 Title:"标题" Content:"内容"} <nil>

	msg := &protobuf.Response{
		Code: 0,
		Msg:  "success",
		Data: any,
	}
	fmt.Println(msg) // Msg:"success" data:{[type.googleapis.com/rpc.TestAny]:{Id:1 Title:"标题" Content:"内容"}}

	unmarshal := &protobuf.TestAny{}
	err = ptypes.UnmarshalAny(msg.Data, unmarshal)
	fmt.Println(unmarshal, err) // Id:1 Title:"标题" Content:"内容" <nil>
}
