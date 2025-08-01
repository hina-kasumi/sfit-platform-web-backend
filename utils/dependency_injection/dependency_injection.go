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
	TagRepo  *repositories.TagRepository
	TeamRepo *repositories.TeamRepository

	// service
	UserService    *services.UserService
	TagService     *services.TagService
	TeamService    *services.TeamService
	RedisService   *services.RedisService
	JwtService     *services.JwtService
	RefreshService *services.RefreshTokenService
	AuthService    *services.AuthService

	//handler
	BaseHandler *handlers.BaseHandler
	AuthHandler *handlers.AuthHandler
	TagHandler  *handlers.TagHandler
	TeamHandler *handlers.TeamHandler
}

func NewDI(db *gorm.DB, redisClient *redis.Client, redisCtx context.Context) *DI {
	// Khởi tạo Repository
	userRepo := repositories.NewUserRepository(db)
	tagRepo := repositories.NewTagRepository(db)
	teamRepo := repositories.NewTeamRepository(db)

	// Khởi tạo Service
	userSer := services.NewUserService(userRepo)
	tagSer := services.NewTagService(tagRepo)
	teamSer := services.NewTeamService(teamRepo)
	redisSer := services.NewRedisService(redisClient, redisCtx)
	jwtSer := services.NewJwtService(redisSer)
	refreshSer := services.NewRefreshTokenService()
	authSer := services.NewAuthService(userSer, jwtSer, refreshSer)

	// Khởi tạo Hander
	baseHandler := handlers.NewBaseHandler()
	authHandler := handlers.NewAuthHandler(baseHandler, authSer, jwtSer, refreshSer)
	tagHandler := handlers.NewTagHandler(baseHandler, tagSer)
	teamHandler := handlers.NewTeamHandler(baseHandler, teamSer)
	return &DI{
		UserRepo:       userRepo,
		TagRepo:        tagRepo,
		TeamRepo:       teamRepo,
		UserService:    userSer,
		TagService:     tagSer,
		TeamService:    teamSer,
		RedisService:   redisSer,
		JwtService:     jwtSer,
		RefreshService: refreshSer,
		AuthService:    authSer,
		BaseHandler:    baseHandler,
		AuthHandler:    authHandler,
		TagHandler:     tagHandler,
		TeamHandler:    teamHandler,
	}
}
