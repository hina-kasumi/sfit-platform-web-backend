package services

import (
	"context"
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
		return err
	}

	return nil
}

func (redis_ser *RedisService) GetRedisValue(key string) (any, error) {
	ctx := redis_ser.ctx
	redisClient := redis_ser.redisC

	val, err := redisClient.Get(ctx, key).Result()

	return val, err
}
