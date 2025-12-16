package data

import (
	"context"

	"github.com/redis/go-redis/v9"

	"hmdp-backend/internal/config"
)

// NewRedis returns a configured go-redis client.
func NewRedis(cfg config.RedisConfig) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})
}

// Ping ensures the redis connection is healthy.
func Ping(ctx context.Context, client *redis.Client) error {
	return client.Ping(ctx).Err()
}
