package routes

import (
	"sfit-platform-web-backend/entities"
	"sfit-platform-web-backend/handlers"
	"sfit-platform-web-backend/middlewares"

	"github.com/gin-gonic/gin"
)

type LessonRoute struct {
	lessonHandler *handlers.LessonHandler
}

func NewLessonRoute(lessonHandler *handlers.LessonHandler) *LessonRoute {
	return &LessonRoute{
		lessonHandler: lessonHandler,
	}
}

func (lessonRou *LessonRoute) RegisterRoutes(router *gin.Engine) {
	lessonGroup := router.Group("/modules/:module_id/lessons")
	lessonGroup.Use(middlewares.EnforceAuthenticatedMiddleware())
	lessonGroup.GET("/:lesson_id", lessonRou.lessonHandler.GetLessonByID)
	lessonGroup.Use(middlewares.RequireRoles(
		string(entities.RoleEnumAdmin),
		string(entities.RoleEnumHead),
		string(entities.RoleEnumVice),
	))
	lessonGroup.POST("/", lessonRou.lessonHandler.CreateLesson)
	lessonGroup.PUT("/:lesson_id", lessonRou.lessonHandler.UpdateLesson)
	lessonGroup.DELETE("/:lesson_id", lessonRou.lessonHandler.DeleteLesson)
}
