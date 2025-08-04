package handlers

import (
	"github.com/gin-gonic/gin"
	"sfit-platform-web-backend/services"
	"sfit-platform-web-backend/utils/response"
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
