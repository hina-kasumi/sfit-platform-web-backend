package routes

import (
	"sfit-platform-web-backend/handlers"
	"sfit-platform-web-backend/entities"
	"sfit-platform-web-backend/middlewares"

	"github.com/gin-gonic/gin"
)

type CourseRoute struct {
    handler *handlers.CourseHandler
}

func NewCourseRoute(handler *handlers.CourseHandler) *CourseRoute {
    return &CourseRoute{handler: handler}
}

func (r *CourseRoute) RegisterRoutes(router *gin.Engine) {

    publicCourse := router.Group("course")
    publicCourse.Use(middlewares.EnforceAuthenticatedMiddleware())
    {
        publicCourse.GET("", r.handler.GetListCourse)
        publicCourse.GET("/:course_id", r.handler.GetCourseDetailByID)
        publicCourse.POST("/:course_id/favourite-course", r.handler.MarkCourseAsFavourite)
        publicCourse.DELETE("/:course_id/favourite-course", r.handler.UnmarkCourseAsFavourite)
    }

    protectedCourse := router.Group("/course")
    protectedCourse.Use(middlewares.EnforceAuthenticatedMiddleware())
    protectedCourse.Use(middlewares.RequireRoles(
        string(entities.RoleEnumAdmin),
        string(entities.RoleEnumHead),
        string(entities.RoleEnumVice),
    ))
    {
        protectedCourse.POST("", r.handler.CreateCourse)
        protectedCourse.PUT("", r.handler.UpdateCourse)
        protectedCourse.POST("/module", r.handler.AddModuleToCourse)
    }
}