package routes

import (
	"sfit-platform-web-backend/handlers"

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
	course.GET("/registered/:user_id", courseHandler.GetRegisteredCourses)
	course.GET("/lessons/:course_id", courseHandler.GetCourseLessons)
	course.POST("/rate", courseHandler.RateCourse)
	course.DELETE("/:course_id", courseHandler.DeleteCourse)
	course.GET("/list-registered", courseHandler.GetRegisteredUsers)
	course.POST("/register", courseHandler.RegisterUserToCourse)
}
