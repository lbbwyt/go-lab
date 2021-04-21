package cache

import (
	"sync"
	"time"
)

type LocalCache struct {
	m sync.Map
}

type Entry struct {
	Expire time.Time
	Value  interface{}
}

func (c *LocalCache) Put(key string, value interface{}, expiration time.Duration) error {
	entry := &Entry{
		Expire: time.Now().Add(expiration),
		Value:  value,
	}
	c.m.Store(key, entry)
	return nil
}

func (c *LocalCache) Get(key string) (value interface{}, expired bool, err error) {
	v, ok := c.m.Load(key)
	e, ok := v.(*Entry)
	if !ok {
		return nil, false, nil
	}
	return e.Value, e.Expire.Before(time.Now()), nil
}
