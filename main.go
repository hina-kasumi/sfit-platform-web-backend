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
		db.AutoMigrate(&entities.Users{})
	}

	// Khởi tạo Redis
	infrastructures.InitRedis(os.Getenv("REDIS_ADDRESS"))

	r := gin.Default()
	r.Use(middlewares.Cors())

	// 1. Khởi tạo các router
	routers.RegisterAuthRoutes(r)

	// 4. Chạy server
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Không thể khởi động server:", err)
	}
}
