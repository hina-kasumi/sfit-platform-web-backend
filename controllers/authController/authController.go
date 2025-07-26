package authcontroller

import (
	"sfit-platform-web-backend/dtos"
	authservice "sfit-platform-web-backend/services/authService"

	"github.com/gin-gonic/gin"
)

func Register(ctx *gin.Context) {
	var userDto dtos.RegisterRequest
	if err := ctx.ShouldBindJSON(&userDto); err != nil {
		ctx.JSON(400, dtos.NewErrorResponse(400, "Invalid input"))
		return
	}

	accessToken, refreshToken, err := authservice.Register(userDto.Username, userDto.Email, userDto.Password)
	if err != nil {
		ctx.JSON(500, dtos.NewErrorResponse(500, err.Error()))
		return
	}

	ctx.JSON(200, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}
func Login(ctx *gin.Context) {
	var userDto dtos.LoginRequest
	if err := ctx.ShouldBindJSON(&userDto); err != nil {
		ctx.JSON(400, dtos.NewErrorResponse(400, "Invalid input"))
		return
	}

	accessToken, refreshToken, err := authservice.Login(userDto.Username, userDto.Email, userDto.Password)
	if err != nil {
		ctx.JSON(500, dtos.NewErrorResponse(500, err.Error()))
		return
	}

	ctx.JSON(200, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}
func Logout(ctx *gin.Context) {
	token := ctx.GetHeader("Authorization")
	if token == "" {
		ctx.JSON(400, dtos.NewErrorResponse(400, "Authorization token is required"))
		return
	}

	err := authservice.Logout(token)
	if err != nil {
		ctx.JSON(500, dtos.NewErrorResponse(500, err.Error()))
		return
	}

	ctx.JSON(200, gin.H{"message": "Logged out successfully"})
}
func RefreshToken(ctx *gin.Context) {
	// Logic to refresh user token
}
