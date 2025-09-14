package handlers

import (
	"sfit-platform-web-backend/services"
	"sfit-platform-web-backend/utils/response"

	"github.com/gin-gonic/gin"
)

type RoleHandler struct {
	*BaseHandler
	roleService *services.RoleService
}

func NewRoleHandler(baseHandler *BaseHandler, roleService *services.RoleService) *RoleHandler {
	return &RoleHandler{
		BaseHandler: baseHandler,
		roleService: roleService,
	}
}

func (rh *RoleHandler) DeleteUserRole(ctx *gin.Context) {
	userID := ctx.Param("user_id")

	var roles []string
	if !rh.canBindJSON(ctx, &roles) {
		return
	}

	err := rh.roleService.RemoveUserRole(userID, roles...)
	if rh.isErrorWithMessage(ctx, err, 500, "Failed to remove user roles") {
		return
	}

	response.Success(ctx, "User roles removed successfully", nil)
}

func (rh *RoleHandler) AddUserRole(ctx *gin.Context) {
	userID := ctx.Param("user_id")

	var roles []string
	if !rh.canBindJSON(ctx, &roles) {
		return
	}

	err := rh.roleService.AddUserRole(userID, roles...)
	if rh.isErrorWithMessage(ctx, err, 500, "Failed to add user roles") {
		return
	}

	response.Success(ctx, "User roles added successfully", nil)
}

func (rh *RoleHandler) GetUserRoles(ctx *gin.Context) {
	userID := ctx.Param("user_id")
	roles, err := rh.roleService.GetUserRoles(userID)
	if rh.isErrorWithMessage(ctx, err, 500, "Failed to get user roles") {
		return
	}
	response.Success(ctx, "Get user roles successfully", roles)
}


