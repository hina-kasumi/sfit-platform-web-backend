package cmd

import (
	"context"
	"log"
	"sfit-platform-web-backend/entities"
	"sfit-platform-web-backend/middlewares"
	"sfit-platform-web-backend/routes"
	dependencyinjection "sfit-platform-web-backend/utits/dependency_injection"
	"sfit-platform-web-backend/utits/validation"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var depInject *dependencyinjection.DI
var rou []routes.IRoute

func StartServer(db *gorm.DB, redisClient *redis.Client, redisCtx context.Context) {
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
	depInject = dependencyinjection.NewDI(db, redisClient, redisCtx)

	rou = []routes.IRoute{
		routes.NewAuthRoute(depInject.AuthHandler),
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
