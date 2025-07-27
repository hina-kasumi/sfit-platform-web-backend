package main

import (
	"log"
	"os"
	"sfit-platform-web-backend/cmd"
	"sfit-platform-web-backend/entities"
	"sfit-platform-web-backend/infrastructures"
	"sfit-platform-web-backend/middlewares"
	"sfit-platform-web-backend/utits/validation"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		return
	} // khai báo để đọc từ .env

	username := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	host := os.Getenv("DB_HOST")

	db := infrastructures.OpenDbConnection(username, password, dbName, host)
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

	// Khởi tạo Redis
	redisClient, redisCtx := infrastructures.InitRedis(os.Getenv("REDIS_ADDRESS"))

	cmd.InitServer(db, redisClient, redisCtx)

	//Khởi tạo routes
	cmd.InitRoutes()

	// Cài đặt gin
	r := gin.Default()
	r.Use(middlewares.Cors())
	r.Use(middlewares.UserLoaderMiddleware(cmd.GetDepInject().JwtService))

	// Đăng ký custom_validate
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		valid := validation.NewCustomValidator()
		v.RegisterValidation("password", valid.PasswordValidator)
	}

	// Khởi tạo các router
	cmd.RegisterRoutes(r)

	// Chạy server
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Không thể khởi động server:", err)
	}
}
