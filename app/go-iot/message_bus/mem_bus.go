package message_bus

import (
	"sync"
	"time"
)

// Bus 内存总线
type MemBus struct {
	mu   sync.RWMutex
	subs map[string][]chan Message // topic -> []channel
}

func NewMemBus() *MemBus {
	return &MemBus{subs: make(map[string][]chan Message)}
}

// Publish 同步发布（无缓冲）
func (b *MemBus) Publish(topic string, payload interface{}) error {
	b.mu.RLock()
	subs := b.subs[topic]
	b.mu.RUnlock()

	msg := Message{
		ID:      Uuid(),
		Topic:   topic,
		Payload: payload,
		TS:      time.Now(),
	}
	for _, ch := range subs {
		select {
		case ch <- msg:
		default:
			// 慢消费者丢弃，避免阻塞
		}
	}

	return nil
}

// Subscribe 返回只读通道, 允许重复订阅，每调一次 Subscribe() 就新建一条 独立 channel，各自收同一份消息； 发布的时候同一条消息被同时写入多个channel
func (b *MemBus) Subscribe(topic string) (interface{}, error) {
	ch := make(chan Message, 64)
	b.mu.Lock()
	b.subs[topic] = append(b.subs[topic], ch)
	b.mu.Unlock()
	return ch, nil
}

func (b *MemBus) Close() error {
	for _, v := range b.subs {
		for _, c := range v {
			close(c)
		}
	}
	return nil
}

func Uuid() string {
	return time.Now().Format("060102150405.000000")
}
