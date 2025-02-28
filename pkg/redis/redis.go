package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"os"
	"time"
)

var Client *redis.Client
var Ctx = context.Background()

func InitRedis() {
	Client = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URL"),
		Password: "",
		DB:       0,
	})

	_, err := Client.Ping(Ctx).Result()
	if err != nil {
		panic("failed to connect to redis")
	}
}

func Set(key string, value interface{}, expiration time.Duration) error {
	return Client.Set(Ctx, key, value, expiration).Err()
}

func Get(key string) (string, error) {
	return Client.Get(Ctx, key).Result()
}

func Delete(keys ...string) error {
	return Client.Del(Ctx, keys...).Err()
}