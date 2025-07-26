package authcontroller

import (
	"sfit-platform-web-backend/dtos"
	authservice "sfit-platform-web-backend/services/auth_service"
	jwtservice "sfit-platform-web-backend/services/jwt_service"

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

	setRefreshTokenCookie(ctx, refreshToken)
	ctx.JSON(200, gin.H{
		"access_token": accessToken,
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

	setRefreshTokenCookie(ctx, refreshToken)
	ctx.JSON(200, gin.H{
		"access_token": accessToken,
	})
}
func Logout(ctx *gin.Context) {
	token := ctx.GetHeader("Authorization")
	if token == "" {
		ctx.JSON(400, dtos.NewErrorResponse(400, "Authorization token is required"))
		return
	}

	err := authservice.Logout(token[7:]) // Remove "Bearer " prefix
	if err != nil {
		ctx.JSON(500, dtos.NewErrorResponse(500, err.Error()))
		return
	}

	ctx.JSON(200, gin.H{"message": "Logged out successfully"})
}
func RefreshToken(ctx *gin.Context) {
	refreshToken, err := ctx.Cookie("refresh_token")
	if err != nil {
		ctx.JSON(400, dtos.NewErrorResponse(400, "Refresh token is required"))
		return
	}

	accessToken, newRefreshToken, err := authservice.RefreshToken(refreshToken)
	if err != nil {
		ctx.JSON(500, dtos.NewErrorResponse(500, err.Error()))
		return
	}

	setRefreshTokenCookie(ctx, newRefreshToken)
	ctx.JSON(200, gin.H{
		"access_token": accessToken,
	})
}

func setRefreshTokenCookie(ctx *gin.Context, refreshToken string) {
	ctx.SetCookie("refresh_token",
		refreshToken,
		int(jwtservice.GetRefreshTokenExp()), // thời gian sống
		"/auth/refresh",                      // cookie sẽ được gửi ở refresh
		"",
		false,
		true, // http only
	)
}
