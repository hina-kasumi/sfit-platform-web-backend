package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"sfit-platform-web-backend/dtos"
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
	response.Success(ctx, "Get event list successfully", gin.H{
		"events":   events,
		"page":     pageNum,
		"pageSize": pageSize,
		"total":    len(events),
	})
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
	response.Success(ctx, "Get registed event list successfully", gin.H{
		"events":   events,
		"page":     pageNum,
		"pageSize": pageSize,
		"total":    len(events),
	})
}

func (eventHandler *EventHandler) GetEventDetail(ctx *gin.Context) {
	eventID := ctx.Param("event_id")
	event, err := eventHandler.EventSer.GetEventByID(eventID)
	if eventHandler.isError(ctx, err) {
		return
	}
	response.Success(ctx, "Get event detail successfully", event)
}

func (eventHandler *EventHandler) EventAttendance(ctx *gin.Context) {
	eventID := ctx.PostForm("event_id")
	userID := middlewares.GetPrincipal(ctx)
	err := eventHandler.EventSer.EventAttendance(eventID, userID)
	if eventHandler.isError(ctx, err) {
		return
	}
	response.Success(ctx, "Attendance event successfully", nil)
}
func (eventHandler *EventHandler) SubscribeEvent(ctx *gin.Context) {
	eventID := ctx.PostForm("event_id")
	userID := middlewares.GetPrincipal(ctx)
	err := eventHandler.EventSer.SubscribeEvent(eventID, userID)
	if eventHandler.isError(ctx, err) {
	}
}

func (eventHandler *EventHandler) CreateEvent(ctx *gin.Context) {
	var eventReq dtos.NewEventRequest
	if err := ctx.ShouldBindJSON(&eventReq); err != nil {
		response.Error(ctx, http.StatusBadRequest, "Invalid input")
		return
	}
	// Gọi service tạo event
	eventResponse, err := eventHandler.EventSer.CreateEvent(&eventReq)
	if eventHandler.isError(ctx, err) {
		return
	}

	response.Success(ctx, "Create new event successfully", gin.H{
		"id":         eventResponse.ID,
		"created_at": eventResponse.CreatedAt,
	})
}

func (eventHandler *EventHandler) DeleteEvent(ctx *gin.Context) {
	eventID := ctx.Param("event_id")
	err := eventHandler.EventSer.DeleteEvent(eventID)
	if eventHandler.isError(ctx, err) {
		return
	}
	response.Success(ctx, "Deleted successfully", nil)
}
func (eventHandler *EventHandler) UpdateEvent(ctx *gin.Context) {
	var eventReq dtos.UpdateEventRequest

	if err := ctx.ShouldBindJSON(&eventReq); err != nil {
		response.Error(ctx, http.StatusBadRequest, "Invalid input")
		return
	}
	// Gọi service tạo event
	eventResponse, err := eventHandler.EventSer.UpdateEvent(&eventReq)
	if eventHandler.isError(ctx, err) {
		return
	}
	response.Success(ctx, "Update event successfully", eventResponse)
}
func (eventHandler *EventHandler) UnsubscribeEvent(ctx *gin.Context) {
	eventID := ctx.PostForm("event_id")
	userID := middlewares.GetPrincipal(ctx)
	err := eventHandler.EventSer.UnsubscribeEvent(eventID, userID)
	if err != nil {
		response.Error(ctx, 400, "Unsubscribe event failed")
		return
	}
	response.Success(ctx, "Unsubscribe event successfully", nil)
}

func (eventHandler *EventHandler) GetEventRegistedList(ctx *gin.Context) {
	page, _ := strconv.ParseInt(ctx.Query("page"), 10, 64)
	size, _ := strconv.ParseInt(ctx.Query("pageSize"), 10, 64)
	eventID := ctx.Query("eventID")
	users, err := eventHandler.EventSer.GetEventRegisted(int(page), int(size), eventID)
	if eventHandler.isError(ctx, err) {
		return
	}
	response.Success(ctx, "Get event registed list successfully", gin.H{
		"users":    users,
		"page":     page,
		"pageSize": size,
		"total":    len(users),
	})
}

func (eventHandler *EventHandler) GetEventAttendanceList(ctx *gin.Context) {
	page, _ := strconv.ParseInt(ctx.Query("page"), 10, 64)
	size, _ := strconv.ParseInt(ctx.Query("pageSize"), 10, 64)
	eventID := ctx.Query("eventID")
	users, err := eventHandler.EventSer.GetEventAttendance(int(page), int(size), eventID)
	if eventHandler.isError(ctx, err) {
		return
	}
	response.Success(ctx, "Get event attendances list successfully", gin.H{
		"users":    users,
		"page":     page,
		"pageSize": size,
		"total":    len(users),
	})
}
