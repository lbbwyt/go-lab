package lock

import (
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestLock_Lock(t *testing.T) {

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	lock := NewLocker(client)
	locked, err := lock.Lock("lock1", time.Minute)
	require.NoError(t, err)
	require.Equal(t, true, locked)

	lock = NewLocker(client)
	locked, err = lock.Lock("lock2", time.Minute)
	require.NoError(t, err)
	require.Equal(t, true, locked)

	lock = NewLocker(client)
	locked, err = lock.Lock("lock3", time.Minute)
	require.NoError(t, err)
	require.Equal(t, true, locked)

	lock = NewLocker(client)
	locked, err = lock.Lock("lock1", time.Minute)
	require.NoError(t, err)
	require.Equal(t, false, locked)

}
