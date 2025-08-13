package handlers

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"sfit-platform-web-backend/dtos"
	"sfit-platform-web-backend/middlewares"
	"sfit-platform-web-backend/services"
	"sfit-platform-web-backend/utils/response"

	"github.com/gin-gonic/gin"
)

type CourseHandler struct {
	*BaseHandler
	courseService  *services.CourseService
	tagService     *services.TagService
	tagTempService *services.TagTempService
}

func NewCourseHandler(base *BaseHandler, course *services.CourseService, tag *services.TagService, tagTemp *services.TagTempService) *CourseHandler {
	return &CourseHandler{
		BaseHandler:     base,
		courseService:   course,
		tagService:      tag,
		tagTempService:  tagTemp,
	}
}

// POST /courses
func (h *CourseHandler) CreateCourse(ctx *gin.Context) {
	var req dtos.CreateCourseRequest
	if !h.canBindJSON(ctx, &req) {
		return
	}

	tags, err := h.tagService.EnsureTags(req.Tags)
	if h.isError(ctx, err) {
		return
	}

	courseID, createdAt, err := h.courseService.CreateCourse(
		req.Title, req.Description, req.Type, req.Target,
		req.Require, req.Teachers, req.Language, req.Certificate, req.Level,
	)
	if h.isError(ctx, err) {
		return
	}

	for i, tag := range tags {
		if _, err := h.tagTempService.CreateTagTemp(tag.ID, courseID); err != nil {
			log.Printf("Failed to create TagTemp for tag %s (index %d) and course %s: %v",
				tag.ID, i, courseID.String(), err)
			if h.isError(ctx, err) {
				return
			}
		}
	}

	resp := dtos.CreateCourseResponse{
		ID:        courseID.String(),
		CreatedAt: createdAt.Format(time.RFC3339),
	}
	response.Success(ctx, "Course created successfully", resp)
}

// GET /courses
func (h *CourseHandler) GetListCourse(ctx *gin.Context) {
	// userID := utils.GetUserIDFromContext(ctx)
	userID := middlewares.GetPrincipal(ctx)

	pageNum, pageSizeNum, valid := parsePagination(ctx)
	if !valid {
		return
	}

	title := ctx.Query("title")
	onlyRegisted := ctx.Query("onlyRegisted") == "true"
	courseType := ctx.Query("type")
	level := ctx.Query("level")

	courses, err := h.courseService.GetListCourse(
		userID, pageNum, pageSizeNum, title, onlyRegisted, courseType, level,
	)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, "Failed to get courses")
		return
	}

	response.Success(ctx, "Courses retrieved successfully", courses)
}

// GET /courses/:course_id
func (h *CourseHandler) GetCourseDetailByID(ctx *gin.Context) {
	// userID := utils.GetUserIDFromContext(ctx)
	userID := middlewares.GetPrincipal(ctx)
	courseID := ctx.Param("course_id")

	if courseID == "" {
		response.Error(ctx, http.StatusBadRequest, "Course ID is required")
		return
	}

	courseDetail, err := h.courseService.GetCourseDetailByID(courseID, userID)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, "Failed to get course detail")
		return
	}

	response.Success(ctx, "Course detail retrieved successfully", courseDetail)
}


