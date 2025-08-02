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
	UserRepo        *repositories.UserRepository
	TagRepo         *repositories.TagRepository
	TeamRepo        *repositories.TeamRepository
	TeamMembersRepo *repositories.TeamMembersRepository

	// service
	UserService        *services.UserService
	TagService         *services.TagService
	TeamService        *services.TeamService
	TeamMembersService *services.TeamMembersService
	RedisService       *services.RedisService
	JwtService         *services.JwtService
	RefreshService     *services.RefreshTokenService
	AuthService        *services.AuthService

	//handler
	BaseHandler        *handlers.BaseHandler
	AuthHandler        *handlers.AuthHandler
	TagHandler         *handlers.TagHandler
	TeamHandler        *handlers.TeamHandler
	TeamMembersHandler *handlers.TeamMembersHandler
}

func NewDI(db *gorm.DB, redisClient *redis.Client, redisCtx context.Context) *DI {
	// Khởi tạo Repository
	userRepo := repositories.NewUserRepository(db)
	tagRepo := repositories.NewTagRepository(db)
	teamRepo := repositories.NewTeamRepository(db)
	teamMembersRepo := repositories.NewTeamMembersRepository(db)

	// Khởi tạo Service
	userSer := services.NewUserService(userRepo)
	tagSer := services.NewTagService(tagRepo)
	teamSer := services.NewTeamService(teamRepo)
	teamMembersService := services.NewTeamMembersService(teamMembersRepo, userRepo)
	redisSer := services.NewRedisService(redisClient, redisCtx)
	jwtSer := services.NewJwtService(redisSer)
	refreshSer := services.NewRefreshTokenService()
	authSer := services.NewAuthService(userSer, jwtSer, refreshSer)

	// Khởi tạo Hander
	baseHandler := handlers.NewBaseHandler()
	authHandler := handlers.NewAuthHandler(baseHandler, authSer, jwtSer, refreshSer)
	tagHandler := handlers.NewTagHandler(baseHandler, tagSer)
	teamHandler := handlers.NewTeamHandler(baseHandler, teamSer)
	teamMembersHandler := handlers.NewTeamMembersHandler(handlers.NewBaseHandler(), teamMembersService)

	return &DI{
		UserRepo:           userRepo,
		TagRepo:            tagRepo,
		TeamRepo:           teamRepo,
		TeamMembersRepo:    teamMembersRepo,
		UserService:        userSer,
		TagService:         tagSer,
		TeamService:        teamSer,
		TeamMembersService: teamMembersService,
		RedisService:       redisSer,
		JwtService:         jwtSer,
		RefreshService:     refreshSer,
		AuthService:        authSer,
		BaseHandler:        baseHandler,
		AuthHandler:        authHandler,
		TagHandler:         tagHandler,
		TeamHandler:        teamHandler,
		TeamMembersHandler: teamMembersHandler,
	}
}
