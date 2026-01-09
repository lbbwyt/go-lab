package gateway

import (
	"context"
)

type Transport interface {
	Listen(addr string) error
	SetHandler(func(ctx context.Context, msg *RawMessage))
	Close() error
}

// RawMessage 网关上行的消息
type RawMessage struct {
	Network string // "mqtt", "tcp"
	Product string // 产品 Key
	Device  string // 设备名
	Payload []byte
}
