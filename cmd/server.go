package cmd

import (
	"context"
	"log"
	"sfit-platform-web-backend/entities"
	"sfit-platform-web-backend/middlewares"
	"sfit-platform-web-backend/routes"
	dependencyinjection "sfit-platform-web-backend/utils/dependency_injection"
	"sfit-platform-web-backend/utils/validation"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// khai báo các biến cần thiết cho server
// 'depInject' là một object chứa tất cả các thành phần được inject, giúp gom tất cả lại một chỗ để dễ quản lý và tái sử dụng ở mọi nơi trong ứng dụng.
// 'rou' là một mảng chứa tất cả các route được regis
var depInject *dependencyinjection.DI
var rou []routes.IRoute

func StartServer(db *gorm.DB, redisClient *redis.Client, redisCtx context.Context) {
	// migrate db để sử dụng các bảng, nếu không có thì sẽ tạo ra bảng mới
	if db != nil {
		err := db.AutoMigrate(
			&entities.Log{}, &entities.Device{}, &entities.Users{}, &entities.UserProfile{},
			&entities.Teams{}, &entities.Event{}, &entities.Course{}, &entities.Module{},
			&entities.Task{}, &entities.Lesson{}, &entities.Tag{}, &entities.FavoriteCourse{},
			&entities.EventAttendance{}, &entities.TagTemp{}, &entities.TeamMembers{}, &entities.UserEvent{},
			&entities.UserCourse{}, &entities.LessonAttendance{}, &entities.Newsfeed{}, &entities.UserRate{},
		)
		if err != nil {
			log.Fatalf("AutoMigrate failed: %v", err)
		}
	}
	// khởi tạo DI
	depInject = dependencyinjection.NewDI(db, redisClient, redisCtx)

	rou = []routes.IRoute{
		routes.NewAuthRoute(depInject.AuthHandler),

		routes.NewTagRoute(depInject.TagHandler),
		routes.NewTeamRoute(depInject.TeamHandler),
		routes.NewTeamMembersRoute(depInject.TeamMembersHandler),

		routes.NewEventRoute(depInject.EventHandler),
		routes.NewUserProfileRoute(depInject.UserProfileHandler),
		routes.NewUserRoute(depInject.UserHandler),

	}

	r := gin.Default()
	r.Use(middlewares.Cors())
	r.Use(middlewares.UserLoaderMiddleware(depInject.JwtService))

	RegisterRoutes(r)
	RegisterValidation()

	if err := r.Run(":8080"); err != nil {
		log.Fatal("Không thể khởi động server:", err)
	}
}

func RegisterRoutes(r *gin.Engine) {
	for _, v := range rou {
		v.RegisterRoutes(r)
	}
}
func RegisterValidation() {
	// Đăng ký custom_validate
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("password", validation.PasswordValidator)
	}
}
