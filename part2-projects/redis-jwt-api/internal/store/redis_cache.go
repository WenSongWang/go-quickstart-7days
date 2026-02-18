package store

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisCache：go-redis 的极薄封装，便于替换/测试
type RedisCache struct {
	Client *redis.Client
}

func (c *RedisCache) Get(ctx context.Context, key string) (string, error) {
	return c.Client.Get(ctx, key).Result()
}

func (c *RedisCache) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	return c.Client.Set(ctx, key, value, ttl).Err()
}

func (c *RedisCache) Delete(ctx context.Context, key string) error {
	return c.Client.Del(ctx, key).Err()
}

