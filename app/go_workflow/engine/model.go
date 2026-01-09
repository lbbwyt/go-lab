package engine

// 流程定义（Process Definition）
type ProcessDef struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Nodes []Node `json:"nodes"`
}

// 节点
type Node struct {
	ID       string            `json:"id"`
	Name     string            `json:"name"`
	Type     NodeType          `json:"type"`     // task / script / gateway / subprocess
	To       []string          `json:"to"`       // 出口
	Config   map[string]string `json:"config"`   // 节点私有配置
	Executor string            `json:"executor"` // 执行器 Bean 名称
}

type NodeType string

const (
	Task      NodeType = "TASK"
	Script    NodeType = "SCRIPT"
	SubProc   NodeType = "SUBPROC"
	Exclusive NodeType = "EXCLUSIVE_GATEWAY"
)

// 流程实例
type ProcessInst struct {
	ID        string                 `json:"id"`
	DefID     string                 `json:"defId"`
	State     InstState              `json:"state"`
	Variables map[string]interface{} `json:"variables"`
	tokens    map[string]*Token
}

type InstState string

const (
	ACTIVE   InstState = "ACTIVE"
	COMPLETE InstState = "COMPLETE"
	ABORT    InstState = "ABORT"
)

// Token（当前执行指针）
type Token struct {
	ID        string `json:"id"`
	InstID    string `json:"instId"`
	NodeID    string `json:"nodeId"`
	ParentID  string `json:"parentId,omitempty"` // 子流程用
	Variables map[string]interface{}
}
