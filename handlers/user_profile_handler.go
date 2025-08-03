package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"sfit-platform-web-backend/dtos"
	"sfit-platform-web-backend/entities"
	"sfit-platform-web-backend/middlewares"
	"sfit-platform-web-backend/services"
	"sfit-platform-web-backend/utils/response"
	"time"
)

type UserProfileHandler struct {
	UserProfileService *services.UserProfileService
}

func NewUserProfileHandler(userProfileService *services.UserProfileService) *UserProfileHandler {
	return &UserProfileHandler{
		UserProfileService: userProfileService,
	}
}

func (profileHandler *UserProfileHandler) UpdateUserProfile(ctx *gin.Context) {
	userIDStr := middlewares.GetPrincipal(ctx)
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
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, 400, "Invalid request")
		return
	}

	socialLinkJSON, _ := json.Marshal(req.SocialLink)

	profile := entities.UserProfile{
		UserID:       userID,
		FullName:     req.FullName,
		ClassName:    req.ClassName,
		Khoa:         req.Khoa,
		Phone:        req.Phone,
		Introduction: req.Introduction,
		SocialLink:   string(socialLinkJSON),
		UpdatedAt:    time.Now(),
		Email:        req.Email,
	}

	createAt, updateAt, err := profileHandler.UserProfileService.UpdateUserProfile(&profile)
	if err != nil {
		response.Error(ctx, 500, "Failed to update user profile")
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
	if err != nil {
		response.Error(ctx, 400, "Invalid user ID")
		return
	}

	err = profileHandler.UserProfileService.DeleteUser(userID)
	if err != nil {
		response.Error(ctx, 500, "Failed to delete user")
		return
	}

	response.Success(ctx, "Delete user successfully", gin.H{
		"message": "User deleted",
	})
}

func (profileHandler *UserProfileHandler) GetUserProfile(ctx *gin.Context) {
	userIDSTr := ctx.Param("user_id")
	userID, err := uuid.Parse(userIDSTr)
	if err != nil {
		response.Error(ctx, 400, "Invalid user ID")
		return
	}

	profile, err := profileHandler.UserProfileService.GetUserProfile(userID)
	if err != nil {
		response.Error(ctx, 500, "User profile not found")
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
	if err != nil {
		response.Error(ctx, 400, "Invalid user ID")
		return
	}

	var req dtos.CreateUserProfileRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, 400, "Invalid request")
		return
	}

	socialLinkJSON, _ := json.Marshal(req.SocialLink)
	profile := entities.UserProfile{
		UserID:       userID,
		FullName:     req.FullName,
		ClassName:    req.ClassName,
		Khoa:         req.Khoa,
		Phone:        req.Phone,
		Introduction: req.Introduction,
		SocialLink:   string(socialLinkJSON),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	err = profileHandler.UserProfileService.CreateUserProfile(&profile)
	if err != nil {
		response.Error(ctx, 500, "Failed to create user profile")
		return
	}

	response.Success(ctx, "Create user profile successfully", gin.H{
		"user_id":  profile.UserID,
		"createAt": profile.CreatedAt,
	})
}
