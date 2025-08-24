package handlers

import (
	"fmt"
	"sfit-platform-web-backend/utils/response"
	"strings"

	"github.com/gin-gonic/gin"
)

type BaseHandler struct {
}

func NewBaseHandler() *BaseHandler {
	return &BaseHandler{}
}

func (bh *BaseHandler) isErrorWithMessage(ctx *gin.Context, err error, code int, message string) bool {
	if err != nil {
		response.Error(ctx, code, message)
		return true
	}
	return false
}

func (bh *BaseHandler) isError(ctx *gin.Context, err error) bool {
	if err != nil {
		return bh.isErrorWithMessage(ctx, err, 500, err.Error())
	}
	return false
}

func (bh *BaseHandler) canBindJSON(ctx *gin.Context, dto any) bool {
	if err := ctx.ShouldBindJSON(dto); err != nil {
		fmt.Printf("Bind JSON Error: %v\n", err)
		response.Error(ctx, 400, "Invalid input")
		return false
	}
	return true
}

func (bh *BaseHandler) canBindQuery(ctx *gin.Context, dto any) bool {
	if err := ctx.ShouldBindQuery(dto); err != nil {
		response.Error(ctx, 400, "Invalid input")
		return false
	}
	return true
}

func (bh *BaseHandler) isNilOrWhiteSpaceWithMessage(ctx *gin.Context, dto any, message string) bool {
	if dto == nil {
		response.Error(ctx, 400, message)
		return true
	}

	// Nếu là *string hoặc string
	switch v := dto.(type) {
	case string:
		if strings.TrimSpace(v) == "" {
			response.Error(ctx, 400, message)
			return true
		}
		return false
	case *string:
		if v == nil {
			response.Error(ctx, 400, message)
			return true
		}
		return false
	default:
		return false
	}
}
