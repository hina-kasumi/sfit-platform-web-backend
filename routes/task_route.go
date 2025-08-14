package routes

import (
	"sfit-platform-web-backend/entities"
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
	taskGroup.GET("/:task_id", tr.handler.GetTaskDetail)
	taskGroup.Use(middlewares.RequireRoles(
		string(entities.RoleEnumAdmin),
		string(entities.RoleEnumHead),
		string(entities.RoleEnumVice),
	))
	taskGroup.POST("", tr.handler.CreateTask)
	taskGroup.PUT("/:task_id", tr.handler.UpdateTask)
	taskGroup.DELETE("/:task_id", tr.handler.DeleteTask)

	userTasksGroup := router.Group("/users/:user_id/tasks")
	userTasksGroup.Use(middlewares.EnforceAuthenticatedMiddleware())
	userTasksGroup.GET("", tr.handler.ListTasksByUserID)
	userTasksGroup.PATCH("/:task_id", tr.handler.UpdateTaskUserStatus)
	userTasksGroup.Use(middlewares.RequireRoles(
		string(entities.RoleEnumAdmin),
		string(entities.RoleEnumHead),
		string(entities.RoleEnumVice),
	))
	userTasksGroup.POST("", tr.handler.AddUserTask)
	userTasksGroup.DELETE("/:task_id", tr.handler.DeleteUserTask)

	// Lấy danh sách task theo event_id
	router.GET("events/:event_id/tasks", tr.handler.ListTasksByEventID)
}
