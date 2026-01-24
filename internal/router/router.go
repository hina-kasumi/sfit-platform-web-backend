package router

import (
	"context"
	"log"
	"sfit-platform-web-backend/internal/bootstrap"
	"sfit-platform-web-backend/internal/config"
	"sfit-platform-web-backend/internal/handlers"
	middlewares "sfit-platform-web-backend/internal/middleware"
	"sfit-platform-web-backend/internal/model"
	"sfit-platform-web-backend/internal/repositories"
	"sfit-platform-web-backend/internal/services"
	"sfit-platform-web-backend/internal/worker"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func Register(r *gin.Engine, cfg *config.Config, db *gorm.DB, redisClient *redis.Client, ctx context.Context) {
	// declare your handlers and services here
	// repositories
	userRepo := repositories.NewUserRepository(db)
	tagRepo := repositories.NewTagRepository(db)
	tagTempRepo := repositories.NewTagTempRepository(db)
	courseRepo := repositories.NewCourseRepository(db)
	favorCourseRepo := repositories.NewFavoriteCourseRepository(db)
	roleRepo := repositories.NewRoleRepository(cfg.Admin, db)
	lessonRepo := repositories.NewLessonRepository(db)
	userCourseRepo := repositories.NewUserCourseRepository(db)
	userRateRepo := repositories.NewUserRateRepository(db)
	lessonAttendanceRepo := repositories.NewLessonAttendanceRepository(db)
	moduleRepo := repositories.NewModuleRepository(db)
	teamRepo := repositories.NewTeamRepository(db)
	teamMembersRepo := repositories.NewTeamMembersRepository(db)
	eventRepo := repositories.NewEventRepository(db)
	userProfileRepo := repositories.NewUserProfileRepository(db)
	taskRepo := repositories.NewTaskRepository(db)

	// services
	roleSer := services.NewRoleService(roleRepo, userRepo)
	userSer := services.NewUserService(ctx, redisClient, userRepo, userProfileRepo, roleSer)
	tagSer := services.NewTagService(tagRepo)
	teamMembersService := services.NewTeamMembersService(teamMembersRepo, userRepo, roleSer, redisClient, ctx)
	teamSer := services.NewTeamService(teamRepo, teamMembersService, redisClient, ctx)
	jwtSer := services.NewJwtService(&cfg.Jwt, redisClient, ctx)
	refreshSer := services.NewRefreshTokenService(cfg.RefreshToken)
	tagTempSer := services.NewTagTempService(tagTempRepo)
	courseSer := services.NewCourseService(redisClient, ctx, userRepo,
		courseRepo, favorCourseRepo, lessonRepo, tagTempRepo, userCourseRepo, userRateRepo, lessonAttendanceRepo, moduleRepo, userProfileRepo)
	eventSer := services.NewEventService(eventRepo)
	taskSer := services.NewTaskService(taskRepo)
	profileSer := services.NewUserProfileService(userProfileRepo, userSer, eventSer, courseSer, taskSer, redisClient, ctx)
	authSer := services.NewAuthService(userSer, jwtSer, refreshSer, profileSer)
	lessonSer := services.NewLessonService(cfg.Youtube, lessonRepo, courseSer)
	newfeedSer := services.NewNewFeedService(courseSer, userSer, taskSer, eventSer)

	// handlers
	baseHandler := handlers.NewBaseHandler()
	authHandler := handlers.NewAuthHandler(baseHandler, authSer, jwtSer, refreshSer)
	courseHandler := handlers.NewCourseHandler(baseHandler, courseSer, tagSer, tagTempSer)
	tagHandler := handlers.NewTagHandler(baseHandler, tagSer)
	teamHandler := handlers.NewTeamHandler(baseHandler, teamSer, teamMembersService)
	teamMembersHandler := handlers.NewTeamMembersHandler(baseHandler, teamMembersService)
	eventHandler := handlers.NewEventHandler(baseHandler, eventSer, tagSer, tagTempSer)
	profileHandler := handlers.NewUserProfileHandler(baseHandler, profileSer)
	userHandler := handlers.NewUserHandler(baseHandler, userSer)
	taskHandler := handlers.NewTaskHandler(baseHandler, taskSer)
	roleHandler := handlers.NewRoleHandler(baseHandler, roleSer)
	lessonHandler := handlers.NewLessonHandler(baseHandler, lessonSer)
	newFeedHandler := handlers.NewNewFeedHandler(newfeedSer)

	// worker to auto update event status
	worker := worker.NewWorker(eventSer)
	go worker.Start()

	// Register your routes here

	r.Use(middlewares.UserLoaderMiddleware(jwtSer))
	// Auth routes
	authRouteGroup := r.Group("/auth")
	{
		authRouteGroup.POST("/register", authHandler.Register)
		authRouteGroup.POST("/login", authHandler.Login)
		authRouteGroup.POST("/logout", authHandler.Logout)
		authRouteGroup.POST("/refresh", authHandler.RefreshToken)
	}

	// Course routes
	publicCourse := r.Group("")
	publicCourse.Use(middlewares.EnforceAuthenticatedMiddleware())
	{
		publicCourse.GET("/courses", courseHandler.GetListCourse)
		publicCourse.GET("/courses/:course_id", courseHandler.GetCourseDetailByID)
		publicCourse.GET("/courses-v2/:course_id", courseHandler.GetCourseDetailByIDV2)
		publicCourse.GET("/courses-v2", courseHandler.GetListCourseV2)
		publicCourse.POST("courses/:course_id/favourite", courseHandler.MarkCourseAsFavourite)
		publicCourse.DELETE("courses/:course_id/favourite", courseHandler.UnmarkCourseAsFavourite)
		publicCourse.GET("users/:user_id/courses/:course_id/progress", courseHandler.GetUserProgressInCourse)

		publicCourse.GET("/users/:user_id/courses", courseHandler.GetRegisteredCourses)
		publicCourse.GET("/courses/:course_id/lessons", courseHandler.GetCourseLessons)
		publicCourse.POST("/users/:user_id/rate/courses/:course_id", courseHandler.RateCourse)
		publicCourse.PUT("/users/courses", courseHandler.RegisterUserToCourse)
	}

	protectedCourse := r.Group("/courses")
	protectedCourse.Use(middlewares.EnforceAuthenticatedMiddleware())
	protectedCourse.Use(middlewares.RequireRoles(
		string(model.RoleEnumAdmin),
		string(model.RoleEnumHead),
		string(model.RoleEnumVice),
		string(model.RoleEnumTeacher),
	))
	{
		protectedCourse.POST("", courseHandler.CreateCourse)
		protectedCourse.PUT("/:course_id", courseHandler.UpdateCourse)
		protectedCourse.POST("/:course_id/modules", courseHandler.AddModuleToCourse)
		protectedCourse.DELETE("/modules/:module_id", courseHandler.DeleteModuleInCourse)
		protectedCourse.DELETE("/:course_id", courseHandler.DeleteCourse)
		protectedCourse.GET("/:course_id/users", courseHandler.GetRegisteredUsers)
	}

	// event routes
	event := r.Group("/events")
	{
		event.GET("", eventHandler.GetEventList)
		event.GET("/:event_id", eventHandler.GetEventDetail)
		event.GET("/:event_id/users", eventHandler.GetUsersInEvent)
	}
	eventAuth := r.Group("/events")
	{
		eventAuth.Use(middlewares.EnforceAuthenticatedMiddleware())
		eventAuth.Use(middlewares.RequireRoles(
			string(model.RoleEnumAdmin),
			string(model.RoleEnumHead),
			string(model.RoleEnumVice),
		))

		eventAuth.POST("", eventHandler.CreateEvent)
		eventAuth.PUT("/:event_id", eventHandler.UpdateEvent)
		eventAuth.DELETE("/:event_id", eventHandler.DeleteEvent)
	}
	attendEvent := r.Group("/users/:user_id/events/:event_id")
	{
		attendEvent.Use(middlewares.EnforceAuthenticatedMiddleware())
		attendEvent.Use(middlewares.RequireRoles(
			string(model.RoleEnumAdmin),
			string(model.RoleEnumHead),
			string(model.RoleEnumVice),
			string(model.RoleEnumMember),
		))
		attendEvent.PUT("", eventHandler.UpdateStatusUserAttendance)
	}

	// lesson routes
	lessonGroup := r.Group("/modules/:module_id/lessons")
	{
		lessonGroup.Use(middlewares.EnforceAuthenticatedMiddleware())
		lessonGroup.GET("/:lesson_id", lessonHandler.GetLessonByID)
		lessonGroup.Use(middlewares.RequireRoles(
			string(model.RoleEnumAdmin),
			string(model.RoleEnumHead),
			string(model.RoleEnumVice),
			string(model.RoleEnumTeacher),
		))
		lessonGroup.POST("/", lessonHandler.CreateLesson)
		lessonGroup.PUT("/:lesson_id", lessonHandler.UpdateLesson)
		lessonGroup.DELETE("/:lesson_id", lessonHandler.DeleteLesson)
	}
	r.PUT("/users/:user_id/lessons/:lesson_id",
		middlewares.EnforceAuthenticatedMiddleware(),
		lessonHandler.UpdateStatusLessonAttendance,
	)
	r.GET("/lessons/:lesson_id/users", lessonHandler.GetUsersByLessonID)

	// new feed routes
	newfeed := r.Group("/newfeed")
	{
		newfeed.GET("", newFeedHandler.GetNewFeed)
	}

	// role routes
	publicRoleRoute := r.Group("/users/:user_id/roles")
	{
		publicRoleRoute.GET("", roleHandler.GetUserRoles)
	}

	adminRoleRoute := r.Group("/users/:user_id/roles")
	{
		adminRoleRoute.Use(middlewares.EnforceAuthenticatedMiddleware())
		adminRoleRoute.Use(middlewares.RequireRoles(string(model.RoleEnumAdmin)))
		adminRoleRoute.POST("", roleHandler.AddUserRole)
		adminRoleRoute.DELETE("", roleHandler.DeleteUserRole)
	}

	// tag routes
	tagRoute := r.Group("/tags")
	{
		tagRoute.GET("", tagHandler.GetAllTags)
	}

	// team members routes
	teamMemerRoute := r.Group("/teams")
	{
		teamMemerRoute.GET("/:team_id/users", teamMembersHandler.GetTeamMembers)

		teamMemerRoute.Use(middlewares.EnforceAuthenticatedMiddleware())
		teamMemerRoute.Use(middlewares.RequireRoles(
			string(model.RoleEnumAdmin),
			string(model.RoleEnumHead),
			string(model.RoleEnumVice),
		))

		teamMemerRoute.PUT("/:team_id/users/:user_id", teamMembersHandler.SaveMember)
		// group.POST("/:team_id/users/:user_id", teamMembersHandler.AddMember)
		teamMemerRoute.DELETE("/:team_id/users/:user_id", teamMembersHandler.DeleteMember)
	}
	r.GET("/users/:user_id/teams", teamMembersHandler.GetTeamsJoinedByUser)

	// team routes
	teamRoute := r.Group("/teams")
	{
		teamRoute.GET("", teamHandler.GetTeamList)
		teamRoute.Use(middlewares.EnforceAuthenticatedMiddleware())
		teamRoute.Use(middlewares.RequireRoles(string(model.RoleEnumAdmin), string(model.RoleEnumHead)))

		teamRoute.POST("", teamHandler.CreateTeam)
		teamRoute.PUT("/:team_id", teamHandler.UpdateTeam)
		teamRoute.DELETE("/:team_id", teamHandler.DeleteTeam)
	}

	// user profile routes
	userprofileRoute := r.Group("/user-profiles")
	{
		userprofileRoute.GET("/:user_id", profileHandler.GetUserProfile)

		userprofileRoute.Use(middlewares.EnforceAuthenticatedMiddleware())
		userprofileRoute.DELETE("/:user_id", middlewares.RequireRoles(string(model.RoleEnumAdmin)), profileHandler.DeleteUser)
		userprofileRoute.PUT("", profileHandler.UpdateUserProfile)
		userprofileRoute.POST("", profileHandler.CreateUserProfile)
	}

	// user routes
	userRoute := r.Group("/users")
	userRoute.Use(middlewares.EnforceAuthenticatedMiddleware())
	userRoute.GET("", userHandler.GetUserList)
	userRoute.PATCH("/:user_id", userHandler.UpdateUser)

	// task routes
	taskGroup := r.Group("/tasks")
	{
		taskGroup.Use(middlewares.EnforceAuthenticatedMiddleware())
		taskGroup.GET("", taskHandler.ListTasks)
		taskGroup.GET("/:task_id", taskHandler.GetTaskDetail)
		taskGroup.Use(middlewares.RequireRoles(
			string(model.RoleEnumAdmin),
			string(model.RoleEnumHead),
			string(model.RoleEnumVice),
		))
		taskGroup.POST("", taskHandler.CreateTask)
		taskGroup.PUT("/:task_id", taskHandler.UpdateTask)
		taskGroup.DELETE("/:task_id", taskHandler.DeleteTask)
	}

	userTasksGroup := r.Group("/users/:user_id/tasks")
	{
		userTasksGroup.Use(middlewares.EnforceAuthenticatedMiddleware())
		userTasksGroup.GET("", taskHandler.ListTasksByUserID)
		userTasksGroup.PATCH("/:task_id", taskHandler.UpdateTaskUserStatus)
		userTasksGroup.Use(middlewares.RequireRoles(
			string(model.RoleEnumAdmin),
			string(model.RoleEnumHead),
			string(model.RoleEnumVice),
		))
		userTasksGroup.POST("", taskHandler.AddUserTask)
		userTasksGroup.DELETE("/:task_id", taskHandler.DeleteUserTask)
	}

	// Lấy danh sách task theo event_id
	r.GET("events/:event_id/tasks", taskHandler.ListTasksByEventID)

	// Bootstrap
	err := bootstrapApp(cfg, userSer, roleSer)
	if err != nil {
		log.Fatalf("bootstrap admin failed: %v", err)
	}
}

func bootstrapApp(cfg *config.Config, userSer *services.UserService, roleSer *services.RoleService) error {
	err := bootstrap.InitRoles(roleSer)
	if err != nil {
		return err
	}

	err = bootstrap.InitAdmin(cfg.Admin, userSer)
	if err != nil {
		return err
	}
	return nil
}
