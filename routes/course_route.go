package routes

import (
	"sfit-platform-web-backend/handlers"
	"sfit-platform-web-backend/middlewares"

	"github.com/gin-gonic/gin"
)

type CourseRoutes struct {
	courseHandler *handlers.CourseHandler
}

func NewCourseRoute(courseHandler *handlers.CourseHandler) *CourseRoutes {
	return &CourseRoutes{
		courseHandler: courseHandler,
	}
}

func (courseRou *CourseRoutes) RegisterRoutes(router *gin.Engine) {
	courseHandler := courseRou.courseHandler

	course := router.Group("/course")

	// Public routes (không cần authentication)
	course.POST("/create", courseHandler.CreateCourse)//ok
	course.GET("/registered/:user_id", courseHandler.GetRegisteredCourses)//ok
	course.GET("/lessons/:course_id", courseHandler.GetCourseLessons)//ok
	course.GET("/list-registered", courseHandler.GetRegisteredUsers)//
	course.DELETE("/:course_id", courseHandler.DeleteCourse)

	// Protected routes (cần authentication)
	protected := course.Group("")
	protected.Use(middlewares.EnforceAuthenticatedMiddleware())
	protected.POST("/rate", courseHandler.RateCourse)//chưa check xem có tồn tại khóa học hay không
	protected.POST("/register", courseHandler.RegisterUserToCourse)
}
