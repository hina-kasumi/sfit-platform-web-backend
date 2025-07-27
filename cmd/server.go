package cmd

import (
	"context"
	"sfit-platform-web-backend/routes"
	dependencyinjection "sfit-platform-web-backend/utits/dependency_injection"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var depInject *dependencyinjection.DI
var rou []routes.IRoute

func InitServer(db *gorm.DB, redisClient *redis.Client, redisCtx context.Context) {
	depInject = dependencyinjection.NewDI(db, redisClient, redisCtx)
}

func InitRoutes() {
	rou = []routes.IRoute{
		routes.NewAuthRoute(depInject.AuthHandler),
	}
}

func RegisterRoutes(r *gin.Engine) {
	for _, v := range rou {
		v.RegisterRoutes(r)
	}
}

func GetDepInject() *dependencyinjection.DI {
	return depInject
}
