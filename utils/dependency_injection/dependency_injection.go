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
	UserRepo        *repositories.UserRepository
	TagRepo         *repositories.TagRepository
	TeamRepo        *repositories.TeamRepository
	TeamMembersRepo *repositories.TeamMembersRepository
	EventRepo       *repositories.EventRepository
	UserProfileRepo *repositories.UserProfileRepository
	// service
	UserService        *services.UserService
	TagService         *services.TagService
	TeamService        *services.TeamService
	TeamMembersService *services.TeamMembersService
	RedisService       *services.RedisService
	JwtService         *services.JwtService
	RefreshService     *services.RefreshTokenService
	AuthService        *services.AuthService
	EventService       *services.EventService
	UserProfileService *services.UserProfileService
	//handler
	BaseHandler        *handlers.BaseHandler
	AuthHandler        *handlers.AuthHandler
	TagHandler         *handlers.TagHandler
	TeamHandler        *handlers.TeamHandler
	TeamMembersHandler *handlers.TeamMembersHandler
	EventHandler       *handlers.EventHandler
	UserProfileHandler *handlers.UserProfileHandler
	UserHandler        *handlers.UserHandler
}

func NewDI(db *gorm.DB, redisClient *redis.Client, redisCtx context.Context) *DI {
	// Khởi tạo Repository
	userRepo := repositories.NewUserRepository(db)
	tagRepo := repositories.NewTagRepository(db)
	teamRepo := repositories.NewTeamRepository(db)
	teamMembersRepo := repositories.NewTeamMembersRepository(db)
	eventRepo := repositories.NewEventRepository(db)
	userProfileRepo := repositories.NewUserProfileRepository(db)

	// Khởi tạo Service
	userSer := services.NewUserService(userRepo)
	tagSer := services.NewTagService(tagRepo)
	teamSer := services.NewTeamService(teamRepo)
	teamMembersService := services.NewTeamMembersService(teamMembersRepo, userRepo)
	redisSer := services.NewRedisService(redisClient, redisCtx)
	jwtSer := services.NewJwtService(redisSer)
	refreshSer := services.NewRefreshTokenService()
	authSer := services.NewAuthService(userSer, jwtSer, refreshSer)
	eventSer := services.NewEventService(eventRepo)
	profileSer := services.NewUserProfileService(userProfileRepo, userSer)

	// Khởi tạo Hander
	baseHandler := handlers.NewBaseHandler()
	authHandler := handlers.NewAuthHandler(baseHandler, authSer, jwtSer, refreshSer)
	tagHandler := handlers.NewTagHandler(baseHandler, tagSer)
	teamHandler := handlers.NewTeamHandler(baseHandler, teamSer)
	teamMembersHandler := handlers.NewTeamMembersHandler(handlers.NewBaseHandler(), teamMembersService)
	eventHandler := handlers.NewEventHandler(baseHandler, eventSer)
	profileHandler := handlers.NewUserProfileHandler(baseHandler, profileSer)
	userHandler := handlers.NewUserHandler(baseHandler, userSer)

	return &DI{
		UserRepo:        userRepo,
		TagRepo:         tagRepo,
		TeamRepo:        teamRepo,
		TeamMembersRepo: teamMembersRepo,
		EventRepo:       eventRepo,
		UserProfileRepo: userProfileRepo,

		UserService:        userSer,
		TagService:         tagSer,
		TeamService:        teamSer,
		TeamMembersService: teamMembersService,
		RedisService:       redisSer,
		JwtService:         jwtSer,
		RefreshService:     refreshSer,
		AuthService:        authSer,
		EventService:       eventSer,
		UserProfileService: profileSer,

		BaseHandler:        baseHandler,
		AuthHandler:        authHandler,
		TagHandler:         tagHandler,
		TeamHandler:        teamHandler,
		TeamMembersHandler: teamMembersHandler,
		EventHandler:       eventHandler,
		UserProfileHandler: profileHandler,
		UserHandler:        userHandler,
	}
}
