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

    publicCourse := router.Group("")
    publicCourse.Use(middlewares.EnforceAuthenticatedMiddleware())
    {
        publicCourse.GET("/course", r.handler.GetListCourse)
        publicCourse.GET("/course/:course_id", r.handler.GetCourseDetailByID)
        publicCourse.POST("course/:course_id/favourite", r.handler.MarkCourseAsFavourite)
        publicCourse.DELETE("course/:course_id/favourite", r.handler.UnmarkCourseAsFavourite)
        publicCourse.GET("users/:user_id/courses/:course_id/progress", r.handler.GetUserProgressInCourse)
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