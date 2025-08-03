package handlers

import (
	"log"
	"net/http"
	"sfit-platform-web-backend/dtos"
	"sfit-platform-web-backend/services"
	"sfit-platform-web-backend/utils/response"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CourseHandler struct {
    *BaseHandler
    course_ser   *services.CourseService
    tag_ser      *services.TagService
    tagTemp_ser  *services.TagTempService
}

func NewCourseHandler(baseHandler *BaseHandler, course_ser *services.CourseService, tag_ser *services.TagService, tagTemp_ser *services.TagTempService) *CourseHandler {
    return &CourseHandler{
        BaseHandler: baseHandler,
        course_ser:   course_ser,
        tag_ser:      tag_ser,
        tagTemp_ser:  tagTemp_ser,
    }
}

func (h *CourseHandler) CreateCourse(ctx *gin.Context) {
    var req dtos.CreateCourseRequest
    if !h.canBindJSON(ctx, &req) {
        return
    }

    tags, err := h.tag_ser.EnsureTags(req.Tags)
    if h.isError(ctx, err) {
        return
    }

    courseID, createdAt, err := h.course_ser.CreateCourse(
        req.Title, req.Description, req.Type, req.Target, req.Require,
        req.Teachers, req.Language, req.Certificate, req.Level,
    )
    if h.isError(ctx, err) {
        return
    }

    // Tạo TagTemp với tag, course
    for i, tag := range tags {
        _, err := h.tagTemp_ser.CreateTagTemp(tag.ID, courseID)
        if err != nil {
            // Log error với context
            log.Printf("Failed to create TagTemp for tag %s (index %d) and course %s: %v", 
                tag.ID, i, courseID.String(), err)
            
            // Sử dụng helper method có sẵn
            if h.isError(ctx, err) {
                return
            }
        }
    }

    // Ok
    resp := dtos.CreateCourseResponse{
        ID:        courseID.String(),
        CreatedAt: createdAt.Format(time.RFC3339),
    }
    response.Success(ctx, "Course created successfully", resp)
}

func (h *CourseHandler) GetListCourse(ctx *gin.Context) {
    // Lấy user_id từ context (đã được set bởi JWT middleware)
	userIDInterface, exists := ctx.Get("user_id")
	var userID uuid.UUID
	
	if exists {
		switch v := userIDInterface.(type) {
		case string:
			parsedID, err := uuid.Parse(v)
			if err != nil {
				response.Error(ctx, http.StatusBadRequest, "Invalid user ID format")
				return
			}
			userID = parsedID
		case uuid.UUID:
			userID = v
		default:
			// Nếu không có user_id, có thể để uuid.Nil cho guest user
			userID = uuid.Nil
		}
	} else {
		// Không có user_id trong context, có thể là guest user
		userID = uuid.Nil
	}



    //page=number&pageSize=number&title=string&onlyRegisted=boolean&type=string&level=string
    page := ctx.Query("page")
    pageSize := ctx.Query("pageSize")
    title := ctx.Query("title")
    onlyRegisted := ctx.Query("onlyRegisted")
    courseType := ctx.Query("type")
    level := ctx.Query("level")

    // Validate and parse query parameters
    if page == "" || pageSize == "" {
        response.Error(ctx, 400, "Page and pageSize are required")
        return
    }
    pageNum, err := strconv.Atoi(page)
    if err != nil || pageNum < 1 {
        response.Error(ctx, 400, "Invalid page number")
        return
    }
    pageSizeNum, err := strconv.Atoi(pageSize)
    if err != nil || pageSizeNum < 1 {
        response.Error(ctx, 400, "Invalid page size")
        return
    }
    // Call service to get courses
    courses, err := h.course_ser.GetListCourse(userID.String(), pageNum, pageSizeNum, title, onlyRegisted == "true", courseType, level)
    if err != nil { 
        response.Error(ctx, http.StatusInternalServerError, "Failed to get courses")
        return
    }
    // Return response
    response.Success(ctx, "Courses retrieved successfully", courses)
}