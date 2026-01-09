package main

import (
	"context"
	"fmt"
	"go-lab/app/go_workflow/engine"
)

func main() {
	store := engine.NewMemStore()
	rt := engine.NewRuntime(store, store)
	// 注册执行器
	rt.RegisterExecutor("task", engine.TaskExecutor{})
	rt.RegisterExecutor("script", engine.ScriptExecutor{})
	rt.RegisterExecutor("exclusive", engine.ExclusiveGatewayExecutor{})

	// 定义：请假流程
	def := &engine.ProcessDef{
		ID:   "leave-v1",
		Name: "请假流程",
		Nodes: []engine.Node{
			{ID: "start", Name: "开始", Type: engine.Task, Executor: "task", To: []string{"apply1"}},
			{ID: "apply1", Name: "填写请假单1", Type: engine.Task, Executor: "task", To: []string{"apply2"}},

			{ID: "apply2", Name: "填写请假单2", Type: engine.Task, Executor: "task", To: []string{"check"}},
			{ID: "check", Name: "主管审批", Type: engine.Exclusive, Executor: "exclusive",
				Config: map[string]string{"expr": "approved"}, To: []string{"end", "modify"}},
			{ID: "modify", Name: "修改单据", Type: engine.Task, Executor: "task", To: []string{"check"}},
			{ID: "end", Name: "结束", Type: engine.Task, Executor: "task", To: []string{}},
		},
	}
	store.Save(def)

	// 启动实例
	inst, err := rt.Start(context.Background(), def.ID, map[string]interface{}{"days": 3})
	if err != nil {
		panic(err)
	}
	fmt.Println("启动实例", inst.ID)

	go func() {
		// 模拟完成节点
		tokens := store.Tokens()
		for _, tok := range tokens {
			fmt.Println("当前节点", tok.NodeID)
		}

		// 主管审批通过
		_ = rt.CompleteTask(context.Background(), tokens[0], map[string]interface{}{"approved": "modify"})

		//_ = rt.CompleteTask(context.Background(), tokens[0], map[string]interface{}{"approved": "end"})

		fmt.Println("流程结束")

	}()

	select {}
}
