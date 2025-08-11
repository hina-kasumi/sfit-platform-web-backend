package handlers

import (
	"net/http"
	"sfit-platform-web-backend/dtos"
	"sfit-platform-web-backend/entities"
	"sfit-platform-web-backend/middlewares"
	"sfit-platform-web-backend/services"
	"sfit-platform-web-backend/utils/response"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type EventHandler struct {
	*BaseHandler
	EventSer   *services.EventService
	TagSer     *services.TagService
	TagTempSer *services.TagTempService
}

func NewEventHandler(baseHandler *BaseHandler, eventSer *services.EventService, tagSer *services.TagService, tagTempSer *services.TagTempService) *EventHandler {
	return &EventHandler{
		BaseHandler: baseHandler,
		EventSer:    eventSer,
		TagSer:      tagSer,
		TagTempSer:  tagTempSer,
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
	var rq dtos.ListEventReq
	if !eventHandler.canBindQuery(ctx, &rq) {
		return
	}

	userID := middlewares.GetPrincipal(ctx)

	events, total, err := eventHandler.EventSer.
		GetEvents(
			rq.Page, rq.PageSize, rq.Title, rq.Type, rq.Status, userID,
		)

	if eventHandler.isError(ctx, err) {
		return
	}
	type EventResponse struct {
		ID          string    `json:"id"`
		Title       string    `json:"title"`
		Description string    `json:"description"`
		Location    string    `json:"location"`
		BeginAt     time.Time `json:"begin_at"`
		Agency      string    `json:"agency"`
		MaxPeople   int       `json:"max_people"`
		Type        string    `json:"type"`
		Priority    int       `json:"priority"`
		Registed    bool      `json:"registed"`
		Tag         []string  `json:"tag"`
	}

	eventResponses := make([]EventResponse, 0, len(events))
	for _, e := range events {
		Registed, _ := eventHandler.EventSer.CheckRegisted(userID, e)
		eventResponses = append(eventResponses, EventResponse{
			ID:          e.ID.String(),
			Title:       e.Title,
			Description: e.Description,
			Location:    e.Location,
			BeginAt:     e.BeginAt, // hoặc format theo nhu cầu
			Agency:      e.Agency,
			MaxPeople:   e.MaxPeople,
			Type:        e.Type,
			Priority:    e.Priority,
			Registed:    Registed, //TODO: logic check user đã đăng ký hay chưa
			Tag:         eventHandler.TagTempSer.GetByEventOrCourse(e.ID.String(), ""),
		})
	}
	response.Success(ctx, "Get event list successfully", gin.H{
		"events":   eventResponses,
		"page":     rq.Page,
		"pageSize": rq.PageSize,
		"total":    total,
	})
}

// Lấy chi tiết event theo id
func (eventHandler *EventHandler) GetEventDetail(ctx *gin.Context) {
	eventID := ctx.Param("event_id")

	event, err := eventHandler.EventSer.GetEventByID(eventID)
	tags := eventHandler.TagTempSer.GetByEventOrCourse(eventID, "")

	if eventHandler.isError(ctx, err) {
		return
	}

	eventRp := dtos.EventDetailRp{
		Event: *event,
		Tags:  tags,
	}

	response.Success(ctx, "Get event detail successfully", eventRp)
}

// cập nhật trạng thái người dùng tham gia sự kiện
func (eventHandler *EventHandler) UpdateStatusUserAttendance(ctx *gin.Context) {
	eventID := ctx.Param("event_id")
	if eventID == "" {
		response.Error(ctx, http.StatusBadRequest, "Missing event_id")
		return
	}
	userID := ctx.Param("user_id")
	if userID == "" {
		response.Error(ctx, http.StatusUnauthorized, "Missing user_id")
		return
	}

	var updateReq dtos.UpdateUserAttendanceReq
	err := ctx.ShouldBindJSON(&updateReq)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "Invalid input")
		return
	}

	// Kiểm tra quyền hạn chỉ người dùng có vai trò Admin, Vice, Head mới được cập nhật trạng thái điểm danh
	if updateReq.Status == string(entities.Attended) &&
		!middlewares.HasRole(ctx,
			string(entities.RoleEnumAdmin),
			string(entities.RoleEnumVice),
			string(entities.RoleEnumHead)) {
		response.Error(ctx, http.StatusForbidden, "Forbidden")
		return
	}

	err = eventHandler.EventSer.UpdateStatusUserAttendance(userID, eventID, updateReq.Status)

	if eventHandler.isError(ctx, err) {
		return
	}

	response.Success(ctx, "Update user attendance status successfully", nil)
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
	// Tạo tag cho event
	eventHandler.TagSer.EnsureTags(eventReq.Tags)
	// Tạo TagTemp cho event
	eventHandler.TagTempSer.CreateTagTempEvent(eventResponse.ID, eventReq.Tags)

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

// Lấy danh sách các user đã đăng ký sự kiện
func (eventHandler *EventHandler) GetUsersInEvent(ctx *gin.Context) {
	eventID := ctx.Param("event_id")
	var query dtos.QueryUsersInEvent
	if !eventHandler.canBindQuery(ctx, &query) {
		return
	}

	users, total, err := eventHandler.EventSer.GetUsersInEvent(
		eventID, query.Page, query.PageSize, query.Status,
	)
	if eventHandler.isError(ctx, err) {
		return
	}
	type UserResponse struct {
		ID       uuid.UUID `json:"id"`
		Username string    `json:"username"`
		Email    string    `json:"email"`
	}

	userResponses := make([]UserResponse, 0, len(users))
	for _, u := range users {
		userResponses = append(userResponses, UserResponse{
			ID:       u.ID,
			Username: u.Username,
			Email:    u.Email,
		})
	}
	response.Success(ctx, "Get event registed list successfully", gin.H{
		"users":    userResponses,
		"page":     query.Page,
		"pageSize": query.PageSize,
		"total":    total,
	})
}
