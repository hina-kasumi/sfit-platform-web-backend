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
	if middlewares.HasRole(ctx, string(entities.RoleEnumMember)) {

	}
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
