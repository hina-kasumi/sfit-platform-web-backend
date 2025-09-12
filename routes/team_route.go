package routes

import (
	"sfit-platform-web-backend/entities"
	"sfit-platform-web-backend/handlers"
	"sfit-platform-web-backend/middlewares"

	"github.com/gin-gonic/gin"
)

type TeamRoute struct {
	teamHandler *handlers.TeamHandler
}

func NewTeamRoute(teamHandler *handlers.TeamHandler) *TeamRoute {
	return &TeamRoute{teamHandler: teamHandler}
}

func (r *TeamRoute) RegisterRoutes(router *gin.Engine) {
	group := router.Group("/teams")
	group.GET("", r.teamHandler.GetTeamList)
	group.Use(middlewares.EnforceAuthenticatedMiddleware())
	group.Use(middlewares.RequireRoles(string(entities.RoleEnumAdmin), string(entities.RoleEnumHead)))

	group.POST("", r.teamHandler.CreateTeam)
	group.PUT("/:team_id", r.teamHandler.UpdateTeam)
	group.DELETE("/:team_id", r.teamHandler.DeleteTeam)
}
