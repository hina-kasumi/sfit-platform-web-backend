package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		return
	} // khai báo để đọc từ .env

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// 4. Chạy server
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Không thể khởi động server:", err)
	}
}
