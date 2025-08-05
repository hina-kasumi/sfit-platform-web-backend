package routes

import (
	"github.com/gin-gonic/gin"
	"sfit-platform-web-backend/handlers"
)

type TagRoute struct {
	tagHandler *handlers.TagHandler
}

func NewTagRoute(tagHandler *handlers.TagHandler) *TagRoute {
	return &TagRoute{tagHandler: tagHandler}
}

func (r *TagRoute) RegisterRoutes(router *gin.Engine) {
	group := router.Group("/tags")
	{
		group.GET("", r.tagHandler.GetAllTags)
	}
}
