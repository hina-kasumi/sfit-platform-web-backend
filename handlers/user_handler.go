package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"sfit-platform-web-backend/dtos"
	"sfit-platform-web-backend/services"
	"sfit-platform-web-backend/utils/response"
)

type UserHandler struct {
	userSer *services.UserService
}

func NewUserHandler(userSer *services.UserService) *UserHandler {
	return &UserHandler{userSer: userSer}
}

func (userHandler *UserHandler) ChangePassword(ctx *gin.Context) {
	var req dtos.ChangePasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, 401, "Invalid request")
		return
	}

	subRaw, exists := ctx.Get("subject")
	if !exists {
		response.Error(ctx, 401, "Unauthorized - subject not found")
		return
	}

	claims, ok := subRaw.(jwt.MapClaims)
	if !ok {
		response.Error(ctx, 401, "Invalid subject type")
		return
	}

	userId, ok := claims["sub"].(string)
	if !ok || userId == "" {
		response.Error(ctx, 401, "Subject claim missing or not a string")
		return
	}

	err := userHandler.userSer.ChangePassword(userId, req.OldPassword, req.NewPassword)
	if err != nil {
		response.Error(ctx, 401, "Change password error")
		return
	}

	response.Success(ctx, "Change password success")
}
