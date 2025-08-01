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

	team, err := h.teamService.CreateTeam(req.Name, req.Description)
	if h.isError(ctx, err) {
		return
	}

	res := dtos.CreateTeamResponse{
		ID:       team.ID,
		CreateAt: team.CreatedAt.Format(time.RFC3339),
	}
	response.Success(ctx, res)
}

func (h *TeamHandler) UpdateTeam(ctx *gin.Context) {
	var req dtos.UpdateTeamRequest
	if !h.canBindJSON(ctx, &req) {
		return
	}

	updatedTeam, err := h.teamService.UpdateTeam(req.ID, req.Name, req.Description)
	if h.isError(ctx, err) {
		return
	}

	res := dtos.UpdateTeamResponse{
		UpdatedAt: updatedTeam.UpdatedAt.Format(time.RFC3339),
	}
	response.Success(ctx, res)
}
