package routes

import (
	"github.com/gin-gonic/gin"
	"sfit-platform-web-backend/handlers"
	"sfit-platform-web-backend/middlewares"
)

type UserProfileRoute struct {
	handler *handlers.UserProfileHandler
}

func NewUserProfileRoute(handler *handlers.UserProfileHandler) *UserProfileRoute {
	return &UserProfileRoute{
		handler: handler,
	}
}

func (userprofileRoute *UserProfileRoute) RegisterRoutes(router *gin.Engine) {
	group := router.Group("/user-profile")
	group.DELETE("/:user_id", userprofileRoute.handler.DeleteUser)
	group.Use(middlewares.EnforceAuthenticatedMiddleware())
	group.PUT("/update", userprofileRoute.handler.UpdateUserProfile)

}
