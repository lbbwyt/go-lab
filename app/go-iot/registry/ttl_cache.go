package registry

import (
	"sync"
	"time"
)

type item struct {
	value      interface{}
	expiration int64 // UnixNano
}

type Cache struct {
	mu   sync.RWMutex
	data map[string]item
}

func NewTTLCache() *Cache {
	c := &Cache{data: make(map[string]item)}
	go c.gc()
	return c
}

// Set 写入
func (c *Cache) Set(k string, v interface{}, ttl time.Duration) {
	c.mu.Lock()
	c.data[k] = item{
		value:      v,
		expiration: time.Now().Add(ttl).UnixNano(),
	}
	c.mu.Unlock()
}

// Get 读取
func (c *Cache) Get(k string) (interface{}, bool) {
	c.mu.RLock()
	it, ok := c.data[k]
	c.mu.RUnlock()
	if !ok || time.Now().UnixNano() > it.expiration {
		return nil, false
	}
	return it.value, true
}

// Del 删除
func (c *Cache) Del(k string) {
	c.mu.Lock()
	delete(c.data, k)
	c.mu.Unlock()
}

// 后台 GC
func (c *Cache) gc() {
	for {
		time.Sleep(time.Second)
		now := time.Now().UnixNano()
		c.mu.Lock()
		for k, v := range c.data {
			if now > v.expiration {
				delete(c.data, k)
			}
		}
		c.mu.Unlock()
	}
}
