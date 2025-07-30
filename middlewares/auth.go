package middlewares

import (
	"sfit-platform-web-backend/services"
	"sfit-platform-web-backend/utils/response"

	"github.com/gin-gonic/gin"
)

func UserLoaderMiddleware(jwtSer *services.JwtService) gin.HandlerFunc {
	return func(c *gin.Context) {
		bearerToken := c.GetHeader("Authorization")
		if bearerToken == "" {
			c.Next()
			return
		}
		tokenPart := bearerToken[len("Bearer "):]
		sub, err := jwtSer.ParseToken(tokenPart)
		if err != nil {
			response.Error(c, 401, "Invalid token")
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
			response.Error(c, 401, "Unauthorized")
			c.Abort()
			return
		}
		c.Next()
	}
}
