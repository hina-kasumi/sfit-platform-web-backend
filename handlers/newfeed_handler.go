package handlers

import (
	"sfit-platform-web-backend/middlewares"
	"sfit-platform-web-backend/services"
	"sfit-platform-web-backend/utils/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type NewFeedHandler struct {
	*BaseHandler
	newfeedService *services.NewFeedService
}

func NewNewFeedHandler(newfeedService *services.NewFeedService) *NewFeedHandler {
	return &NewFeedHandler{
		newfeedService: newfeedService,
	}
}

func (h *NewFeedHandler) GetNewFeed(ctx *gin.Context) {
	userID := middlewares.GetPrincipal(ctx)
	userUUID, err := uuid.Parse(userID)
	if h.isError(ctx, err) {
		return
	}
	newfeed, err := h.newfeedService.GetNewFeed(userUUID)
	if h.isError(ctx, err) {
		return
	}

	response.Success(ctx, "Get newfeed success", newfeed)
}
