package routes

import (
	"sfit-platform-web-backend/handlers"

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
	roleRoutes := router.Group("/users/:user_id/roles")

	roleRoutes.POST("", r.Handler.AddUserRole)
	roleRoutes.DELETE("", r.Handler.DeleteUserRole)
}
