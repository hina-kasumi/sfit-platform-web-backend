package controllers

import (
	"sfit-platform-web-backend/dtos"
	"sfit-platform-web-backend/services"
	"sfit-platform-web-backend/utits/response"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authSer    *services.AuthService
	jwtSer     *services.JwtService
	refreshSer *services.RefreshTokenService
}

func NewAuthController(authSer *services.AuthService, jwtSer *services.JwtService, refreshSer *services.RefreshTokenService) *AuthController {
	return &AuthController{
		authSer:    authSer,
		jwtSer:     jwtSer,
		refreshSer: refreshSer,
	}
}

func (authController *AuthController) RegisterRoutes(router *gin.Engine) {
	task := router.Group("/auth")
	task.POST("/register", authController.Register)
	task.POST("/login", authController.Login)
	task.POST("/logout", authController.Logout)
	task.POST("/refresh", authController.RefreshToken)
}

func (authController *AuthController) Register(ctx *gin.Context) {
	var userDto dtos.RegisterRequest
	if err := ctx.ShouldBindJSON(&userDto); err != nil {
		response.Error(ctx, 400, "Invalid input")
		return
	}

	accessToken, refreshToken, err := authController.authSer.Register(userDto.Username, userDto.Email, userDto.Password)
	if err != nil {
		response.Error(ctx, 500, err.Error())
		return
	}

	authController.setRefreshTokenCookie(ctx, refreshToken)
	response.Success(ctx, accessToken)
}

func (authController *AuthController) Login(ctx *gin.Context) {
	var userDto dtos.LoginRequest
	if err := ctx.ShouldBindJSON(&userDto); err != nil {
		response.Error(ctx, 400, "Invalid input")
		return
	}

	accessToken, refreshToken, err := authController.authSer.Login(userDto.Username, userDto.Email, userDto.Password)
	if err != nil {
		response.Error(ctx, 500, err.Error())
		return
	}

	authController.setRefreshTokenCookie(ctx, refreshToken)
	response.Success(ctx, accessToken)
}
func (authController *AuthController) Logout(ctx *gin.Context) {
	token := ctx.GetHeader("Authorization")
	if token == "" {
		response.Error(ctx, 400, "Authorization token is required")
		return
	}

	err := authController.authSer.Logout(token[7:]) // Remove "Bearer " prefix
	if err != nil {
		response.Error(ctx, 500, err.Error())
		return
	}

	response.Success(ctx, "Logged out successfully")
}
func (authController *AuthController) RefreshToken(ctx *gin.Context) {
	refreshToken, err := ctx.Cookie("refresh_token")
	if err != nil {
		response.Error(ctx, 400, "Refresh token is required")
		return
	}

	accessToken, newRefreshToken, err := authController.authSer.RefreshToken(refreshToken)
	if err != nil {
		response.Error(ctx, 500, err.Error())
		return
	}

	authController.setRefreshTokenCookie(ctx, newRefreshToken)
	response.Success(ctx, accessToken)
}

func (authController *AuthController) setRefreshTokenCookie(ctx *gin.Context, refreshToken string) {
	ctx.SetCookie("refresh_token",
		refreshToken,
		int(authController.refreshSer.GetRefreshTokenExp()), // thời gian sống
		"/auth/refresh", // cookie sẽ được gửi ở refresh
		"",
		false,
		true, // http only
	)
}
