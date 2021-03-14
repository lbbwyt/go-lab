package redislock

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"
	"testing"
	"time"
)

func TestClient_Obtain(t *testing.T) {
	client := redis.NewClient(&redis.Options{
		Network: "tcp",
		Addr:    "127.0.0.1:6379",
	})
	defer client.Close()
	locker := New(client)
	ctx := context.Background()

	lock, err := locker.Obtain(ctx, "my.key", 100*time.Millisecond, nil)
	if err == ErrLockNotHeld {
		log.Info("can not obtion locker")
	} else if err != nil {
		log.Error(err)
	}

	defer lock.Release(ctx)
	log.Info("obtain a lock ")

	time.Sleep(50 * time.Millisecond)
	if ttl, err := lock.TTL(ctx); err != nil {
		log.Error(err)
	} else if ttl > 0 {
		fmt.Println("Yay, I still have my lock!")
	}

	// Extend my lock.
	if err := lock.Refresh(ctx, 100*time.Millisecond, nil); err != nil {
		log.Fatalln(err)
	}
	time.Sleep(100 * time.Millisecond)
	if ttl, err := lock.TTL(ctx); err != nil {
		log.Fatalln(err)
	} else if ttl == 0 {
		fmt.Println("Now, my lock has expired!")
	}
}
