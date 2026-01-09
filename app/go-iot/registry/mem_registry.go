package registry

import (
	"errors"
	"time"
)

type MemRegistry struct {
	TTLCache *Cache
}

// NewMemoryRegistry 创建并启动后台 GC
func NewMemRegistry() *MemRegistry {
	return &MemRegistry{TTLCache: NewTTLCache()}
}

const ttl time.Duration = time.Minute * 30

func (m *MemRegistry) Register(deviceID string, s *Session) error {
	m.TTLCache.Set(deviceID, s, ttl)
	return nil
}

func (m *MemRegistry) Unregister(deviceID string) error {
	m.TTLCache.Del(deviceID)
	return nil
}

func (m *MemRegistry) Lookup(deviceID string) (*Session, bool) {
	d, ok := m.TTLCache.Get(deviceID)
	if ok {
		return d.(*Session), true
	}
	return nil, false
}

func (m *MemRegistry) Heartbeat(deviceID string) error {
	if s, ok := m.TTLCache.Get(deviceID); ok {
		ss := s.(*Session)
		m.TTLCache.Set(deviceID, ss, ttl) // 从新设置过期时间
		return nil
	}

	return errors.New("设备未上线")
}
