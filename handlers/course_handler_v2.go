package handlers

import (
	"sfit-platform-web-backend/dtos"
	"sfit-platform-web-backend/middlewares"
	"sfit-platform-web-backend/utils/response"

	"github.com/gin-gonic/gin"
)

func (h *CourseHandler) GetCourseDetailByIDV2(c *gin.Context) {
	courseID := c.Param("course_id")
	userID := middlewares.GetPrincipal(c)
	course, err := h.courseService.GetCourseDetailByIDV2(userID, courseID)
	if h.isError(c, err) {
		return
	}
	response.Success(c, "get course detail success", course)
}

func (h *CourseHandler) GetListCourseV2(c *gin.Context) {
	userID := middlewares.GetPrincipal(c)
	var filter dtos.GetListCoursesForm
	if !h.canBindQuery(c, &filter) {
		return
	}

	courses, total, err := h.courseService.GetCourses(userID, filter)
	if h.isError(c, err) {
		return
	}
	response.GetListResp(c, "get list success", filter.Page, filter.PageSize, total, courses)
}
