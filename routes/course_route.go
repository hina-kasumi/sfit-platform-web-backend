package routes

import (
    "sfit-platform-web-backend/handlers"
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
}