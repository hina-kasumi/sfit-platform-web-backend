package response

import (
	"net/http"
	"sfit-platform-web-backend/dtos"

	"github.com/gin-gonic/gin"
)

func Success(c *gin.Context, message string, data any) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": message,
		"data":    data,
	})
}

// Lỗi (ví dụ 400 Bad Request)
func Error(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{
		"status":  "error",
		"message": message,
	})
}

func GetListResp(c *gin.Context, message string, page, pageSize int, totalCount int64, data any) {
	Success(c, message, dtos.PageListResp{
		TotalCount: totalCount,
		Page:       page,
		PageSize:   pageSize,
		Items:      data,
	})
}
