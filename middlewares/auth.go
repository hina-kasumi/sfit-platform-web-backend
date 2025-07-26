package middlewares

import (
	"sfit-platform-web-backend/dtos"
	jwtservice "sfit-platform-web-backend/services/jwt_service"

	"github.com/gin-gonic/gin"
)

func UserLoaderMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		bearerToken := c.GetHeader("Authorization")
		if bearerToken == "" {
			c.Next()
			return
		}
		tokenPart := bearerToken[len("Bearer "):]
		err := jwtservice.ValidateToken(tokenPart)
		if err != nil {
			c.JSON(401, dtos.ErrorResponse{
				Code:    401,
				Message: "Invalid token",
			})
			c.Abort()
			return
		}

		sub, err := jwtservice.ParseToken(tokenPart)
		if err != nil {
			c.JSON(401, dtos.ErrorResponse{
				Code:    401,
				Message: "Invalid token"})
			c.Abort()
			return
		}
		c.Set("subject", sub)
	}
}

func EnforceAuthenticatedMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		sub, exists := c.Get("subject")

		if !exists || sub == "" {
			c.JSON(401, dtos.ErrorResponse{
				Code:    401,
				Message: "Unauthorized",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
