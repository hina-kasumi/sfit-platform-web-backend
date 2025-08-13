package routes

import (
	"sfit-platform-web-backend/handlers"
	"sfit-platform-web-backend/middlewares"

	"github.com/gin-gonic/gin"
)

type TaskRoute struct {
	handler *handlers.TaskHandler
}

func NewTaskRoute(handler *handlers.TaskHandler) *TaskRoute {
	return &TaskRoute{handler: handler}
}

func (r *TaskRoute) RegisterRoutes(router *gin.Engine) {
	publicTask := router.Group("/tasks")
	publicTask.Use(middlewares.EnforceAuthenticatedMiddleware())
	{
		publicTask.GET("", r.handler.GetListTasks)
	}
}
