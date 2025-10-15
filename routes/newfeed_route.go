package routes

import (
	"sfit-platform-web-backend/handlers"

	"github.com/gin-gonic/gin"
)

type NewFeedRoute struct {
	newFeedHandler *handlers.NewFeedHandler
}

func NewNewFeedRoute(newFeedHandler *handlers.NewFeedHandler) *NewFeedRoute {
	return &NewFeedRoute{
		newFeedHandler: newFeedHandler,
	}
}

func (r *NewFeedRoute) RegisterRoutes(router *gin.Engine) {
	newfeed := router.Group("/newfeed")
	{
		newfeed.GET("", r.newFeedHandler.GetNewFeed)
	}
}
