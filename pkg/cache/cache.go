package cache

import (
	"context"
	"github.com/go-redis/redis/v8"
	"sync"
	"time"
)

//适配器模式， 目前支持本地缓存和redis缓存
type Cache interface {
	Put(key string, value interface{}, expiration time.Duration) error
	Get(key string) (interface{}, bool, error)
	// 接口待完善。。。。
}

func GetRedisCache() Cache {
	cch := &RedisCache{
		conn: redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "",
			DB:       0,
		}),
	}

	cch.ctx = context.Background()

	return cch
}

func GetLocalCache() Cache {
	return &LocalCache{
		m: sync.Map{},
	}
}
