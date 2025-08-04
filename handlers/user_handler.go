package handlers

import (
	"sfit-platform-web-backend/dtos"
	"sfit-platform-web-backend/middlewares"
	"sfit-platform-web-backend/services"
	"sfit-platform-web-backend/utils/response"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	*BaseHandler
	userSer *services.UserService
}

func NewUserHandler(baseHandler *BaseHandler, userSer *services.UserService) *UserHandler {
	return &UserHandler{
		BaseHandler: baseHandler,
		userSer:     userSer,
	}
}

func (userHandler *UserHandler) ChangePassword(ctx *gin.Context) {
	var req dtos.ChangePasswordRequest
	if !userHandler.canBindJSON(ctx, &req) {
		return
	}

	userID := middlewares.GetPrincipal(ctx)
	if userID == "" {
		response.Error(ctx, 401, "Unauthorized")
		return
	}

	err := userHandler.userSer.ChangePassword(userID, req.OldPassword, req.NewPassword)
	if userHandler.isErrorWithMessage(ctx, err, 500, "Change password error") {
		return
	}

	response.Success(ctx, "Change password successfully", nil)
}

func (userHandler *UserHandler) GetUserList(ctx *gin.Context) {
	var query dtos.UserListQuery
	if !userHandler.canBindQuery(ctx, &query) {
		return
	}

	users, page, pageSize, total, err := userHandler.userSer.GetUserList(query.Page, query.PageSize)
	if userHandler.isErrorWithMessage(ctx, err, 500, "Failed to get user list") {
		return
	}

	res := dtos.UserListResponse{
		Users:    users,
		Page:     page,
		PageSize: pageSize,
		Total:    total,
	}

	response.Success(ctx, "Get user list successfully", res)
}
