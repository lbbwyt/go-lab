package cache_middleware

import (
	"errors"
	"time"
)

var ErrCacheMiss = errors.New("persist cache miss error")

type CacheStore interface {
	Get(key string, value interface{}) error
	Set(key string, value interface{}, expire time.Duration) error
	Delete(key string) error
}
