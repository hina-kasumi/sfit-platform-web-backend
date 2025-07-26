package redisservice

import (
	"log"
	"sfit-platform-web-backend/infrastructures"
	"time"

	"github.com/redis/go-redis/v9"
)

func SetRedisValue(key string, value any) error {
	ctx := infrastructures.GetRedisContext()
	redisClient := infrastructures.GetRedisClient()

	err := redisClient.Set(ctx, key, value, 0).Err()
	if err != nil {
		log.Println("Can not set value: ", err)
		return err
	}

	return nil
}

func SetRedisExpire(key string, value any, exp int64) error {
	ctx := infrastructures.GetRedisContext()
	redisClient := infrastructures.GetRedisClient()

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

func GetRedisValue(key string) (any, error) {
	ctx := infrastructures.GetRedisContext()
	redisClient := infrastructures.GetRedisClient()

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
