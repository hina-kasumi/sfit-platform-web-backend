package infrastructures

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

func InitRedis(addr string) (*redis.Client, context.Context) {
	var ctx = context.Background()

	redisClient := redis.NewClient(&redis.Options{
		Addr: addr,
		DB:   0,
	})

	// Kiểm tra kết nối
	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Không thể kết nối Redis: %v", err)
	}

	log.Println("Kết nối Redis thành công")

	return redisClient, ctx
}
