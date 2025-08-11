package middlewares

import (
	"sfit-platform-web-backend/services"
	"sfit-platform-web-backend/utils/response"
	"slices"

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

/* RequireRoles checks if the user has one of the required roles
* Have to be used after UserLoaderMiddleware
* input: roles ...string (recommend use entities.RoleEnum)
 */
func RequireRoles(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRoles := GetRoles(c)
		if len(userRoles) == 0 {
			response.Error(c, 403, "Forbidden")
			c.Abort()
			return
		}

		for _, role := range roles {
			if slices.Contains(userRoles, role) {
				c.Next()
				return
			}
		}
		response.Error(c, 403, "Forbidden")
		c.Abort()
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

func HasRole(c *gin.Context, role ...string) bool {
	userRoles := GetRoles(c)
	for _, r := range role {
		if slices.Contains(userRoles, r) {
			return true
		}
	}
	return false
}
