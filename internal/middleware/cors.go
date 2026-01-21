package middlewares

import (
	"net/http"
	"sfit-platform-web-backend/internal/config"

	"github.com/gin-gonic/gin"
)

func Cors(cfg *config.CorsConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		for _, origin := range cfg.AllowedOrigins {
			c.Writer.Header().Add("Access-Control-Allow-Origin", origin)
		}
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization, Accept")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
