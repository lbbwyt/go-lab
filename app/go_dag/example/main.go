package main

import (
	"context"
	"fmt"
	"go-lab/app/go_dag/model"
	"go-lab/app/go_dag/scheduler"
	"go-lab/app/go_dag/store"
	"gopkg.in/yaml.v2"
	"os"
)

func main() {
	// 1. 加载流程
	data, _ := os.ReadFile("./app/go_dag/example/workflow.yaml")
	var wf model.Workflow
	_ = yaml.Unmarshal(data, &wf)

	// 2. 注册任务实现
	s := scheduler.NewScheduler(store.NewMemStore())
	s.Register("validateUser", validateUser)
	s.Register("saveUser", saveUser)
	s.Register("sendEmail", sendEmail)

	// 3. 运行
	s.Run(context.Background(), &wf)

}

func validateUser(ctx context.Context, p map[string]interface{}) error {
	// TODO: 参数校验

	fmt.Println("do task validateUser")
	return nil
}
func saveUser(ctx context.Context, p map[string]interface{}) error {
	// TODO: 写库

	fmt.Println("do task saveUser")
	return nil
}
func sendEmail(ctx context.Context, p map[string]interface{}) error {
	// TODO: 发邮件

	fmt.Println("do task sendEmail")
	return nil
}
