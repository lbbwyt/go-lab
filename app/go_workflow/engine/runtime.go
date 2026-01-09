package engine

import (
	"context"
	"errors"
	"fmt"
	"go-lab/pkg/utils/uuid_utils"
)

type Runtime struct {
	defStore  DefinitionStore
	instStore InstanceStore
	executors map[string]Executor
	//dispatcher *Dispatcher
}

func NewRuntime(defStore DefinitionStore, instStore InstanceStore) *Runtime {
	return &Runtime{
		defStore:  defStore,
		instStore: instStore,
		executors: make(map[string]Executor),
	}
}

func (rt *Runtime) RegisterExecutor(name string, ex Executor) {
	rt.executors[name] = ex
}

// 基于模板创建实例并启动
func (rt *Runtime) Start(ctx context.Context, defID string, vars map[string]interface{}) (*ProcessInst, error) {
	def, ok := rt.defStore.Get(defID)
	if !ok {
		return nil, errors.New("process definition not found")
	}
	inst := &ProcessInst{
		ID:        newID(),
		DefID:     defID,
		State:     ACTIVE,
		Variables: vars,
	}

	startNode := def.Nodes[0] // 约定第一个节点为开始
	tok := &Token{
		ID:        newID(),
		InstID:    inst.ID,
		NodeID:    startNode.ID,
		Variables: vars,
	}

	if err := rt.instStore.SaveInst(inst); err != nil {
		return nil, err
	}
	if err := rt.instStore.SaveToken(tok); err != nil {
		return nil, err
	}
	rt.dispatch(tok)
	return inst, nil

}

// 分发 token（异步可改为队列）
func (rt *Runtime) dispatch(tok *Token) {
	//go func() {
	_ = rt.CompleteTask(context.Background(), tok, nil)
	//}()
}

// 推进流程
func (rt *Runtime) CompleteTask(ctx context.Context, tok *Token, vars map[string]interface{}) error {

	inst, ok := rt.instStore.GetInst(tok.InstID)
	if !ok {
		return errors.New("instance not found")
	}
	def, ok := rt.defStore.Get(inst.DefID)
	if !ok {
		return errors.New("definition not found")
	}

	// 合并变量
	for k, v := range vars {
		tok.Variables[k] = v
	}
	// 找到当前节点
	var node *Node
	for i := range def.Nodes {
		if def.Nodes[i].ID == tok.NodeID {
			node = &def.Nodes[i]
			break
		}
	}

	if node == nil {
		return errors.New("node not found")
	}
	// 执行
	ex, ok := rt.executors[node.Executor]
	if !ok {
		return errors.New("executor not found")
	}
	out, err := ex.Execute(ctx, Execution{
		Token: tok,
		Node:  node,
	})
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	// 流转
	rt.takeAll(out, tok)
	return nil
}

// 流转
func (rt *Runtime) takeAll(out *ExecutionResult, parent *Token) {
	for _, nextID := range out.Outgoing {
		child := &Token{
			ID:        newID(),
			InstID:    parent.InstID,
			NodeID:    nextID,
			Variables: parent.Variables,
		}
		_ = rt.instStore.SaveToken(child)
		rt.dispatch(child)
	}
	// 父 token 已结束
	_ = rt.instStore.DeleteToken(parent.ID)
}

func newID() string {
	return uuid_utils.New()
}
