package routes

import (
	"github.com/gin-gonic/gin"
	"sfit-platform-web-backend/handlers"
)

type TeamRoute struct {
	teamHandler *handlers.TeamHandler
}

func NewTeamRoute(teamHandler *handlers.TeamHandler) *TeamRoute {
	return &TeamRoute{teamHandler: teamHandler}
}

func (r *TeamRoute) RegisterRoutes(router *gin.Engine) {
	group := router.Group("/teams")
	{
		group.POST("", r.teamHandler.CreateTeam)
	}
}
