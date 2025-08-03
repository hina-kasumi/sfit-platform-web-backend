package routes

import (
	"github.com/gin-gonic/gin"
	"sfit-platform-web-backend/handlers"
	"sfit-platform-web-backend/middlewares"
)

type UserRoute struct {
	userHandler *handlers.UserHandler
}

func NewUserRoute(userHandler *handlers.UserHandler) *UserRoute {
	return &UserRoute{userHandler: userHandler}
}

func (userRoute *UserRoute) RegisterRoutes(router *gin.Engine) {
	auth := router.Group("/")
	auth.GET("/user", userRoute.userHandler.GetUserList)
	auth.Use(
		middlewares.EnforceAuthenticatedMiddleware(),
	)

	auth.POST("/change-password", userRoute.userHandler.ChangePassword)
}
