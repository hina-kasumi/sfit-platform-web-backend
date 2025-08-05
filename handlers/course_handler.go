package handlers

import (
	"net/http"
	"sfit-platform-web-backend/dtos"
	"sfit-platform-web-backend/services"
	"sfit-platform-web-backend/utils/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CourseHandler struct {
	BaseHandler   *BaseHandler
	CourseService *services.CourseService
}

func NewCourseHandler(baseHandler *BaseHandler, courseService *services.CourseService) *CourseHandler {
	return &CourseHandler{
		BaseHandler:   baseHandler,
		CourseService: courseService,
	}
}

// Lấy danh sách khóa học đã đăng ký của người dùng với phân trang
func (ch *CourseHandler) GetRegisteredCourses(c *gin.Context) {
	// Lấy userID từ URL parameter
	userID := c.Param("user_id")
	if userID == "" {
		response.Error(c, http.StatusBadRequest, "User ID is required")
		return
	}

	// Lấy page và pageSize từ query parameters
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		response.Error(c, http.StatusBadRequest, "Invalid page parameter")
		return
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 || pageSize > 100 {
		response.Error(c, http.StatusBadRequest, "Invalid page_size parameter")
		return
	}

	// Gọi service để lấy danh sách course
	result, err := ch.CourseService.GetRegisteredCourses(userID, page, pageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to get registered courses")
		return
	}

	response.Success(c, result)
}

// Lấy danh sách bài học của khóa học
func (ch *CourseHandler) GetCourseLessons(c *gin.Context) {
	// Lấy course_id từ URL parameter
	courseID := c.Param("course_id")
	if courseID == "" {
		response.Error(c, http.StatusBadRequest, "Course ID is required")
		return
	}

	// Lấy user_id từ query parameter
	userID := c.Query("user_id")
	if userID == "" {
		response.Error(c, http.StatusBadRequest, "User ID is required")
		return
	}

	// Gọi service để lấy thông tin lessons
	result, err := ch.CourseService.GetCourseLessons(courseID, userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to get course lessons")
		return
	}

	response.Success(c, result)
}

// Đánh giá khóa học
func (ch *CourseHandler) RateCourse(c *gin.Context) {
	// Lấy userID từ JWT token
	userID := c.GetString("user_id")
	if userID == "" {
		response.Error(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req dtos.CourseRateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Gọi service để rate course
	err := ch.CourseService.RateCourse(userID, req.Course, req.Star, req.Comment)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to rate course")
		return
	}

	response.Success(c, gin.H{"message": "Course rated successfully"})
}

// Xóa khóa học theo ID
func (ch *CourseHandler) DeleteCourse(c *gin.Context) {
	// Lấy course_id từ URL parameter
	courseID := c.Param("course_id")
	if courseID == "" {
		response.Error(c, http.StatusBadRequest, "Course ID is required")
		return
	}

	// Gọi service để delete course
	err := ch.CourseService.DeleteCourse(courseID)
	if err != nil {
		if err.Error() == "record not found" {
			response.Error(c, http.StatusNotFound, "Course not found")
			return
		}
		response.Error(c, http.StatusInternalServerError, "Failed to delete course")
		return
	}

	response.Success(c, gin.H{"message": "Course deleted successfully"})
}

// Lấy danh sách người dùng đã đăng ký khóa học với phân trang
func (ch *CourseHandler) GetRegisteredUsers(c *gin.Context) {
	// Lấy course_id từ query parameter
	courseID := c.Query("course_id")
	if courseID == "" {
		response.Error(c, http.StatusBadRequest, "Course ID is required")
		return
	}

	// Lấy page và pageSize từ query parameters
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("pageSize", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		response.Error(c, http.StatusBadRequest, "Invalid page parameter")
		return
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 || pageSize > 100 {
		response.Error(c, http.StatusBadRequest, "Invalid pageSize parameter")
		return
	}

	// Gọi service để lấy danh sách người dùng đã đăng ký
	result, err := ch.CourseService.GetRegisteredUsers(courseID, page, pageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to get registered users")
		return
	}

	response.Success(c, result)
}

// Đăng ký người dùng vào khóa học
func (ch *CourseHandler) RegisterUserToCourse(c *gin.Context) {
	// Lấy userID từ JWT token
	userID := c.GetString("user_id")
	if userID == "" {
		response.Error(c, http.StatusUnauthorized, "Unauthorized")
		return
	}
	var req dtos.CourseRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Gọi service để đăng ký người dùng vào khóa học
	err := ch.CourseService.RegisterUserToCourse(userID, req.CourseID)
	if err != nil {
		if err.Error() == "record not found" {
			response.Error(c, http.StatusNotFound, "Course not found")
			return
		}
		if err.Error() == "UNIQUE constraint failed" || err.Error() == "duplicated key not allowed" {
			response.Error(c, http.StatusConflict, "User already registered for this course")
			return
		}
		response.Error(c, http.StatusInternalServerError, "Failed to register for course")
		return
	}

	response.Success(c, gin.H{"message": "Successfully registered for course"})
}
