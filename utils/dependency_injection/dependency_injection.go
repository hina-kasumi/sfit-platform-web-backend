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
	UserRepo  *repositories.UserRepository
	EventRepo *repositories.EventRepository

	// service
	UserService        *services.UserService
	RedisService       *services.RedisService
	JwtService         *services.JwtService
	RefreshService     *services.RefreshTokenService
	AuthService        *services.AuthService
	EventService       *services.EventService
	UserProfileService *services.UserProfileService

	//handler
	BaseHandler        *handlers.BaseHandler
	AuthHandler        *handlers.AuthHandler
	EventHandler       *handlers.EventHandler
	UserProfileHandler *handlers.UserProfileHandler
	UserHandler        *handlers.UserHandler
}

func NewDI(db *gorm.DB, redisClient *redis.Client, redisCtx context.Context) *DI {
	// Khởi tạo Repository
	userRepo := repositories.NewUserRepository(db)
	eventRepo := repositories.NewEventRepository(db)
	// Khởi tạo Service
	userSer := services.NewUserService(userRepo)
	redisSer := services.NewRedisService(redisClient, redisCtx)
	jwtSer := services.NewJwtService(redisSer)
	refreshSer := services.NewRefreshTokenService()
	authSer := services.NewAuthService(userSer, jwtSer, refreshSer)
	eventSer := services.NewEventService(eventRepo)
	profileSer := services.NewUserProfileService(db, userSer)

	// Khởi tạo Hander
	baseHandler := handlers.NewBaseHandler()
	authHandler := handlers.NewAuthHandler(baseHandler, authSer, jwtSer, refreshSer)
	eventHandler := handlers.NewEventHandler(baseHandler, eventSer)
	profileHandler := handlers.NewUserProfileHandler(profileSer)
	userHandler := handlers.NewUserHandler(userSer)

	return &DI{
		UserRepo:           userRepo,
		EventRepo:          eventRepo,
		UserService:        userSer,
		RedisService:       redisSer,
		JwtService:         jwtSer,
		RefreshService:     refreshSer,
		AuthService:        authSer,
		EventService:       eventSer,
		BaseHandler:        baseHandler,
		AuthHandler:        authHandler,
		EventHandler:       eventHandler,
		UserProfileService: profileSer,
		UserProfileHandler: profileHandler,
		UserHandler:        userHandler,
	}
}
