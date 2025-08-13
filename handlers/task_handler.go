package handlers

import (
	"sfit-platform-web-backend/dtos"
	"sfit-platform-web-backend/services"
	"sfit-platform-web-backend/utils/response"

	"github.com/gin-gonic/gin"
)

type TaskHandler struct {
	*BaseHandler
	taskService *services.TaskService
}

func NewTaskHandler(base *BaseHandler, task *services.TaskService) *TaskHandler {
	return &TaskHandler{
		BaseHandler: base,
		taskService: task,
	}
}

// GET /tasks?page=...&pageSize=...
func (h *TaskHandler) GetListTasks(ctx *gin.Context) {
	pageNum, pageSizeNum, valid := parsePagination(ctx)
	if !valid {
		return
	}

	tasks, pagination, err := h.taskService.GetListTasks(pageNum, pageSizeNum)
	if h.isError(ctx, err) {
		return
	}

	responseData := dtos.TaskListResponse{
		Tasks:      tasks,
		Pagination: dtos.TaskPaginationResponse{
			CurrentPage: pagination.CurrentPage,
			TotalPages:  pagination.TotalPages,
			TotalTasks:  pagination.TotalTasks,
		},
	}
	response.Success(ctx, "Get list tasks successfully", responseData)
}