// POST /courses/favourite
func (h *CourseHandler) MarkCourseAsFavourite(ctx *gin.Context) {
	// userID := utils.GetUserIDFromContext(ctx)
	userID := middlewares.GetPrincipal(ctx)
	if userID == "" {
		response.Error(ctx, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req dtos.SetFavouriteCourseRequest
	if !h.canBindJSON(ctx, &req) {
		return
	}

	if err := h.courseService.MarkCourseAsFavourite(userID, req.CourseID); err != nil {
		response.Error(ctx, http.StatusInternalServerError, "Failed to mark course as favourite")
		return
	}

	response.Success(ctx, "Course marked as favourite successfully", nil)
}

// DELETE /courses/favourite
func (h *CourseHandler) UnmarkCourseAsFavourite(ctx *gin.Context) {
	// userID := utils.GetUserIDFromContext(ctx)
	userID := middlewares.GetPrincipal(ctx)
	if userID == "" {
		response.Error(ctx, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req dtos.SetFavouriteCourseRequest
	if !h.canBindJSON(ctx, &req) {
		return
	}

	if err := h.courseService.UnmarkCourseAsFavourite(userID, req.CourseID); err != nil {
		response.Error(ctx, http.StatusInternalServerError, "Failed to unmark course as favourite")
		return
	}

	response.Success(ctx, "Course unmarked as favourite successfully", nil)
}

// PUT /courses
func (h *CourseHandler) UpdateCourse(ctx *gin.Context) {
	userID := middlewares.GetPrincipal(ctx)
	if userID == "" {
		response.Error(ctx, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req dtos.UpdateCourseRequest
	if !h.canBindJSON(ctx, &req) {
		return
	}

	if req.ID == "" {
		response.Error(ctx, http.StatusBadRequest, "Course ID is required")
		return
	}

	tags, err := h.tagService.EnsureTags(req.Tags)
	if h.isError(ctx, err) {
		return
	}

	updatedAt, err := h.courseService.UpdateCourse(
		req.ID, req.Title, req.Description, req.Type, req.Target,
		req.Require, req.Teachers, req.Language, req.Certificate, req.Level,
	)
	if h.isError(ctx, err) {
		return
	}

	err = h.tagTempService.UpdateTagTemp(req.ID, tags)
	if h.isError(ctx, err) {
		return
	}

	resp := dtos.UpdateCourseResponse{
		// ID:        req.ID,
		UpdatedAt: updatedAt.Format(time.RFC3339),
	}
	response.Success(ctx, "Course update successfully", resp)
}

//
// Helper
//

func parsePagination(ctx *gin.Context) (page, pageSize int, ok bool) {
	pageStr := ctx.Query("page")
	pageSizeStr := ctx.Query("pageSize")

	if pageStr == "" || pageSizeStr == "" {
		response.Error(ctx, http.StatusBadRequest, "Page and pageSize are required")
		return 0, 0, false
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		response.Error(ctx, http.StatusBadRequest, "Invalid page number")
		return 0, 0, false
	}

	pageSize, err = strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		response.Error(ctx, http.StatusBadRequest, "Invalid page size")
		return 0, 0, false
	}

	return page, pageSize, true
}

// Lấy danh sách khóa học đã đăng ký của người dùng với phân trang
func (ch *CourseHandler) GetRegisteredCourses(c *gin.Context) {
	// Lấy userID từ JWT token
	userID := middlewares.GetPrincipal(c)
	if userID == "" {
		response.Error(c, http.StatusUnauthorized, "Unauthorized")
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
	result, err := ch.courseService.GetRegisteredCourses(userID, page, pageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to get registered courses")
		return
	}

	response.Success(c, "Get registered courses successfully", result)
}

// Lấy danh sách bài học của khóa học
func (ch *CourseHandler) GetCourseLessons(c *gin.Context) {
	// Lấy course_id từ URL parameter
	courseID := c.Param("course_id")
	if courseID == "" {
		response.Error(c, http.StatusBadRequest, "Course ID is required")
		return
	}

	// Lấy userID từ JWT token
	userID := middlewares.GetPrincipal(c)
	if userID == "" {
		response.Error(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Gọi service để lấy thông tin lessons
	result, err := ch.courseService.GetCourseLessons(courseID, userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to get course lessons")
		return
	}

	response.Success(c, "Get course lessons successfully", result)
}

// Đánh giá khóa học
func (ch *CourseHandler) RateCourse(c *gin.Context) {

	// Lấy userID từ JWT token
	userID := middlewares.GetPrincipal(c)
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
	err := ch.courseService.RateCourse(userID, req.Course, req.Star, req.Comment)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to rate course")
		return
	}

	response.Success(c, "Course rated successfully", nil)
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
	err := ch.courseService.DeleteCourse(courseID)
	if err != nil {
		if err.Error() == "record not found" {
			response.Error(c, http.StatusNotFound, "Course not found")
			return
		}
		response.Error(c, http.StatusInternalServerError, "Failed to delete course")
		return
	}

	response.Success(c, "Course deleted successfully", nil)
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
	result, err := ch.courseService.GetRegisteredUsers(courseID, page, pageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to get registered users")
		return
	}

	response.Success(c, "Get registered users successfully", result)
}

// Đăng ký người dùng vào khóa học
func (ch *CourseHandler) RegisterUserToCourse(c *gin.Context) {
	// Lấy userID từ JWT token
	userID := middlewares.GetPrincipal(c)
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
	err := ch.courseService.RegisterUserToCourse(userID, req.CourseID)
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

	response.Success(c, "Successfully registered for course", nil)
}


