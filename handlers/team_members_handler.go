package handlers

import (
	"sfit-platform-web-backend/dtos"
	"sfit-platform-web-backend/services"
	"sfit-platform-web-backend/utils/response"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
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

	if !h.canBindJSON(ctx, &req) {
		return
	}

	member, err := h.service.AddMember(req.UserID, req.TeamID, req.Role)
	if h.isError(ctx, err) {
		return
	}

	res := dtos.AddTeamMemberResponse{
		ID:        member.UserID.String() + "_" + member.TeamID.String(),
		CreatedAt: member.CreatedAt.Format(time.RFC3339),
	}
	response.Success(ctx, "add member success", res)
}

func (h *TeamMembersHandler) DeleteMember(ctx *gin.Context) {
	teamID := ctx.Param("team_id")
	if teamID == "" {
		response.Error(ctx, 400, "team_id is required in URL")
		return
	}

	userID := ctx.Param("user_id")
	if userID == "" {
		response.Error(ctx, 400, "user_id is required in URL")
		return
	}

	err := h.service.DeleteMember(userID, teamID)
	if h.isError(ctx, err) {
		return
	}

	response.Success(ctx, "success", gin.H{"message": "Member removed from team successfully"})
}

func (h *TeamMembersHandler) UpdateMemberRole(ctx *gin.Context) {
	var req dtos.UpdateTeamMemberRequest
	if !h.canBindJSON(ctx, &req) {
		return
	}

	err := h.service.UpdateMemberRole(req.UserID, req.TeamID, req.Role)
	if h.isError(ctx, err) {
		return
	}

	response.Success(ctx, "success", gin.H{
		"message":   "Role updated successfully",
		"updatedAt": time.Now().Format(time.RFC3339),
	})
}

func (h *TeamMembersHandler) GetTeamsJoinedByUser(ctx *gin.Context) {
	userID := ctx.Param("user_id")
	if userID == "" {
		response.Error(ctx, 400, "user_id is required")
		return
	}

	teams, err := h.service.GetTeamsJoinedByUser(userID)
	if h.isError(ctx, err) {
		return
	}

	response.Success(ctx, "Get team success", teams)
}

func (h *TeamMembersHandler) GetTeamMembers(ctx *gin.Context) {
	teamID := ctx.Query("teamid")
	pageStr := ctx.Query("page")
	pageSizeStr := ctx.Query("pageSize")

	if teamID == "" {
		response.Error(ctx, 400, "teamid query parameter is required")
		return
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	result, err := h.service.GetMembers(teamID, page, pageSize)
	if err != nil {
		response.Error(ctx, 500, err.Error())
		return
	}

	response.Success(ctx, "Get team member success", result)
}
