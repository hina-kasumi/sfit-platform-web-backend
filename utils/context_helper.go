package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Lấy userID từ context
func GetUserIDFromContext(ctx *gin.Context) uuid.UUID {
	userIDInterface, exists := ctx.Get("user_id")

	if !exists {
		return uuid.Nil
	}

	switch v := userIDInterface.(type) {
	case string:
		if parsedID, err := uuid.Parse(v); err == nil {
			return parsedID
		}
	case uuid.UUID:
		return v
	}

	return uuid.Nil
}
