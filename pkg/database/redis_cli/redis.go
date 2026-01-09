package redis_cli

import (
	"context"
	"fmt"
	"github.com/nitishm/go-rejson/v4"
	"github.com/redis/go-redis/v9"
)

type RedisStackClient struct {
	Client *redis.Client
	ctx    context.Context
	Rh     *rejson.Handler
}

func NewRedisStackClient(addr string, password string, db int) *RedisStackClient {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
		Protocol: 2,
		//示例中的连接选项在 Protocol 字段中指定了 RESP2。我们建议您在 go-redis_cli 中使用 RESP2 进行 Redis 查询引擎操作，因为默认 RESP3 的某些响应结构目前不完整，因此您必须在自己的代码中处理“原始”响应。
		//如果您确实想使用 RESP3，连接时应设置 UnstableResp3 选项,您还必须使用 RawResult() 和 RawVal() 方法而不是通常的 Result() 和 Val() 来访问命令结果
		UnstableResp3: false,
	})
	ctx := context.Background()
	rh := rejson.NewReJSONHandler()
	rh.SetGoRedisClientWithContext(ctx, rdb)
	return &RedisStackClient{
		Client: rdb,
		ctx:    ctx,
		Rh:     rh,
	}
}

// 基础命令封装
func (r *RedisStackClient) Set(ctx context.Context, key string, value interface{}) error {
	return r.Client.Set(ctx, key, value, 0).Err()
}

func (r *RedisStackClient) Get(ctx context.Context, key string) (string, error) {
	return r.Client.Get(ctx, key).Result()
}

func (r *RedisStackClient) Do(ctx context.Context, cmds []interface{}) error {
	if err := r.Client.Do(ctx, cmds...).Err(); err != nil {
		return fmt.Errorf("failed to create index: %w", err)
	}
	return nil
}
