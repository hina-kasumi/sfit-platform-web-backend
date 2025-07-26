package infrastructures

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client
var ctx = context.Background()

func InitRedis(addr string) {
	redisClient = redis.NewClient(&redis.Options{
		Addr: addr,
		DB:   0,
	})

	// Kiểm tra kết nối
	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Không thể kết nối Redis: %v", err)
	}

	log.Println("Kết nối Redis thành công")
}

func GetRedisClient() *redis.Client {
	if redisClient == nil {
		log.Fatal("Redis client is not initialized. Call InitRedis first.")
	}
	return redisClient
}

func GetRedisContext() context.Context {
	return ctx
}
