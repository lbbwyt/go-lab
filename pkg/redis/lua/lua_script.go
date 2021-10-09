package lua

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

// Redis Represents a Redis counter with minimum threshold and TTL support.
type Redis struct {
	client *redis.Client
	script *redis.Script
}

// New Creates a new Redis counter.
func New(host string, db int) *Redis {
	return &Redis{
		client: redis.NewClient(&redis.Options{
			Addr: host,
			DB:   db,
		}),

		//计数器用于将给定key减去参数值，如果给定可以小于参数指定的最小值，删除该key

		script: redis.NewScript(`
			local min=tonumber(ARGV[2])
			local val=redis.call("DECRBY", KEYS[1], ARGV[1]) 
			if val <= min then redis.call("del", KEYS[1]) 
				return min 
			else 
				redis.call("EXPIRE", KEYS[1], ARGV[3]) 
				return val 
			end`),
	}
}

// DecrToMin Decreases the counter by `value` and removes it once `counter value <= min`.
//
// Returns `min` if the counter does not exist.
func (r *Redis) DecrToMin(ctx context.Context, key string, value, min int64, ttl time.Duration) (int64, error) {
	result, err := r.script.Run(ctx, r.client, []string{key}, value, min, ttl.Seconds()).Int64()
	if err != nil {
		return 0, err
	}
	return result, nil
}

// IncrBy Increments the counter by the specified value and sets the TTL of the key.
func (r *Redis) IncrBy(ctx context.Context, key string, value int64, ttl time.Duration) (int64, error) {
	pipe := r.client.WithContext(ctx).TxPipeline()
	incResult := pipe.IncrBy(ctx, key, value)
	pipe.Expire(ctx, key, ttl)
	_, err := pipe.Exec(ctx)
	return incResult.Val(), err
}
