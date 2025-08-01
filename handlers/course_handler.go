package handlers

import (
    "sfit-platform-web-backend/dtos"
    "sfit-platform-web-backend/services"
    "sfit-platform-web-backend/utils/response"
    "time"

    "github.com/gin-gonic/gin"
)

type CourseHandler struct {
    *BaseHandler
    courseSer *services.CourseService
    tagSer    *services.TagService
}

func NewCourseHandler(baseHandler *BaseHandler, courseSer *services.CourseService, tagSer *services.TagService) *CourseHandler {
    return &CourseHandler{
        BaseHandler: baseHandler,
        courseSer:   courseSer,
        tagSer:      tagSer,
    }
}

func (h *CourseHandler) CreateCourse(ctx *gin.Context) {
    var req dtos.CreateCourseRequest
    if !h.canBindJSON(ctx, &req) {
        return
    }

    tags, err := h.tagSer.EnsureTags(req.Tags)
    if h.isError(ctx, err) {
        return
    }

    courseID, createdAt, err := h.courseSer.CreateCourse(
        req.Title, req.Description, req.Type, req.Target, req.Require,
        req.Teachers, req.Language, req.Certificate, req.Level, tags,
    )
    if h.isError(ctx, err) {
        return
    }

    // Ok
    resp := dtos.CreateCourseResponse{
        ID:        courseID.String(),
        CreatedAt: createdAt.Format(time.RFC3339),
    }
    response.Success(ctx, "Course created successfully", resp)
}