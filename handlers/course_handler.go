package handlers

import (
	"log"
	"sfit-platform-web-backend/dtos"
	"sfit-platform-web-backend/services"
	"sfit-platform-web-backend/utils/response"
	"time"

	"github.com/gin-gonic/gin"
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

    // Create TagTemp for each tag
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