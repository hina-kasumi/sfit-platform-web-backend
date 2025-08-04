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
	TagRepo      *repositories.TagRepository
	CourseRepo   *repositories.CourseRepository
	TagTempRepo *repositories.TagTempRepository

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
	TagService     *services.TagService
	CourseService  *services.CourseService
	TagTempService *services.TagTempService

	EventService       *services.EventService
	UserProfileService *services.UserProfileService
	//handler
	BaseHandler        *handlers.BaseHandler
	AuthHandler        *handlers.AuthHandler
	CourseHandler *handlers.CourseHandler
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
	tagTempRepo := repositories.NewTagTempRepository(db)
	courseRepo := repositories.NewCourseRepository(db)
	lessonRepo := repositories.NewLessonRepository(db)
	userCourseRepo := repositories.NewUserCourseRepository(db)
	userRateRepo := repositories.NewUserRateRepository(db)
	lessonAttendanceRepo := repositories.NewLessonAttendanceRepository(db)
	moduleRepo := repositories.NewModuleRepository(db)
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
	tagSer := services.NewTagService(tagRepo)
	tagTempSer := services.NewTagTempService(tagTempRepo)
	courseSer := services.NewCourseService(userRepo, courseRepo, lessonRepo, tagTempRepo, userCourseRepo, userRateRepo, lessonAttendanceRepo, moduleRepo)
	eventSer := services.NewEventService(eventRepo)
	profileSer := services.NewUserProfileService(userProfileRepo, userSer)

	// Khởi tạo Hander
	baseHandler := handlers.NewBaseHandler()
	authHandler := handlers.NewAuthHandler(baseHandler, authSer, jwtSer, refreshSer)
	courseHandler := handlers.NewCourseHandler(baseHandler, courseSer, tagSer, tagTempSer)
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

		TagRepo:       tagRepo,
		TagService:    tagSer,
		CourseRepo:    courseRepo,
		CourseService: courseSer,
		TagTempRepo:   tagTempRepo,
		TagTempService: tagTempSer,
		CourseHandler:  courseHandler,
		TagHandler:         tagHandler,
		TeamHandler:        teamHandler,
		TeamMembersHandler: teamMembersHandler,
		EventHandler:       eventHandler,
		UserProfileHandler: profileHandler,
		UserHandler:        userHandler,
	}
}
