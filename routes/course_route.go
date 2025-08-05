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
    router.GET("/course", r.handler.GetListCourse)
    router.GET("/course/:course_id", r.handler.GetCourseDetailByID)
    router.POST("/course/favourite", r.handler.MarkCourseAsFavourite)
    router.DELETE("/course/favourite", r.handler.UnmarkCourseAsFavourite)
}