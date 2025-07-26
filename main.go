package main

import (
	"log"
	"os"
	"sfit-platform-web-backend/entities"
	"sfit-platform-web-backend/infrastructures"
	"sfit-platform-web-backend/middlewares"
	"sfit-platform-web-backend/routers"

	"github.com/gin-gonic/gin"
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
			&entities.Log{},
			&entities.Device{},
			&entities.Users{},
			&entities.UserProfile{},
			&entities.Teams{},
			&entities.Event{},
			&entities.Course{},
			&entities.Module{},
			&entities.Task{},
			&entities.Lesson{},
			&entities.Tag{},
			&entities.FavoriteCourse{},
			&entities.EventAttendance{},
			&entities.TagTemp{},
			&entities.TeamMembers{},
			&entities.UserEvent{},
			&entities.UserCourse{},
			&entities.LessonAttendance{},
			&entities.Newsfeed{},
			&entities.UserRate{},
		)

		if err != nil {
			log.Fatalf("AutoMigrate failed: %v", err)
		}
	}

	// Khởi tạo Redis
	infrastructures.InitRedis(os.Getenv("REDIS_ADDRESS"))

	r := gin.Default()
	r.Use(middlewares.Cors())
	r.Use(middlewares.UserLoaderMiddleware())

	// 1. Khởi tạo các router
	routers.RegisterAuthRoutes(r)

	// 4. Chạy server
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Không thể khởi động server:", err)
	}
}
