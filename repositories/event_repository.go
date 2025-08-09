package repositories

import (
	"errors"
	"sfit-platform-web-backend/dtos"
	"sfit-platform-web-backend/entities"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EventRepository struct {
	db *gorm.DB
}

func NewEventRepository(db *gorm.DB) *EventRepository {
	return &EventRepository{
		db: db,
	}
}

// lấy danh sách event (lọc có điều kiện)
func (er *EventRepository) GetEvents(page int, size int, title string, etype string, status string, registed bool, userID string) ([]entities.Event, error) {
	var events []entities.Event
	query := er.db.Model(&entities.Event{})
	if title != "" {
		query = query.Where("title ILIKE ?", "%"+title+"%")
	}
	if etype != "" {
		query = query.Where("event_type = ?", etype)
	}
	if status != "" {
		query = query.Where("status ILIKE ?", status)
	}
	if registed {
		userID, _ := uuid.Parse(userID)
		query = query.Joins("JOIN user_events ue ON ue.event_id = events.id").
			Where("ue.user_id = ?", userID)
	}
	query = query.Order("begin_at, priority DESC")
	offset := (page - 1) * size
	result := query.Offset(offset).Limit(size).Find(&events)
	if result.Error != nil {
		return nil, result.Error
	}
	return events, nil
}

// user đã đăng kí event chưa
func (er *EventRepository) CheckRegisted(userID string, event entities.Event) (bool, error) {
	UserID, err := uuid.Parse(userID)
	if err != nil {
		return false, err
	}
	var userEvent entities.UserEvent
	result := er.db.Where("user_id = ? AND event_id = ?", UserID, event.ID).First(&userEvent)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false, nil
	}
	if result.Error != nil {
		return false, result.Error
	}
	return true, nil
}

// kiểm tra user đã điểm danh chưa
func (er *EventRepository) CheckAttendance(userID string, event entities.Event) (bool, error) {
	UserID, err := uuid.Parse(userID)
	if err != nil {
		return false, err
	}
	var eventAttendance entities.EventAttendance
	result := er.db.Where("user_id = ? AND event_id = ?", UserID, event.ID).First(&eventAttendance)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false, nil
	}
	if result.Error != nil {
		return false, result.Error
	}
	return true, nil
}

// đếm số lượng user đã đăng kí event
func (er *EventRepository) CountRegistedEvent(eventID string) (int64, error) {
	EventID, err := uuid.Parse(eventID)
	if err != nil {
		return 0, err
	}
	var count int64
	result := er.db.Model(&entities.UserEvent{}).Where("event_id = ?", EventID).Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}
	return count, nil
}

// lấy thông tin event theo id
func (er *EventRepository) GetEventByID(id string) (*entities.Event, error) {
	eventID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	event := entities.Event{ID: eventID}
	result := er.db.First(&event)
	if result.Error != nil {
		return nil, result.Error
	}
	return &event, nil
}

// tạo mới event
func (er *EventRepository) CreateEvent(eventReq *dtos.NewEventRequest) (*entities.Event, error) {
	eventID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	event := entities.Event{
		ID:          eventID,
		Title:       eventReq.Title,
		Type:        eventReq.Type,
		Description: eventReq.Description,
		Priority:    eventReq.Priority,
		MaxPeople:   eventReq.MaxPeople,
		Agency:      eventReq.Agency,
		EventType:   eventReq.EventType,
		Status:      eventReq.Status,
		BeginAt:     eventReq.BeginAt,
		EndAt:       eventReq.EndAt,
		Location:    eventReq.Location,
	}
	result := er.db.Create(&event)
	if result.Error != nil {
		return nil, result.Error
	}
	return &event, nil
}

// câp nhật event
func (er *EventRepository) UpdateEvent(eventReq *dtos.UpdateEventRequest) (*entities.Event, error) {
	var event = entities.Event{
		ID:          eventReq.ID,
		Title:       eventReq.Title,
		Type:        eventReq.Type,
		Description: eventReq.Description,
		Priority:    eventReq.Priority,
		MaxPeople:   eventReq.MaxPeople,
		Agency:      eventReq.Agency,
		EventType:   eventReq.EventType,
		Status:      eventReq.Status,
		BeginAt:     eventReq.BeginAt,
		EndAt:       eventReq.EndAt,
		Location:    eventReq.Location,
		UpdatedAt:   time.Now(),
	}
	result := er.db.Save(&event)
	if result.Error != nil {
		return nil, result.Error
	}
	return &event, nil
}

// xóa event
func (er *EventRepository) DeleteEvent(id string) error {
	eventID, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	event := entities.Event{
		ID: eventID,
	}
	result := er.db.Delete(&event)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("event not found")
	}
	return nil
}

// hủy đăng kí event
func (er *EventRepository) UnsubscribeEvent(userID string, eventID string) error {
	UserEvent := entities.UserEvent{
		UserID:  uuid.MustParse(userID),
		EventID: uuid.MustParse(eventID),
	}
	result := er.db.Delete(&UserEvent)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("event not found")
	}
	return nil
}

// đăng kí event
func (er *EventRepository) SubscribeEvent(userID string, eventID string) error {
	UserEvent := entities.UserEvent{
		UserID:  uuid.MustParse(userID),
		EventID: uuid.MustParse(eventID),
	}
	return er.db.Create(&UserEvent).Error
}

// điểm danh sự kiện
func (er *EventRepository) EventAttendance(userID string, eventID string) error {
	EventAttendance := entities.EventAttendance{
		UserID:  uuid.MustParse(userID),
		EventID: uuid.MustParse(eventID),
	}
	return er.db.Create(&EventAttendance).Error
}

// lấy danh sách những nguoi đã đăng kí vào sự kiện
func (er *EventRepository) GetEventRegisted(page int, size int, eventID string) ([]entities.Users, error) {
	EventID, err := uuid.Parse(eventID)
	if err != nil {
		return nil, err
	}
	var users []entities.Users
	query := er.db.Model(&entities.Users{}).Joins("JOIN user_events ue ON ue.user_id = users.id").
		Where("ue.event_id = ?", EventID)
	offset := (page - 1) * size
	result := query.Offset(offset).Limit(size).Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}

// lấy danh sách những người đã điểm danh vào sự kiện
func (er *EventRepository) GetEventAttendance(page int, size int, eventID string) ([]entities.Users, error) {
	EventID, err := uuid.Parse(eventID)
	if err != nil {
		return nil, err
	}
	var users []entities.Users
	query := er.db.Model(&entities.Users{}).Joins("JOIN event_attendances ea ON ea.user_id = users.id").
		Where("ea.event_id = ?", EventID)
	offset := (page - 1) * size
	result := query.Offset(offset).Limit(size).Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}
