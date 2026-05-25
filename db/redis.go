package db

import (
	"context"
	"os"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client
var Ctx = context.Background()

func InitRedis() {
	addr := os.Getenv("REDIS_ADDR")

	if addr == "" {
		addr = "localhost:6379"
	}

	RedisClient = redis.NewClient(&redis.Options{
		Addr: addr,
	})

	// test connection
	_, err := RedisClient.Ping(Ctx).Result()

	if err != nil {
		panic("❌ Redis connection failed.")
	}

	println("✅ Redis connected successfully.")
}
