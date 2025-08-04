package handlers

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"sfit-platform-web-backend/dtos"
	"sfit-platform-web-backend/services"
	"sfit-platform-web-backend/utils"
	"sfit-platform-web-backend/utils/response"
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
	userID := utils.GetUserIDFromContext(ctx)

	pageNum, pageSizeNum, valid := parsePagination(ctx)
	if !valid {
		return
	}

	title := ctx.Query("title")
	onlyRegisted := ctx.Query("onlyRegisted") == "true"
	courseType := ctx.Query("type")
	level := ctx.Query("level")

	courses, err := h.courseService.GetListCourse(
		userID.String(), pageNum, pageSizeNum, title, onlyRegisted, courseType, level,
	)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, "Failed to get courses")
		return
	}

	response.Success(ctx, "Courses retrieved successfully", courses)
}

// GET /courses/:course_id
func (h *CourseHandler) GetCourseDetailByID(ctx *gin.Context) {
	userID := utils.GetUserIDFromContext(ctx)
	courseID := ctx.Param("course_id")

	if courseID == "" {
		response.Error(ctx, http.StatusBadRequest, "Course ID is required")
		return
	}

	courseDetail, err := h.courseService.GetCourseDetailByID(courseID, userID.String())
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, "Failed to get course detail")
		return
	}

	response.Success(ctx, "Course detail retrieved successfully", courseDetail)
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
