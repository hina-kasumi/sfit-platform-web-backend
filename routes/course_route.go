package routes

import (
	"sfit-platform-web-backend/handlers"
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
	router.POST("/course", r.handler.CreateCourse)
	router.GET("/course", r.handler.GetListCourse)
	router.GET("/course/:course_id", r.handler.GetCourseDetailByID)
	router.POST("/course/favourite", r.handler.MarkCourseAsFavourite)
	router.DELETE("/course/favourite", r.handler.UnmarkCourseAsFavourite)
	router.PUT("/course", r.handler.UpdateCourse)

	protected := router.Group("")
    protected.Use(middlewares.EnforceAuthenticatedMiddleware())
    protected.GET("/users/:user_id/courses", r.handler.GetRegisteredCourses)
    protected.GET("/course/:course_id/lessons", r.handler.GetCourseLessons)
	protected.POST("/users/:user_id/rate/courses/:course_id", r.handler.RateCourse)
	protected.DELETE("/course/:course_id", middlewares.RequireRoles("ADMIN","HEAD"), r.handler.DeleteCourse)
	protected.POST("/users/:user_id/courses", r.handler.RegisterUserToCourse)
	protected.GET("/course/:course_id/users", middlewares.RequireRoles("ADMIN","HEAD"), r.handler.GetRegisteredUsers)
}
