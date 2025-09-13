package handlers

import (
	"sfit-platform-web-backend/dtos"
	"sfit-platform-web-backend/entities"
	"sfit-platform-web-backend/middlewares"
	"sfit-platform-web-backend/services"
	"sfit-platform-web-backend/utils/response"

	"github.com/gin-gonic/gin"
)

type LessonHandler struct {
	*BaseHandler
	lessonService *services.LessonService
}

func NewLessonHandler(baseHandler *BaseHandler, lessonService *services.LessonService) *LessonHandler {
	return &LessonHandler{
		BaseHandler:   baseHandler,
		lessonService: lessonService,
	}
}

func (h *LessonHandler) GetLessonByID(ctx *gin.Context) {
	lessonID := ctx.Param("lesson_id")
	lesson, err := h.lessonService.GetLessonByID(lessonID)
	if h.isError(ctx, err) {
		return
	}
	response.Success(ctx, "Lesson retrieved successfully", lesson)
}

func (h *LessonHandler) CreateLesson(ctx *gin.Context) {
	moduleID := ctx.Param("module_id")
	if h.isNilOrWhiteSpaceWithMessage(ctx, moduleID, "Module ID is required") {
		return
	}

	createLessonRequest := dtos.LessonRequest{}
	if !h.canBindJSON(ctx, &createLessonRequest) {
		return
	}

	lesson, err := h.lessonService.CreateNewLesson(moduleID, createLessonRequest)
	if h.isError(ctx, err) {
		return
	}

	response.Success(ctx, "Lesson created successfully", lesson)
}

func (h *LessonHandler) UpdateLesson(ctx *gin.Context) {
	moduleID := ctx.Param("module_id")
	lessonID := ctx.Param("lesson_id")
	if h.isNilOrWhiteSpaceWithMessage(ctx, lessonID, "Lesson ID is required") {
		return
	}

	var req dtos.LessonRequest
	if !h.canBindJSON(ctx, &req) {
		return
	}

	if err := h.lessonService.UpdateLesson(moduleID, lessonID, req); err != nil {
		h.isError(ctx, err)
		return
	}

	response.Success(ctx, "Lesson updated successfully", nil)
}

func (h *LessonHandler) DeleteLesson(ctx *gin.Context) {
	lessonID := ctx.Param("lesson_id")
	if err := h.lessonService.DeleteLessonByID(lessonID); err != nil {
		h.isError(ctx, err)
		return
	}
	response.Success(ctx, "Lesson deleted successfully", nil)
}

func (h *LessonHandler) UpdateStatusLessonAttendance(ctx *gin.Context) {
	userID := ctx.Param("user_id")
	lessonID := ctx.Param("lesson_id")
	var req dtos.UpdateStatusLessonAttendanceReq
	if !h.canBindJSON(ctx, &req) {
		return
	}

	currentUserID := middlewares.GetPrincipal(ctx)
	lesson, err := h.lessonService.GetLessonByID(lessonID)
	if h.isError(ctx, err) {
		return
	}
	if !middlewares.HasRole(ctx,
		string(entities.RoleEnumAdmin),
		string(entities.RoleEnumHead),
		string(entities.RoleEnumVice),
		string(entities.RoleEnumTeacher),
	) {
		if lesson.Type == entities.OfflineLesson {
			response.Error(ctx, 403, "Only admins or heads can update offline lesson attendance")
		} else if lesson.Type == entities.QuizLesson && req.Answer == nil {
			response.Error(ctx, 403, "you can't update quiz lesson attendance with out an answer")
		}
		if userID != currentUserID {
			response.Error(ctx, 403, "you can't update another user's lesson attendance")
		}
	}
	err = h.lessonService.UpdateStatusLessonAttendance(
		userID,
		lesson,
		req.Status,
		req.DeviceID,
		req.Answer,
		currentUserID,
		req.Duration,
	)
	if h.isError(ctx, err) {
		return
	}
	response.Success(ctx, "Lesson attendance status updated successfully", nil)
}

func (h *LessonHandler) GetUsersByLessonID(ctx *gin.Context) {
	lessonID := ctx.Param("lesson_id")
	query := dtos.GetUserAttendanceLessonReq{}
	if !h.canBindQuery(ctx, &query) {
		return
	}
	users, total, err := h.lessonService.GetUsersByLessonID(lessonID, query)
	if h.isError(ctx, err) {
		return
	}
	response.GetListResp(ctx, "Users retrieved successfully", query.Page, query.PageSize, total, users)
}
