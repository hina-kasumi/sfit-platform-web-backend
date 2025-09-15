package main

import (
	"os"
	"sfit-platform-web-backend/cmd"
	"sfit-platform-web-backend/infrastructures"
)

func main() {
	username := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	host := os.Getenv("DB_HOST")

	//kết nối database, cấu hình database
	db := infrastructures.OpenDbConnection(username, password, dbName, host)

	//kết nối redis
	redisClient, redisCtx := infrastructures.InitRedis(os.Getenv("REDIS_ADDRESS"))

	cmd.StartServer(db, redisClient, redisCtx)
}
