package infrastructures

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

// về cơ bản nó là một cơ sở dữ liệu key-value, nó được sử dụng để lưu trữ dữ liệu trong trường hợp bạn cần nó để truy cập nhanh hơn.
// InitRedis khởi tạo một Redis client với địa chỉ được truyền vào
func InitRedis(addr string) (*redis.Client, context.Context) {
	//khởi tạo context mặc định
	var ctx = context.Background()
	// Khởi tạo Redis client với địa chỉ được truyền vào
	redisClient := redis.NewClient(&redis.Options{
		Addr: addr, // địa chỉ Redis
		DB:   0,    //sử dụng DB 0, ở đây dùng DB 0 của redis
		//redis gồm 16 DB, mỗi DB có thể lưu trữ 5 triệu dữ liệu, được đánh số từ 0 đến 15
		//tùy chọn này sẽ giúp bạn tối ưu hóa tài nguyên của Redis
		//nếu bạn không cần nó thì bạn có thể bỏ qua
	})
	// Kiểm tra kết nối
	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Không thể kết nối Redis: %v", err)
	}

	log.Println("Kết nối Redis thành công")

	return redisClient, ctx
}
