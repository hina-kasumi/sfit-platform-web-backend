package handlers

import (
	"github.com/gin-gonic/gin"
	"sfit-platform-web-backend/dtos"
	"sfit-platform-web-backend/services"
	"sfit-platform-web-backend/utils/response"
	"time"
)

type TeamMembersHandler struct {
	*BaseHandler
	service *services.TeamMembersService
}

func NewTeamMembersHandler(base *BaseHandler, srv *services.TeamMembersService) *TeamMembersHandler {
	return &TeamMembersHandler{
		BaseHandler: base,
		service:     srv,
	}
}

func (h *TeamMembersHandler) AddMember(ctx *gin.Context) {
	var req dtos.AddTeamMemberRequest

	teamID := ctx.Param("team_id")
	if teamID == "" {
		response.Error(ctx, 400, "team_id is required in URL")
		return
	}

	if !h.canBindJSON(ctx, &req) {
		return
	}

	member, err := h.service.AddMember(req.UserID, teamID, req.Role)
	if h.isError(ctx, err) {
		return
	}

	res := dtos.AddTeamMemberResponse{
		ID:        member.UserID.String() + "_" + member.TeamID.String(),
		CreatedAt: member.CreatedAt.Format(time.RFC3339),
	}
	response.Success(ctx, res)
}

func (h *TeamMembersHandler) DeleteMember(ctx *gin.Context) {
	var req dtos.DeleteTeamMemberRequest

	teamID := ctx.Param("team_id")
	if teamID == "" {
		response.Error(ctx, 400, "team_id is required in URL")
		return
	}

	if !h.canBindJSON(ctx, &req) {
		return
	}

	err := h.service.DeleteMember(req.UserID, teamID)
	if h.isError(ctx, err) {
		return
	}

	response.Success(ctx, gin.H{"message": "Member removed from team successfully"})
}
