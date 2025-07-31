package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
	"sfit-platform-web-backend/entities"
	"sfit-platform-web-backend/middlewares"
	"sfit-platform-web-backend/services"
	"sfit-platform-web-backend/utils/response"
	"strconv"
	"time"
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
	response.Success(ctx, "Get event list successfully", events)
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
	response.Success(ctx, "Get registed event list successfully", events)
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
	title := ctx.PostForm("title")
	eventType := ctx.PostForm("type")
	description := ctx.PostForm("description")
	priority, _ := strconv.ParseInt(ctx.PostForm("priority"), 10, 64)
	location := ctx.PostForm("location")
	maxPeople, _ := strconv.ParseInt(ctx.PostForm("max_people"), 10, 64)
	agency := ctx.PostForm("agency")
	status := ctx.PostForm("status")
	beginDate := ctx.PostForm("begin_at")
	endDate := ctx.PostForm("end_at")
	layout := "2006-01-02" // layout chuẩn của Go (yyyy-mm-dd)
	beginTime, err := time.Parse(layout, beginDate)
	if err != nil {
		log.Fatal(err)
	}

	endTime, err := time.Parse(layout, endDate)
	if err != nil {
		log.Fatal(err)
	}

	event := entities.Event{
		Title:       title,
		EventType:   eventType,
		Description: description,
		Priority:    int(priority),
		Location:    location,
		MaxPeople:   int(maxPeople),
		Agency:      agency,
		Status:      entities.EventStatus(status),
		BeginAt:     beginTime,
		EndAt:       endTime,
	}
	eventRespone, err := eventHandler.EventSer.CreateEvent(&event)
	if eventHandler.isError(ctx, err) {
		return
	}
	response.Success(ctx, "Create new event successfully", eventRespone)
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
	response.Success(ctx, "Update event successfully", event)
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
	response.Success(ctx, "Get event registed list successfully", users)
}

func (eventHandler *EventHandler) GetEventAttendanceList(ctx *gin.Context) {
	page, _ := strconv.ParseInt(ctx.Query("page"), 10, 64)
	size, _ := strconv.ParseInt(ctx.Query("pageSize"), 10, 64)
	eventID := ctx.Query("eventID")
	users, err := eventHandler.EventSer.GetEventAttendance(int(page), int(size), eventID)
	if eventHandler.isError(ctx, err) {
		return
	}
	response.Success(ctx, "Get event attendances list successfully", users)
}
