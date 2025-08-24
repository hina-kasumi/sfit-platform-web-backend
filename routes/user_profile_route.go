package routes

import (
	"sfit-platform-web-backend/entities"
	"sfit-platform-web-backend/handlers"
	"sfit-platform-web-backend/middlewares"

	"github.com/gin-gonic/gin"
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
	group := router.Group("/user-profiles")
	group.GET("/:user_id", userprofileRoute.handler.GetUserProfile)

	group.Use(middlewares.EnforceAuthenticatedMiddleware())
	group.DELETE("/:user_id", middlewares.RequireRoles(string(entities.RoleEnumAdmin)), userprofileRoute.handler.DeleteUser)
	group.PUT("", userprofileRoute.handler.UpdateUserProfile)
	group.POST("", userprofileRoute.handler.CreateUserProfile)
}
