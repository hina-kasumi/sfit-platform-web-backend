package routes

import (
	"sfit-platform-web-backend/entities"
	"sfit-platform-web-backend/handlers"
	"sfit-platform-web-backend/middlewares"

	"github.com/gin-gonic/gin"
)

type RoleRoutes struct {
	Handler *handlers.RoleHandler
}

func NewRoleRoutes(handler *handlers.RoleHandler) *RoleRoutes {
	return &RoleRoutes{
		Handler: handler,
	}
}

func (r *RoleRoutes) RegisterRoutes(router *gin.Engine) {
	public := router.Group("/users/:user_id/roles")
	public.GET("", r.Handler.GetUserRoles)
	
	admin := router.Group("/users/:user_id/roles")
	admin.Use(middlewares.EnforceAuthenticatedMiddleware())
	admin.Use(middlewares.RequireRoles(string(entities.RoleEnumAdmin)))
	admin.POST("", r.Handler.AddUserRole)
	admin.DELETE("", r.Handler.DeleteUserRole)
}


