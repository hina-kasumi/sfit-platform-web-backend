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

	response.GetListResp(ctx, "Tasks retrieved successfully", query.Page, query.PageSize, totalCount, tasks)
}

func (th *TaskHandler) ListTasksByEventID(ctx *gin.Context) {
	eventID := ctx.Param("event_id")
	if th.isNilOrWhiteSpaceWithMessage(ctx, eventID, "Event ID is required") {
		return
	}

	var query *dtos.PageListQuery
	if !th.canBindQuery(ctx, &query) {
		return
	}

	tasks, totalCount, err := th.taskService.GetTasks(query.Page, query.PageSize, "", eventID)
	if th.isError(ctx, err) {
		return
	}

	response.GetListResp(ctx, "Tasks retrieved successfully", query.Page, query.PageSize, totalCount, tasks)
}

func (th *TaskHandler) ListTasksByUserID(ctx *gin.Context) {
	userID := ctx.Param("user_id")

	var query dtos.PageListQuery
	if !th.canBindQuery(ctx, &query) {
		return
	}

	tasks, totalCount, err := th.taskService.ListTasksByUserID(userID, query.Page, query.PageSize)
	if th.isError(ctx, err) {
		return
	}

	response.GetListResp(ctx, "Tasks retrieved successfully", query.Page, query.PageSize, totalCount, tasks)
}

func (th *TaskHandler) AddUserTask(ctx *gin.Context) {
	userID := ctx.Param("user_id")

	if th.isNilOrWhiteSpaceWithMessage(ctx, userID, "User ID is required") {
		return
	}

	var req dtos.AddUserTaskReq
	if !th.canBindJSON(ctx, &req) {
		return
	}

	_, err := th.taskService.AddUserTask(userID, req.TaskID)
	if th.isError(ctx, err) {
		return
	}

	response.Success(ctx, "Task created successfully", nil)
}

func (th *TaskHandler) DeleteUserTask(ctx *gin.Context) {
	userID := ctx.Param("user_id")
	taskID := ctx.Param("task_id")

	if th.isNilOrWhiteSpaceWithMessage(ctx, userID, "User ID is required") ||
		th.isNilOrWhiteSpaceWithMessage(ctx, taskID, "Task ID is required") {
		return
	}

	err := th.taskService.DeleteUserTask(userID, taskID)
	if th.isError(ctx, err) {
		return
	}

	response.Success(ctx, "Task deleted successfully", nil)
}
