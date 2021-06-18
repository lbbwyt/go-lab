package cache_middleware

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"time"
)

type RedisStore struct {
	RedisClient *redis.Client
}

func (store *RedisStore) Set(key string, value interface{}, expire time.Duration) error {
	payload, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return store.RedisClient.Set(context.Background(), key, string(payload), expire).Err()
}

func (store *RedisStore) Get(key string, value interface{}) error {
	payload, err := store.RedisClient.Get(context.Background(), key).Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(payload, &value)
}

func (store *RedisStore) Delete(key string) error {
	return store.RedisClient.Del(context.Background(), key).Err()
}
