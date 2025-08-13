package handlers

import (
	"sfit-platform-web-backend/dtos"
	"sfit-platform-web-backend/middlewares"
	"sfit-platform-web-backend/services"
	"sfit-platform-web-backend/utils/response"

	"github.com/gin-gonic/gin"
)

type TaskHandler struct {
	*BaseHandler
	taskService *services.TaskService
}

func NewTaskHandler(baseHandler *BaseHandler, taskService *services.TaskService) *TaskHandler {
	return &TaskHandler{
		BaseHandler: baseHandler,
		taskService: taskService,
	}
}

func (th *TaskHandler) CreateTask(ctx *gin.Context) {
	var req dtos.CreateTaskReq
	if !th.canBindJSON(ctx, &req) {
		return
	}

	userID := middlewares.GetPrincipal(ctx)

	task, err := th.taskService.CreateTask(userID, req.Name, req.Description, req.EventID, req.StartDate, req.Deadline)
	if th.isError(ctx, err) {
		return
	}

	response.Success(ctx, "Task created successfully", task)
}

func (th *TaskHandler) GetTaskDetail(ctx *gin.Context) {
	taskID := ctx.Param("task_id")

	task, err := th.taskService.GetTaskByID(taskID)
	if th.isError(ctx, err) {
		return
	}

	response.Success(ctx, "success", task)
}

func (th *TaskHandler) UpdateTask(ctx *gin.Context) {
	var req dtos.UpdateTaskReq
	if !th.canBindJSON(ctx, &req) {
		return
	}

	taskID := ctx.Param("task_id")

	err := th.taskService.UpdateTask(taskID, req.Name, req.Description, req.PercentComplete, req.StartDate, req.Deadline)
	if th.isError(ctx, err) {
		return
	}

	response.Success(ctx, "Task updated successfully", nil)
}

func (th *TaskHandler) DeleteTask(ctx *gin.Context) {
	taskID := ctx.Param("task_id")

	err := th.taskService.DeleteTask(taskID)
	if th.isError(ctx, err) {
		return
	}

	response.Success(ctx, "Task deleted successfully", nil)
}

func (th *TaskHandler) ListTasks(ctx *gin.Context) {
	var query dtos.ListTaskQuery
	if !th.canBindQuery(ctx, &query) {
		return
	}

	tasks, totalCount, err := th.taskService.GetTasks(query.Page, query.PageSize, query.Name, query.EventID)
	if th.isError(ctx, err) {
		return
	}

	response.Success(ctx, "Tasks retrieved successfully", gin.H{
		"tasks":       tasks,
		"total_count": totalCount,
		"page":        query.Page,
		"page_size":   query.PageSize,
	})
}
