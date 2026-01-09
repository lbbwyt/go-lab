package message_bus

import "time"

// Message 统一结构, 用于在消息总线中传递
type Message struct {
	ID      string      `json:"id"` // uuid
	Topic   string      `json:"topic"`
	Payload interface{} `json:"payload"`
	TS      time.Time   `json:"ts"`
}
