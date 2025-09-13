package handlers

import (
	"sfit-platform-web-backend/dtos"
	"sfit-platform-web-backend/entities"
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

func (userHandler *UserHandler) UpdateUser(ctx *gin.Context) {
	var req dtos.UpdateUserDto
	if !userHandler.canBindJSON(ctx, &req) {
		return
	}

	userID := ctx.Param("user_id")
	if userID == "" {
		response.Error(ctx, 401, "Unauthorized")
		return
	}

	user, err := userHandler.userSer.GetUserByID(userID)
	if err != nil {
		response.Error(ctx, 404, "User not found")
		return
	}

	if req.Email != "" {
		user.Email = req.Email
	}
	if req.NewPassword != "" && user.IsValidPassword(req.OldPassword) == nil ||
		middlewares.HasRole(ctx, string(entities.RoleEnumAdmin)) {
		user.SetPassword(req.NewPassword)
	}

	_, err = userHandler.userSer.UpdateUser(user)
	if userHandler.isErrorWithMessage(ctx, err, 500, "Change password error") {
		return
	}

	response.Success(ctx, "Change password successfully", nil)
}

func (userHandler *UserHandler) GetUserList(ctx *gin.Context) {
	var query dtos.PageListQuery
	if !userHandler.canBindQuery(ctx, &query) {
		return
	}

	users, page, pageSize, total, err := userHandler.userSer.GetUserList(query.Page, query.PageSize)
	if userHandler.isErrorWithMessage(ctx, err, 500, "Failed to get user list") {
		return
	}

	res := dtos.PageListResp{
		Items:      users,
		Page:       page,
		PageSize:   pageSize,
		TotalCount: total,
	}

	response.Success(ctx, "Get user list successfully", res)
}
