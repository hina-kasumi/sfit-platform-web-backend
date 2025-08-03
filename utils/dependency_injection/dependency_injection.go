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
	TagRepo      *repositories.TagRepository
	CourseRepo   *repositories.CourseRepository
	TagTempRepo *repositories.TagTempRepository


	// service
	UserService    *services.UserService
	RedisService   *services.RedisService
	JwtService     *services.JwtService
	RefreshService *services.RefreshTokenService
	AuthService    *services.AuthService
	TagService     *services.TagService
	CourseService  *services.CourseService
	TagTempService *services.TagTempService


	//handler
	BaseHandler *handlers.BaseHandler
	AuthHandler *handlers.AuthHandler
	CourseHandler *handlers.CourseHandler
}

func NewDI(db *gorm.DB, redisClient *redis.Client, redisCtx context.Context) *DI {
	// Khởi tạo Repository
	userRepo := repositories.NewUserRepository(db)
	tagRepo := repositories.NewTagRepository(db)
	tagTempRepo := repositories.NewTagTempRepository(db)
	courseRepo := repositories.NewCourseRepository(db)

	// Khởi tạo Service
	userSer := services.NewUserService(userRepo)
	redisSer := services.NewRedisService(redisClient, redisCtx)
	jwtSer := services.NewJwtService(redisSer)
	refreshSer := services.NewRefreshTokenService()
	authSer := services.NewAuthService(userSer, jwtSer, refreshSer)
	tagSer := services.NewTagService(tagRepo)
	tagTempSer := services.NewTagTempService(tagTempRepo)
	courseSer := services.NewCourseService(courseRepo)

	// Khởi tạo Hander
	baseHandler := handlers.NewBaseHandler()
	authHandler := handlers.NewAuthHandler(baseHandler, authSer, jwtSer, refreshSer)
	courseHandler := handlers.NewCourseHandler(baseHandler, courseSer, tagSer, tagTempSer)

	return &DI{
		UserRepo:       userRepo,
		UserService:    userSer,
		RedisService:   redisSer,
		JwtService:     jwtSer,
		RefreshService: refreshSer,
		AuthService:    authSer,
		BaseHandler:    baseHandler,
		AuthHandler:    authHandler,

		TagRepo:       tagRepo,
		TagService:    tagSer,
		CourseRepo:    courseRepo,
		CourseService: courseSer,
		TagTempRepo:   tagTempRepo,
		TagTempService: tagTempSer,
		CourseHandler:  courseHandler,
	}
}
