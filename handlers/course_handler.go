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
		BaseHandler:    base,
		courseService:  course,
		tagService:     tag,
		tagTempService: tagTemp,
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

// GET /course/list-user-complete
func (h *CourseHandler) GetListUserCompleteCourse(ctx *gin.Context) {
	page, pageSize, valid := parsePagination(ctx)
	if !valid {
		return
	}

	courseID := ctx.Query("course_id")
	if courseID == "" {
		response.Error(ctx, http.StatusBadRequest, "Course ID is required")
		return
	}

	listUser, err := h.courseService.GetListUserCompleteCourse(
		courseID, page, pageSize,
	)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, "Failed to get user complete courses")
		return
	}

	response.Success(ctx, "Courses retrieved successfully", listUser)
}

// GET /course/module
func (h *CourseHandler) AddModuleToCourse(ctx *gin.Context) {
	var req dtos.AddModuleToCourseRequest
	if !h.canBindJSON(ctx, &req) {
		return
	}

	if req.CourseID == "" || req.ModuleTitle == "" {
		response.Error(ctx, http.StatusBadRequest, "Course ID and Module Title are required")
		return
	}

	moduleID, create_at, err := h.courseService.AddModuleToCourse(req.CourseID, req.ModuleTitle)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, "Failed to add module to course")
		return
	}
	addResponse := dtos.AddModuleToCourseResponse{
		ModuleID:    moduleID.String(),
		CourseID:    req.CourseID,
		ModuleTitle: req.ModuleTitle,
		CreatedAt:   create_at.Format(time.RFC3339),
	}
	response.Success(ctx, "Module added to course successfully", addResponse)
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
