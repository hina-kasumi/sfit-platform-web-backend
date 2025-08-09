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
		claims, err := jwtSer.ParseToken(tokenPart)
		if err != nil {
			response.Error(c, 401, "Invalid token")
			c.Abort()
			return
		}
		sub, err := claims.GetSubject()
		if err != nil {
			response.Error(c, 401, "Invalid token")
			c.Abort()
			return
		}
		roles, ok := claims["roles"]
		if !ok {
			response.Error(c, 401, "Invalid token")
			c.Abort()
			return
		}
		c.Set("subject", sub)
		c.Set("roles", roles)
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

func GetPrincipal(c *gin.Context) string {
	sub, exists := c.Get("subject")
	if !exists || sub == "" {
		return ""
	}
	return sub.(string)
}

func GetRoles(c *gin.Context) []string {
	raw, exists := c.Get("roles")
	if !exists {
		return nil
	}
	if raw == nil {
		return nil
	}
	inter := raw.([]any)
	var roles []string
	for _, item := range inter {
		if str, ok := item.(string); ok {
			roles = append(roles, str)
		}
	}
	return roles
}
