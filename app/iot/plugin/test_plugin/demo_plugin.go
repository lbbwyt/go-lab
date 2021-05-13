package main

import (
	"fmt"
	"github.com/dullgiulio/pingo"
	"go-lab/pkg/utils/uuid_utils"
	"reflect"
)

type TestPlugin struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func (p *TestPlugin) DeCode(req string, res *string) error {
	fmt.Println(fmt.Sprintf("%s : start decoding", p.Name))
	fmt.Println(reflect.TypeOf(res))
	fmt.Println(req)
	*res = req
	return nil
}

func (p *TestPlugin) EnCode(req string, res *string) error {
	fmt.Println(fmt.Sprintf("%s : start decoding", p.Name))
	return nil
}

func main() {
	var testPlugin = &TestPlugin{
		Id:   uuid_utils.New(),
		Name: "test_plugin",
	}
	pingo.Register(testPlugin)
	pingo.Run()
}
