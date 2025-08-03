package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"sfit-platform-web-backend/dtos"
	"sfit-platform-web-backend/entities"
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
	sub, exists := ctx.Get("subject") //lấy subject từ context
	if !exists {
		response.Error(ctx, 401, "Unauthorized")
		return
	}

	claims, ok := sub.(jwt.MapClaims)
	if !ok {
		response.Error(ctx, 400, "Invalid subject format")
		return
	}

	subStr, ok := claims["sub"].(string) //lấy giá trị chuỗi từ sub
	if !ok {
		response.Error(ctx, 400, "Invalid subject format")
		return
	}

	userID, err := uuid.Parse(subStr) //chuyển thành uuid
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
	}

	createAt, updateAt, err := profileHandler.UserProfileService.UpdateUserProfile(&profile)
	if err != nil {
		response.Error(ctx, 500, "Failed to update user profile")
		return
	}

	response.Success(ctx, gin.H{
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

	response.Success(ctx, gin.H{
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

	response.Success(ctx, profile)
}
