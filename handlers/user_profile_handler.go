package handlers

import (
	"sfit-platform-web-backend/dtos"
	"sfit-platform-web-backend/entities"
	"sfit-platform-web-backend/middlewares"
	"sfit-platform-web-backend/services"
	"sfit-platform-web-backend/utils/response"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/google/uuid"
)

type UserProfileHandler struct {
	*BaseHandler
	UserProfileService *services.UserProfileService
}

func NewUserProfileHandler(baseHandler *BaseHandler, userProfileService *services.UserProfileService) *UserProfileHandler {
	return &UserProfileHandler{
		BaseHandler:        baseHandler,
		UserProfileService: userProfileService,
	}
}

func (profileHandler *UserProfileHandler) UpdateUserProfile(ctx *gin.Context) {
	userIDStr := ctx.Param("user_id")
	if userIDStr == "" {
		userIDStr = middlewares.GetPrincipal(ctx)
	}

	if userIDStr == "" {
		response.Error(ctx, 401, "unauthorized")
		return
	}
	userID, err := uuid.Parse(userIDStr) //chuyển thành uuid
	if err != nil {
		response.Error(ctx, 400, "Invalid user ID")
		return
	}

	var req dtos.UpdateUserProfileRequest
	if !profileHandler.canBindJSON(ctx, &req) {
		return
	}

	socialLinkJSON, _ := json.Marshal(req.SocialLink)

	profile := entities.UserProfile{
		Avatar:       req.Avatar,
		UserID:       userID,
		FullName:     req.FullName,
		ClassName:    req.ClassName,
		Khoa:         req.Khoa,
		Phone:        req.Phone,
		CoverImage:   req.CoverImage,
		Introduction: req.Introduction,
		SocialLink:   string(socialLinkJSON),
		Location:     req.Location,
		MSV:          req.Msv,
		UpdatedAt:    time.Now(),
		Email:        req.Email,
	}

	createAt, updateAt, err := profileHandler.UserProfileService.UpdateUserProfile(&profile)
	if profileHandler.isErrorWithMessage(ctx, err, 500, "Failed to update user profile") {
		return
	}

	response.Success(ctx, "Update user profile successfully", gin.H{
		"createAt": createAt,
		"updateAt": updateAt,
	})

}

func (profileHandler *UserProfileHandler) DeleteUser(ctx *gin.Context) {
	userIDSTr := ctx.Param("user_id")
	userID, err := uuid.Parse(userIDSTr)
	if profileHandler.isErrorWithMessage(ctx, err, 400, "Invalid user ID") {
		return
	}

	err = profileHandler.UserProfileService.DeleteUser(userID)
	if profileHandler.isErrorWithMessage(ctx, err, 500, "Failed to delete user") {
		return
	}

	response.Success(ctx, "Delete user successfully", gin.H{
		"message": "User deleted",
	})
}

func (profileHandler *UserProfileHandler) GetUserProfile(ctx *gin.Context) {
	userIDSTr := ctx.Param("user_id")
	userID, err := uuid.Parse(userIDSTr)
	if profileHandler.isErrorWithMessage(ctx, err, 400, "Invalid user ID") {
		return
	}

	profile, err := profileHandler.UserProfileService.GetUserProfile(userID)
	if profileHandler.isErrorWithMessage(ctx, err, 500, "Failed to get user profile") {
		return
	}

	response.Success(ctx, "Get user profile successfully", profile)
}

func (profileHandler *UserProfileHandler) CreateUserProfile(ctx *gin.Context) {
	userIDSTr := middlewares.GetPrincipal(ctx)
	if userIDSTr == "" {
		response.Error(ctx, 401, "unauthorized")
		return
	}

	userID, err := uuid.Parse(userIDSTr)
	if profileHandler.isErrorWithMessage(ctx, err, 400, "Invalid user ID") {
		return
	}

	var req dtos.CreateUserProfileRequest
	if !profileHandler.canBindJSON(ctx, &req) {
		return
	}

	socialLinkJSON, _ := json.Marshal(req.SocialLink)
	profile := entities.UserProfile{
		Avatar:       req.Avatar,
		CoverImage:   req.CoverImage,
		UserID:       userID,
		FullName:     req.FullName,
		ClassName:    req.ClassName,
		Khoa:         req.Khoa,
		Phone:        req.Phone,
		Introduction: req.Introduction,
		SocialLink:   string(socialLinkJSON),
		Location:     req.Location,
		MSV:          req.Msv,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	err = profileHandler.UserProfileService.CreateUserProfile(&profile)
	if profileHandler.isErrorWithMessage(ctx, err, 500, "Failed to create user profile") {
		return
	}

	response.Success(ctx, "Create user profile successfully", gin.H{
		"user_id":  profile.UserID,
		"createAt": profile.CreatedAt,
	})
}
