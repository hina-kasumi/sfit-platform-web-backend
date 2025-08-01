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

	// service
	UserService    *services.UserService
	RedisService   *services.RedisService
	JwtService     *services.JwtService
	RefreshService *services.RefreshTokenService
	AuthService    *services.AuthService

	//handler
	BaseHandler *handlers.BaseHandler
	AuthHandler *handlers.AuthHandler

		// handler Tag, Course
	TagHandler   *handlers.TagHandler
	CourseHandler *handlers.CourseHandler
}

func NewDI(db *gorm.DB, redisClient *redis.Client, redisCtx context.Context) *DI {
	// Khởi tạo Repository
	userRepo := repositories.NewUserRepository(db)

		// Khởi tạo Repository Tag, Course
	tagRepo := repositories.NewTagRepository(db)
	courseRepo := repositories.NewCourseRepository(db)

	// Khởi tạo Service
	userSer := services.NewUserService(userRepo)
	redisSer := services.NewRedisService(redisClient, redisCtx)
	jwtSer := services.NewJwtService(redisSer)
	refreshSer := services.NewRefreshTokenService()
	authSer := services.NewAuthService(userSer, jwtSer, refreshSer)
	
		// Khởi tạo Service Tag, Course
	tagSer := services.NewTagService(tagRepo)
	courseSer := services.NewCourseService(courseRepo)

	// Khởi tạo Hander
	baseHandler := handlers.NewBaseHandler()
	authHandler := handlers.NewAuthHandler(baseHandler, authSer, jwtSer, refreshSer)

		// Khởi tạo Handler Tag, Course
	courseHandler := handlers.NewCourseHandler(baseHandler, courseSer, tagSer)
	tagHandler := handlers.NewTagHandler(baseHandler, tagSer)

	return &DI{
		UserRepo:       userRepo,
		UserService:    userSer,
		RedisService:   redisSer,
		JwtService:     jwtSer,
		RefreshService: refreshSer,
		AuthService:    authSer,
		BaseHandler:    baseHandler,
		AuthHandler:    authHandler,

			// Handler Tag, Course
		TagHandler: tagHandler,
		CourseHandler:  courseHandler,
	}
}
