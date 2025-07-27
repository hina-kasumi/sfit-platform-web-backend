package routes

import (
	"sfit-platform-web-backend/handlers"

	"github.com/gin-gonic/gin"
)

type AuthRoutes struct {
	authHandler *handlers.AuthHandler
}

func NewAuthRoute(authHandler *handlers.AuthHandler) *AuthRoutes {
	return &AuthRoutes{
		authHandler: authHandler,
	}
}

func (authRou *AuthRoutes) RegisterRoutes(router *gin.Engine) {
	authHandler := authRou.authHandler

	task := router.Group("/auth")
	task.POST("/register", authHandler.Register)
	task.POST("/login", authHandler.Login)
	task.POST("/logout", authHandler.Logout)
	task.POST("/refresh", authHandler.RefreshToken)
}
