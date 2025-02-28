package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/shani34/book-management-system/config"
	"fmt"
	"time"
)

type RedisClient  struct{
  Client *redis.Client
}
var Ctx = context.Background()

func InitRedis() (*RedisClient, error) {
	cfg := config.Get().Redis

	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	_, err := client.Ping(Ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	return &RedisClient{
		Client: client,
	}, nil
}

func(r *RedisClient) Set(key string, value interface{}, expiration time.Duration) error {
	return r.Client.Set(Ctx, key, value, expiration).Err()
}

func(r *RedisClient) Get(key string) (string, error) {
	return r.Client.Get(Ctx, key).Result()
}

func(r *RedisClient) Delete(keys ...string) error {
	return r.Client.Del(Ctx, keys...).Err()
}