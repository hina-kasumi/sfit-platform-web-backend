package handlers

import (
	"sfit-platform-web-backend/dtos"
	"sfit-platform-web-backend/services"
	"sfit-platform-web-backend/utils/response"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	*BaseHandler
	authSer    *services.AuthService
	jwtSer     *services.JwtService
	refreshSer *services.RefreshTokenService
}

func NewAuthHandler(baseHandler *BaseHandler, authSer *services.AuthService, jwtSer *services.JwtService, refreshSer *services.RefreshTokenService) *AuthHandler {
	return &AuthHandler{
		BaseHandler: baseHandler,
		authSer:     authSer,
		jwtSer:      jwtSer,
		refreshSer:  refreshSer,
	}
}

func (authHandler *AuthHandler) Register(ctx *gin.Context) {
	var userDto dtos.RegisterRequest
	if !authHandler.canBindJSON(ctx, &userDto) {
		return
	}

	accessToken, refreshToken, err := authHandler.authSer.Register(userDto.Username, userDto.Email, userDto.Password)
	if authHandler.isError(ctx, err) {
		return
	}

	authHandler.setRefreshTokenCookie(ctx, refreshToken)
	response.Success(ctx, "Register success", accessToken)
}

func (authHandler *AuthHandler) Login(ctx *gin.Context) {
	var userDto dtos.LoginRequest
	if !authHandler.canBindJSON(ctx, &userDto) {
		return
	}

	accessToken, refreshToken, err := authHandler.authSer.Login(userDto.Username, userDto.Email, userDto.Password)
	if authHandler.isError(ctx, err) {
		return
	}

	authHandler.setRefreshTokenCookie(ctx, refreshToken)
	response.Success(ctx, "Login success", accessToken)
}
func (authHandler *AuthHandler) Logout(ctx *gin.Context) {
	token := ctx.GetHeader("Authorization")
	if authHandler.isNilOrWhiteSpaceWithMessage(ctx, token, "Authorization token is required") {
		return
	}

	err := authHandler.authSer.Logout(token[7:]) // Remove "Bearer " prefix
	if authHandler.isError(ctx, err) {
		return
	}

	response.Success(ctx, "Logged out successfully", nil)
}
func (authHandler *AuthHandler) RefreshToken(ctx *gin.Context) {
	refreshToken, err := ctx.Cookie("refresh_token")

	if authHandler.isErrorWithMessage(ctx, err, 400, "Refresh token is required") {
		return
	}

	accessToken, newRefreshToken, err := authHandler.authSer.RefreshToken(refreshToken)
	if authHandler.isError(ctx, err) {
		return
	}

	authHandler.setRefreshTokenCookie(ctx, newRefreshToken)
	response.Success(ctx, "", accessToken)
}

func (authHandler *AuthHandler) setRefreshTokenCookie(ctx *gin.Context, refreshToken string) {
	ctx.SetCookie("refresh_token",
		refreshToken,
		int(authHandler.refreshSer.GetRefreshTokenExp()), // thời gian sống
		"/auth/refresh", // cookie sẽ được gửi ở refresh
		"",
		false,
		true, // http only
	)
}
