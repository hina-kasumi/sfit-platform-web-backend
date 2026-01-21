package handlers

import (
	"sfit-platform-web-backend/internal/services"
	"sfit-platform-web-backend/internal/utils/response"

	"github.com/gin-gonic/gin"
)

type TagHandler struct {
	*BaseHandler
	tagService *services.TagService
}

func NewTagHandler(base *BaseHandler, tagService *services.TagService) *TagHandler {
	return &TagHandler{
		BaseHandler: base,
		tagService:  tagService,
	}
}

func (h *TagHandler) GetAllTags(ctx *gin.Context) {
	tags, err := h.tagService.GetAll()
	if h.isError(ctx, err) {
		return
	}
	response.Success(ctx, "", tags)
}
