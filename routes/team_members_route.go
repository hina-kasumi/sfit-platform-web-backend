package routes

import (
	"github.com/gin-gonic/gin"
	"sfit-platform-web-backend/handlers"
)

type TeamMembersRoute struct {
	handler *handlers.TeamMembersHandler
}

func NewTeamMembersRoute(handler *handlers.TeamMembersHandler) *TeamMembersRoute {
	return &TeamMembersRoute{handler: handler}
}

func (r *TeamMembersRoute) RegisterRoutes(router *gin.Engine) {
	group := router.Group("/teams/:team_id/team_members")
	{
		group.POST("/new", r.handler.AddMember)
		group.DELETE("/delete", r.handler.DeleteMember)
		group.PUT("/update", r.handler.UpdateMemberRole)
	}
}
