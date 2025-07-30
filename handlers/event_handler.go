package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"sfit-platform-web-backend/entities"
	"sfit-platform-web-backend/middlewares"
	"sfit-platform-web-backend/services"
	"sfit-platform-web-backend/utils/response"
	"strconv"
)

type EventHandler struct {
	*BaseHandler
	EventSer *services.EventService
}

func NewEventHandler(baseHandler *BaseHandler, eventSer *services.EventService) *EventHandler {
	return &EventHandler{
		BaseHandler: baseHandler,
		EventSer:    eventSer,
	}
}

func (eventHandler *EventHandler) GetEventList(ctx *gin.Context) {
	// Lấy tham số từ query
	page := ctx.Query("page")
	size := ctx.Query("pageSize")
	title := ctx.Query("title")
	status := ctx.Query("status")
	etype := ctx.Query("type")
	onlyRegisted := ctx.Query("onlyRegisted")
	var pageNum, pageSize = 1, 20
	olRegisted := false
	if page != "" {
		pageNum, _ = strconv.Atoi(page)
	}
	if size != "" {
		pageSize, _ = strconv.Atoi(size)
	}
	if onlyRegisted != "" {
		olRegisted = true
	}
	userID := middlewares.GetPrincipal(ctx)
	events, err := eventHandler.EventSer.GetEvents(pageNum, pageSize, title, etype, status, olRegisted, userID)
	if eventHandler.isError(ctx, err) {
		return
	}
	response.Success(ctx, events)
}

func (eventHandler *EventHandler) GetEventDetail(ctx *gin.Context) {
	eventID := ctx.Param("event_id")
	event, err := eventHandler.EventSer.GetEventByID(eventID)
	if eventHandler.isError(ctx, err) {
		return
	}
	response.Success(ctx, event)
}
func (eventHandler *EventHandler) GetRegistedEventList(ctx *gin.Context) {
	userID := middlewares.GetPrincipal(ctx)
	// Lấy tham số từ query
	page := ctx.Query("page")
	size := ctx.Query("pageSize")
	var pageNum, pageSize = 1, 20
	// Kiểm tra có phải là số
	if page != "" {
		pageNum, _ = strconv.Atoi(page)
	}
	if size != "" {
		pageSize, _ = strconv.Atoi(size)
	}
	events, err := eventHandler.EventSer.GetRegistedEvent(pageNum, pageSize, userID)
	if eventHandler.isError(ctx, err) {
		return
	}
	response.Success(ctx, events)
}

func (eventHandler *EventHandler) SubscribeEvent(ctx *gin.Context) {
	eventID := ctx.PostForm("event_id")
	userID := middlewares.GetPrincipal(ctx)
	err := eventHandler.EventSer.SubscribeEvent(eventID, userID)
	if eventHandler.isError(ctx, err) {
	}
}

func (eventHandler *EventHandler) UnsubscribeEvent(ctx *gin.Context) {
	eventID := ctx.PostForm("event_id")
	userID := middlewares.GetPrincipal(ctx)
	err := eventHandler.EventSer.UnsubscribeEvent(eventID, userID)
	if err != nil {
		response.Error(ctx, 400, "Unsubscribe event failed")
		return
	}
	response.Success(ctx, "Unsubscribe event successfully")
}
func (eventHandler *EventHandler) DeleteEvent(ctx *gin.Context) {
	eventID := ctx.Param("event_id")
	err := eventHandler.EventSer.DeleteEvent(eventID)
	if eventHandler.isError(ctx, err) {
		return
	}
	response.Success(ctx, "Deleted successfully")
}
func (eventHandler *EventHandler) UpdateEvent(ctx *gin.Context) {
	eventID := ctx.PostForm("event_id")
	title := ctx.PostForm("title")
	description := ctx.PostForm("description")
	eventType := ctx.PostForm("eventType")
	status := ctx.PostForm("status")

	eventAfterUpdate := entities.Event{
		ID:          uuid.MustParse(eventID),
		Title:       title,
		Description: description,
		EventType:   eventType,
		Status:      entities.EventStatus(status),
	}
	event, err := eventHandler.EventSer.UpdateEvent(&eventAfterUpdate)
	if eventHandler.isError(ctx, err) {
		return
	}
	response.Success(ctx, event)
}
