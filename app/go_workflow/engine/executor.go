package engine

import (
	"context"
	"fmt"
)

type Execution struct {
	Token *Token
	Node  *Node
}
type ExecutionResult struct {
	Outgoing []string
}

type Executor interface {
	Execute(ctx context.Context, ex Execution) (*ExecutionResult, error)
}

// 简单任务节点
type TaskExecutor struct{}

func (TaskExecutor) Execute(ctx context.Context, ex Execution) (*ExecutionResult, error) {
	fmt.Printf("[Task] %s executed\n", ex.Node.Name)
	return &ExecutionResult{Outgoing: ex.Node.To}, nil
}

// 脚本节点
type ScriptExecutor struct{}

func (ScriptExecutor) Execute(ctx context.Context, ex Execution) (*ExecutionResult, error) {
	script := ex.Node.Config["script"]
	fmt.Printf("[Script] %s --> %s\n", ex.Node.Name, script)

	// 这里可以真正执行脚本（govaluate、expr、otto 等）
	return &ExecutionResult{Outgoing: ex.Node.To}, nil
}

// 排他网关
type ExclusiveGatewayExecutor struct{}

func (ExclusiveGatewayExecutor) Execute(ctx context.Context, ex Execution) (*ExecutionResult, error) {
	fmt.Printf("[ExclusiveGateway] %s executed\n", ex.Node.Name)
	expr := ex.Node.Config["expr"]
	v, ok := ex.Token.Variables[expr]
	if !ok {
		return nil, fmt.Errorf("variable %s not found", expr)
	}
	// 简单实现：expr 变量等于哪个出口就走哪个
	var out string
	for _, to := range ex.Node.To {
		if to == fmt.Sprintf("%v", v) {
			out = to
			break
		}
	}
	if out == "" {
		return nil, fmt.Errorf("no path matches value %v", v)
	}
	return &ExecutionResult{Outgoing: []string{out}}, nil
}
