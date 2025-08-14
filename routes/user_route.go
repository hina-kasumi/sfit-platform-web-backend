package routes

import (
	"sfit-platform-web-backend/handlers"
	"sfit-platform-web-backend/middlewares"

	"github.com/gin-gonic/gin"
)

type UserRoute struct {
	userHandler *handlers.UserHandler
}

func NewUserRoute(userHandler *handlers.UserHandler) *UserRoute {
	return &UserRoute{userHandler: userHandler}
}

func (userRoute *UserRoute) RegisterRoutes(router *gin.Engine) {
	auth := router.Group("/users")
	auth.Use(middlewares.EnforceAuthenticatedMiddleware())
	auth.GET("", userRoute.userHandler.GetUserList)
	auth.PATCH("/:user_id", userRoute.userHandler.UpdateUser)
}
