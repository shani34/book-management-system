package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"os"
	"time"
)

type RedisClient  struct{
  Client *redis.Client
}
var Ctx = context.Background()

func InitRedis()(*RedisClient, error) {
	Client:= redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URL"),
		Password: "",
		DB:       0,
	})

	_, err := Client.Ping(Ctx).Result()
	if err != nil {
		panic("failed to connect to redis")
	}

	return &RedisClient{
		Client: Client,
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