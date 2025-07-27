package services

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisService struct {
	redisC *redis.Client
	ctx    context.Context
}

func NewRedisService(redisC *redis.Client, ctx context.Context) *RedisService {
	return &RedisService{
		redisC: redisC,
		ctx:    ctx,
	}
}

func (redis_ser *RedisService) SetRedisValue(key string, value any) error {
	ctx := redis_ser.ctx
	redisClient := redis_ser.redisC

	err := redisClient.Set(ctx, key, value, 0).Err()
	if err != nil {
		log.Println("Can not set value: ", err)
		return err
	}

	return nil
}

func (redis_ser *RedisService) SetRedisExpire(key string, value any, exp int64) error {
	ctx := redis_ser.ctx
	redisClient := redis_ser.redisC

	// Set key trước nếu chưa có
	duration := time.Duration(exp-time.Now().Unix()) * time.Second
	err := redisClient.Set(ctx, key, value, duration).Err()
	if err != nil {
		log.Println("Can not set key:", err)
		return err
	}

	// Thiết lập thời gian hết hạn tại thời điểm cụ thể
	log.Println("Sẽ hết hạn sau:", duration)

	return nil
}

func (redis_ser *RedisService) GetRedisValue(key string) (any, error) {
	ctx := redis_ser.ctx
	redisClient := redis_ser.redisC

	val, err := redisClient.Get(ctx, key).Result()

	if err == redis.Nil {
		log.Println("Key không tồn tại")
	} else if err != nil {
		log.Println("Lỗi khi lấy dữ liệu:", err)
	} else {
		log.Println("Giá trị của key là:", val)
	}
	return val, err
}
