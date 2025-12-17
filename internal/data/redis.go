package data

import (
	"context"

	"github.com/redis/go-redis/v9"

	"hmdp-backend/internal/config"
)

// NewRedis 返回Redis客户端
func NewRedis(cfg config.RedisConfig) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})
}

// Ping 健康检查
func Ping(ctx context.Context, client *redis.Client) error {
	return client.Ping(ctx).Err()
}
