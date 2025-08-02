package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
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

// Lấy danh sách các event theo trạng thái, loại, tên, ...
// Tham số:
// - page: số trang
// - pageSize: số lượng event mỗi trang
// - title: tên event
// - status: trạng thái event
// - type: loại event
// - onlyRegisted: chỉ lấy event đã đăng ký
// - userID: id người dùng
// - onlyRegisted: chỉ lấy event đã đăng ký của người dùng
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

// Lấy danh sách các event đã đăng ký của người dùng
// Tham số:
// - page: số trang
// - pageSize: số lượng event mỗi trang
// - userID: id người dùng
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

// Lấy chi tiết event theo id
func (eventHandler *EventHandler) GetEventDetail(ctx *gin.Context) {

	eventID := ctx.Param("event_id")

	event, err := eventHandler.EventSer.GetEventByID(eventID)

	if eventHandler.isError(ctx, err) {
		return
	}

	response.Success(ctx, "Get event detail successfully", event)
}

// Điểm danh sự kiện
func (eventHandler *EventHandler) EventAttendance(ctx *gin.Context) {

	raw, _ := ctx.GetRawData()
	eventID := gjson.GetBytes(raw, "event_id").String()
	userID := middlewares.GetPrincipal(ctx)

	err := eventHandler.EventSer.EventAttendance(eventID, userID)

	if eventHandler.isError(ctx, err) {
		return
	}

	response.Success(ctx, "Attendance event successfully", nil)
}

// Đăng kí sự kiện
func (eventHandler *EventHandler) SubscribeEvent(ctx *gin.Context) {

	raw, _ := ctx.GetRawData()
	eventID := gjson.GetBytes(raw, "event_id").String()
	userID := middlewares.GetPrincipal(ctx)

	err := eventHandler.EventSer.SubscribeEvent(eventID, userID)
	if eventHandler.isError(ctx, err) {
		return
	}
	response.Success(ctx, "Subscribe event successfully", nil)
}

// Tạo sự kiện mới
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
		"id":        eventResponse.ID,
		"createdAt": eventResponse.CreatedAt,
	})
}

// Cập nhật sự kiện theo id
// Tham số:
// - event_id: id sự kiện
// - title: tên sự kiện
// - description: mô tả sự kiện
// - start_date: ngày bắt đầu sự kiện
// - end_date: ngày kết thúc sự kiện
// - location: vị trí sự kiện
// - type: loại sự kiện
// - status: trạng thái sự kiện
func (eventHandler *EventHandler) UpdateEvent(ctx *gin.Context) {
	var eventReq dtos.UpdateEventRequest

	if err := ctx.ShouldBindJSON(&eventReq); err != nil {
		response.Error(ctx, http.StatusBadRequest, "Invalid input")
		return
	}
	// Gọi service cập nhật event
	eventResponse, err := eventHandler.EventSer.UpdateEvent(&eventReq)
	if eventHandler.isError(ctx, err) {
		return
	}
	response.Success(ctx, "Update event successfully", eventResponse)
}

// Xóa sự kiện theo id
func (eventHandler *EventHandler) DeleteEvent(ctx *gin.Context) {
	eventID := ctx.Param("event_id")
	err := eventHandler.EventSer.DeleteEvent(eventID)
	if eventHandler.isError(ctx, err) {
		return
	}
	response.Success(ctx, "Deleted successfully", nil)
}

// Hủy đăng kí sự kiện
func (eventHandler *EventHandler) UnsubscribeEvent(ctx *gin.Context) {

	raw, _ := ctx.GetRawData()
	eventID := gjson.GetBytes(raw, "event_id").String()

	userID := middlewares.GetPrincipal(ctx)

	err := eventHandler.EventSer.UnsubscribeEvent(eventID, userID)
	if err != nil {
		response.Error(ctx, 400, "Unsubscribe event failed")
		return
	}

	response.Success(ctx, "Unsubscribe event successfully", nil)
}

// Lấy danh sách các user đã đăng ký sự kiện
// Tham số:
// - page: số trang
// - pageSize: số lượng user mỗi trang
// - eventID: id sự kiện
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

// Lấy danh sách các user đã điểm danh sự kiện
// Tham số:
// - page: số trang
// - pageSize: số lượng user mỗi trang
// - eventID: id sự kiện
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
