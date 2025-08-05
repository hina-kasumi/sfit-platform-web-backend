package routes

import (
	"sfit-platform-web-backend/handlers"

	"github.com/gin-gonic/gin"
)

type TeamMembersRoute struct {
	handler *handlers.TeamMembersHandler
}

func NewTeamMembersRoute(handler *handlers.TeamMembersHandler) *TeamMembersRoute {
	return &TeamMembersRoute{handler: handler}
}

func (r *TeamMembersRoute) RegisterRoutes(router *gin.Engine) {
	group := router.Group("/team-member")
	{
		group.POST("", r.handler.AddMember)
		group.DELETE("/:team_id/:user_id", r.handler.DeleteMember)
		group.PUT("", r.handler.UpdateMemberRole)
	}
	group1 := router.Group("/team")
	{
		group1.GET("/joined/:user_id", r.handler.GetTeamsJoinedByUser)
	}
	router.GET("/team/member", r.handler.GetTeamMembers)

}
