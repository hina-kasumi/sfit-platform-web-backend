package handlers

import (
	"github.com/gin-gonic/gin"
	"sfit-platform-web-backend/dtos"
	"sfit-platform-web-backend/services"
	"sfit-platform-web-backend/utils/response"
	"time"
)

type TeamHandler struct {
	*BaseHandler
	teamService *services.TeamService
}

func NewTeamHandler(base *BaseHandler, teamService *services.TeamService) *TeamHandler {
	return &TeamHandler{
		BaseHandler: base,
		teamService: teamService,
	}
}

func (h *TeamHandler) CreateTeam(ctx *gin.Context) {
	var req dtos.CreateTeamRequest
	if !h.canBindJSON(ctx, &req) {
		return
	}

	team, err := h.teamService.CreateTeam(req.ID.String(), req.Name, req.Description)
	if h.isError(ctx, err) {
		return
	}

	res := dtos.CreateTeamResponse{
		ID:       team.ID,
		CreateAt: team.CreatedAt.Format(time.RFC3339),
	}
	response.Success(ctx, res)
}
