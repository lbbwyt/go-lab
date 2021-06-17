package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"sync"
	"sync/atomic"
)

// Options controls the behavior of cache.
type Options struct {
	DeleteHandler func(key string, oldData interface{})

	ErrLogFunc func(str string)
}

type AsyncCache struct {
	opt  Options
	data sync.Map
}

type entry struct {
	val    atomic.Value
	expire int32 // 0 means useful, 1 will expire
}

func NewAsyncCache(opt Options) *AsyncCache {
	c := &AsyncCache{
		opt:  opt,
		data: sync.Map{},
	}
	if c.opt.ErrLogFunc == nil {
		c.opt.ErrLogFunc = func(str string) {
			log.Error(str)
		}
	}
	if c.opt.DeleteHandler == nil {
		c.opt.DeleteHandler = func(key string, oldData interface{}) {
			log.Info(fmt.Sprintf("delete key %s", key))
		}
	}
	return c
}

func (c *AsyncCache) Dump() map[string]interface{} {
	data := make(map[string]interface{})
	c.data.Range(func(key, value interface{}) bool {
		k, ok := key.(string)
		if !ok {
			c.opt.ErrLogFunc(fmt.Sprintf("invalid key: %v, type: %T is not string", k, k))
			c.data.Delete(key)
			return true
		}
		data[k] = value
		return true
	})
	return data
}

//shouldDelete 方法用户判断是否可删除的
// 从缓存中删除满足条件的的数据
func (c *AsyncCache) Delete(shouldDelete func(key string) bool) {
	c.data.Range(func(key, value interface{}) bool {
		s := key.(string)
		if shouldDelete(s) {
			if c.opt.DeleteHandler != nil {
				c.opt.DeleteHandler(s, value)
			}
			c.data.Delete(key)
		}
		return true
	})
}

func (c *AsyncCache) Get(key string) (val interface{}, ok bool) {
	return c.data.Load(key)
}

func (c *AsyncCache) Set(key string, val interface{}) {
	c.data.Store(key, val)
}

func main() {
	opt := new(Options)
	c := NewAsyncCache(*opt)
	c.Set("lbb", "wyt")
	c.Set("wyt", "sdf")
	v, ok := c.Get("lbb")
	if ok {
		fmt.Println(v.(string))
	}

	c.Delete(func(key string) bool {
		if key == "lbb" {
			return true
		}
		return false
	})

	v, ok = c.Get("lbb")
	if !ok {
		fmt.Println("deleted")
	}
}
