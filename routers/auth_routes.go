package routers

import (
	authcontroller "sfit-platform-web-backend/controllers/auth_controller"

	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(router *gin.Engine) {
	task := router.Group("/auth")
	task.POST("/register", authcontroller.Register)
	task.POST("/login", authcontroller.Login)
	task.POST("/logout", authcontroller.Logout)
	task.POST("/refresh", authcontroller.RefreshToken)
}
