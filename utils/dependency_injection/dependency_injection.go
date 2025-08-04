package dependencyinjection

import (
	"context"
	"sfit-platform-web-backend/handlers"
	"sfit-platform-web-backend/repositories"
	"sfit-platform-web-backend/services"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type DI struct {
	//repository
	UserRepo *repositories.UserRepository
	CourseRepo *repositories.CourseRepository

	// service
	UserService    *services.UserService
	CourseService   *services.CourseService
	RedisService   *services.RedisService
	JwtService     *services.JwtService
	RefreshService *services.RefreshTokenService
	AuthService    *services.AuthService

	//handler
	BaseHandler *handlers.BaseHandler
	AuthHandler *handlers.AuthHandler
	CourseHandler *handlers.CourseHandler
}

func NewDI(db *gorm.DB, redisClient *redis.Client, redisCtx context.Context) *DI {
	// Khởi tạo Repository
	userRepo := repositories.NewUserRepository(db)
	courseRepo := repositories.NewCourseRepository(db)

	// Khởi tạo Service
	userSer := services.NewUserService(userRepo)
	redisSer := services.NewRedisService(redisClient, redisCtx)
	jwtSer := services.NewJwtService(redisSer)
	refreshSer := services.NewRefreshTokenService()
	authSer := services.NewAuthService(userSer, jwtSer, refreshSer)
	courseSer := services.NewCourseService(courseRepo)

	// Khởi tạo Hander
	baseHandler := handlers.NewBaseHandler()
	authHandler := handlers.NewAuthHandler(baseHandler, authSer, jwtSer, refreshSer)
	courseHandler := handlers.NewCourseHandler(baseHandler, courseSer)

	return &DI{
		UserRepo:       userRepo,
		CourseRepo:     courseRepo,
		UserService:    userSer,
		CourseService:  courseSer,
		RedisService:   redisSer,
		JwtService:     jwtSer,
		RefreshService: refreshSer,
		AuthService:    authSer,
		BaseHandler:    baseHandler,
		AuthHandler:    authHandler,
		CourseHandler:  courseHandler,
	}
}
