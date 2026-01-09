package model

type Task struct {
	ID      string                 `json:"id"`
	Action  string                 `json:"action"` // 业务函数名
	Params  map[string]interface{} `json:"params"`
	Deps    []string               `json:"deps"`    //依赖任务, 只有这些依赖任务全部完成后，才能执行当前任务
	Timeout int                    `json:"timeout"` // 秒
	Retry   int                    `json:"retry"`
}

type Workflow struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Tasks []Task `json:"tasks"`
}
