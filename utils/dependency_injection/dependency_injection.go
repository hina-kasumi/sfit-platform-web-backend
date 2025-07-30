package dependencyinjection

import (
	"context"
	"sfit-platform-web-backend/handlers"
	"sfit-platform-web-backend/repositories"
	"sfit-platform-web-backend/services"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// Struct 'DI'sẽ chứa tất cả các thành phần được inject,
// giúp gom tất cả lại một chỗ để dễ quản lý và tái sử dụng ở mọi nơi trong ứng dụng.
// ví dụ: nếu bạn muốn sử dụng một service trong nhiều handler, bạn chỉ cần inject service vào handler đó thôi.
type DI struct {
	//repository
	UserRepo *repositories.UserRepository

	// service
	UserService    *services.UserService
	RedisService   *services.RedisService
	JwtService     *services.JwtService
	RefreshService *services.RefreshTokenService
	AuthService    *services.AuthService

	//handler
	BaseHandler  *handlers.BaseHandler
	AuthHandler  *handlers.AuthHandler
	EventHandler *handlers.EventHandler
}

func NewDI(db *gorm.DB, redisClient *redis.Client, redisCtx context.Context) *DI {
	// Khởi tạo Repository
	userRepo := repositories.NewUserRepository(db)

	// Khởi tạo Service
	userSer := services.NewUserService(userRepo)
	redisSer := services.NewRedisService(redisClient, redisCtx)
	jwtSer := services.NewJwtService(redisSer)
	refreshSer := services.NewRefreshTokenService()
	authSer := services.NewAuthService(userSer, jwtSer, refreshSer)

	// Khởi tạo Hander
	baseHandler := handlers.NewBaseHandler()
	authHandler := handlers.NewAuthHandler(baseHandler, authSer, jwtSer, refreshSer)

	return &DI{
		UserRepo:       userRepo,
		UserService:    userSer,
		RedisService:   redisSer,
		JwtService:     jwtSer,
		RefreshService: refreshSer,
		AuthService:    authSer,
		BaseHandler:    baseHandler,
		AuthHandler:    authHandler,
	}
}
