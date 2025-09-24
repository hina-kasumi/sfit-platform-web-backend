package routes

import (
	"sfit-platform-web-backend/entities"
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

	publicCourse := router.Group("")
	publicCourse.Use(middlewares.EnforceAuthenticatedMiddleware())
	{
		publicCourse.GET("/courses", r.handler.GetListCourse)
		publicCourse.GET("/courses/:course_id", r.handler.GetCourseDetailByID)
		publicCourse.POST("courses/:course_id/favourite", r.handler.MarkCourseAsFavourite)
		publicCourse.DELETE("courses/:course_id/favourite", r.handler.UnmarkCourseAsFavourite)
		publicCourse.GET("users/:user_id/courses/:course_id/progress", r.handler.GetUserProgressInCourse)

		publicCourse.GET("/users/:user_id/courses", r.handler.GetRegisteredCourses)
		publicCourse.GET("/courses/:course_id/lessons", r.handler.GetCourseLessons)
		publicCourse.POST("/users/:user_id/rate/courses/:course_id", r.handler.RateCourse)
		publicCourse.POST("/users/courses", r.handler.RegisterUserToCourse)
	}

	protectedCourse := router.Group("/courses")
	protectedCourse.Use(middlewares.EnforceAuthenticatedMiddleware())
	protectedCourse.Use(middlewares.RequireRoles(
		string(entities.RoleEnumAdmin),
		string(entities.RoleEnumHead),
		string(entities.RoleEnumVice),
		string(entities.RoleEnumTeacher),
	))
	{
		protectedCourse.POST("", r.handler.CreateCourse)
		protectedCourse.PUT("/:course_id", r.handler.UpdateCourse)
		protectedCourse.POST("/:course_id/modules", r.handler.AddModuleToCourse)
		protectedCourse.DELETE("/modules/:module_id", r.handler.DeleteModuleInCourse)
		protectedCourse.DELETE("/:course_id", r.handler.DeleteCourse)
		protectedCourse.GET("/:course_id/users", r.handler.GetRegisteredUsers)
	}

}
