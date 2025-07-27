package main

import (
	"os"
	"sfit-platform-web-backend/cmd"
	"sfit-platform-web-backend/infrastructures"

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
	redisClient, redisCtx := infrastructures.InitRedis(os.Getenv("REDIS_ADDRESS"))

	cmd.StartServer(db, redisClient, redisCtx)
}
