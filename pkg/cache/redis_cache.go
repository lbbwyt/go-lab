package cache

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

type RedisCache struct {
	conn *redis.Client
	ctx  context.Context
}

// Put adds an entry in the cache.
func (rc *RedisCache) Put(key string, value interface{}, expiration time.Duration) error {
	if err := rc.conn.Set(rc.ctx, key, value, expiration); err != nil {
		return err.Err()
	}
	return nil
}

// Get gets an entry from the cache.
func (rc *RedisCache) Get(key string) (interface{}, bool, error) {
	value, err := rc.conn.Get(rc.ctx, key).Result()
	if err != nil {
		if err != redis.Nil {
			return nil, false, err
		}
	}
	return value, false, nil
}
