package routes

import (
	"sfit-platform-web-backend/entities"
	"sfit-platform-web-backend/handlers"
	"sfit-platform-web-backend/middlewares"

	"github.com/gin-gonic/gin"
)

type TeamMembersRoute struct {
	handler *handlers.TeamMembersHandler
}

func NewTeamMembersRoute(handler *handlers.TeamMembersHandler) *TeamMembersRoute {
	return &TeamMembersRoute{handler: handler}
}

func (r *TeamMembersRoute) RegisterRoutes(router *gin.Engine) {
	group := router.Group("/teams")
	group.GET("/:team_id/users", r.handler.GetTeamMembers)

	group.Use(middlewares.EnforceAuthenticatedMiddleware())
	group.Use(middlewares.RequireRoles(
		string(entities.RoleEnumAdmin),
		string(entities.RoleEnumHead),
		string(entities.RoleEnumVice),
	))

	group.PUT("/:team_id/users/:user_id", r.handler.SaveMember)
	// group.POST("/:team_id/users/:user_id", r.handler.AddMember)
	group.DELETE("/:team_id/users/:user_id", r.handler.DeleteMember)

	router.GET("/users/:user_id/teams", r.handler.GetTeamsJoinedByUser)
}
