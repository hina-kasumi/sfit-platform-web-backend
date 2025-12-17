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

	// Connect to database, configure database
	db := infrastructures.OpenDbConnection(username, password, dbName, host)

	// Connect to Redis
	redisClient, redisCtx := infrastructures.InitRedis(os.Getenv("REDIS_ADDRESS"))

	cmd.StartServer(db, redisClient, redisCtx)
}
