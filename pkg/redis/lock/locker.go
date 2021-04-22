package lock

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/rs/xid"
	"time"
)

type Locker interface {
	// Lock can lock one
	Lock(key string, expire time.Duration) (looked bool, err error)
	Unlock() error
}

type lock struct {
	cache *redis.Client
	key   string
	// 是否需要释放key
	needUnlock bool
	isFirst    bool
}

func NewLocker(cache *redis.Client) Locker {
	return &lock{
		isFirst:    false,
		cache:      cache,
		key:        "",
		needUnlock: false,
	}
}

// Lock can lock one, key from GetLockKey function
func (l *lock) Lock(key string, expire time.Duration) (locked bool, err error) {
	if l.isFirst {
		return false, fmt.Errorf("repeat lock")
	}
	l.isFirst = true
	l.key = fmt.Sprintf("%s%s", "locker_", key)

	// 不能一样，一样的话，会提示设置成功
	uuid := xid.New().String()
	locked, err = l.cache.SetNX(context.Background(), l.key, uuid, expire).Result()
	// locked sucess , can unlock
	if locked {
		l.needUnlock = true
	}
	return locked, err
}

func (l *lock) Unlock() error {
	// locked sucess , can unlock
	if !l.needUnlock {
		return nil
	}
	return l.cache.Del(context.Background(), l.key).Err()
}
