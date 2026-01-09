package model

// Meta 物模型根节点
type Meta struct {
	Schema     string     `json:"schema"`
	ProductKey string     `json:"productKey"`
	Version    string     `json:"version"`
	Properties []Property `json:"properties"`
	Services   []Service  `json:"services"`
	Events     []Event    `json:"events"`
}

// 根据 identifier 快速取定义
func (m *Meta) Property(id string) *Property {
	for i := range m.Properties {
		if m.Properties[i].Identifier == id {
			return &m.Properties[i]
		}
	}
	return nil
}

type Property struct {
	Identifier string   `json:"id"`         // 字段名
	Name       string   `json:"name"`       // 中文名
	AccessMode string   `json:"accessMode"` // r / rw
	Required   bool     `json:"required"`
	DataType   DataType `json:"dataType"` // 类型描述
}

// DataType 统一类型描述
type DataType struct {
	Type  string `json:"type"`            // int | float | bool | text | enum | struct
	Specs Specs  `json:"specs,omitempty"` // 范围/单位/长度/枚举
}

// Specs 类型约束
type Specs struct {
	Min    *float64   `json:"min,omitempty"` // 数值范围
	Max    *float64   `json:"max,omitempty"`
	Unit   string     `json:"unit,omitempty"`   // ℃ / % / lx
	Length int        `json:"length,omitempty"` // text 最大长度
	Values []EnumItem `json:"values,omitempty"` // enum 可选值
	Fields []Property `json:"fields,omitempty"` // struct 子字段
}

// EnumItem 枚举项
type EnumItem struct {
	Value int    `json:"value"` // 内部值
	Desc  string `json:"desc"`  // 显示名
}

type Service struct {
	Identifier string  `json:"id"` // 命令名
	Name       string  `json:"name"`
	CallType   string  `json:"callType"` // sync | async
	Input      []Param `json:"input"`    // 入参列表
	Output     []Param `json:"output"`   // 出参列表
}

// Param 参数 = 精简版 Property
type Param struct {
	Identifier string   `json:"id"`
	Name       string   `json:"name"`
	DataType   DataType `json:"dataType"`
}

type Event struct {
	Identifier string  `json:"id"` // 事件名
	Name       string  `json:"name"`
	Type       string  `json:"type"` // info | alert | fault
	Output     []Param `json:"output"`
}
