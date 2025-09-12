package handlers

import (
	"sfit-platform-web-backend/dtos"
	"sfit-platform-web-backend/entities"
	"sfit-platform-web-backend/middlewares"
	"sfit-platform-web-backend/services"
	"sfit-platform-web-backend/utils/response"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TeamHandler struct {
	*BaseHandler
	teamService    *services.TeamService
	teamMembersSer *services.TeamMembersService
}

func NewTeamHandler(base *BaseHandler, teamService *services.TeamService, teamMembersSer *services.TeamMembersService) *TeamHandler {
	return &TeamHandler{
		BaseHandler:    base,
		teamService:    teamService,
		teamMembersSer: teamMembersSer,
	}
}

func (h *TeamHandler) GetTeamList(ctx *gin.Context) {
	teams, err := h.teamService.GetTeamList()
	if h.isError(ctx, err) {
		return
	}

	response.Success(ctx, "get team list success", teams)
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
	response.Success(ctx, "create team suss", res)
}

func (h *TeamHandler) UpdateTeam(ctx *gin.Context) {
	teamID := ctx.Param("team_id")

	var req dtos.UpdateTeamRequest
	if !h.canBindJSON(ctx, &req) {
		return
	}

	userID := middlewares.GetPrincipal(ctx)
	if !middlewares.HasRole(ctx, string(entities.RoleEnumAdmin)) {
		role, err := h.teamMembersSer.GetRoleUserInTeam(userID, teamID)
		if h.isError(ctx, err) {
			return
		} else if role != string(entities.RoleEnumHead) {
			response.Error(ctx, 403, "You are not allowed to update this team")
			return
		}
	}

	uuidTeamID, err := uuid.Parse(teamID)
	if err != nil {
		response.Error(ctx, 400, "Invalid team ID")
		return
	}

	updatedTeam, err := h.teamService.UpdateTeam(uuidTeamID, req.Name, req.Description)
	if h.isError(ctx, err) {
		return
	}

	res := dtos.UpdateTeamResponse{
		UpdatedAt: updatedTeam.UpdatedAt.Format(time.RFC3339),
	}
	response.Success(ctx, "update success", res)
}

func (h *TeamHandler) DeleteTeam(ctx *gin.Context) {
	teamID := ctx.Param("team_id")
	if h.isNilOrWhiteSpaceWithMessage(ctx, teamID, "team id is required") {
		return
	}

	userID := middlewares.GetPrincipal(ctx)
	if !middlewares.HasRole(ctx, string(entities.RoleEnumAdmin)) {
		role, err := h.teamMembersSer.GetRoleUserInTeam(userID, teamID)
		if h.isError(ctx, err) {
			return
		} else if role != string(entities.RoleEnumHead) {
			response.Error(ctx, 403, "You are not allowed to delete this team")
			return
		}
	}

	err := h.teamService.DeleteTeam(teamID)
	if h.isError(ctx, err) {
		return
	}

	response.Success(ctx, "delete team success", nil)
}
