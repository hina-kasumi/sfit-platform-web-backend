package routes

import (
	"sfit-platform-web-backend/handlers"
	"sfit-platform-web-backend/middlewares"

	"github.com/gin-gonic/gin"
)

type TaskRouter struct {
	handler *handlers.TaskHandler
}

func NewTaskRouter(handler *handlers.TaskHandler) *TaskRouter {
	return &TaskRouter{
		handler: handler,
	}
}

func (tr *TaskRouter) RegisterRoutes(router *gin.Engine) {
	taskGroup := router.Group("/tasks")
	taskGroup.Use(middlewares.EnforceAuthenticatedMiddleware())
	taskGroup.GET("", tr.handler.ListTasks)
	taskGroup.POST("", tr.handler.CreateTask)
	taskGroup.GET("/:task_id", tr.handler.GetTaskDetail)
	taskGroup.PUT("/:task_id", tr.handler.UpdateTask)
	taskGroup.DELETE("/:task_id", tr.handler.DeleteTask)
}
