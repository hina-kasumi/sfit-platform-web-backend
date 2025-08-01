package handlers

import (
    "sfit-platform-web-backend/services"
)

type TagHandler struct {
    *BaseHandler
    tagSer *services.TagService
}

func NewTagHandler(baseHandler *BaseHandler, tagSer *services.TagService) *TagHandler {
    return &TagHandler{
        BaseHandler: baseHandler,
        tagSer:   tagSer,
    }
}