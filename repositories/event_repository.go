package repositories

import (
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"sfit-platform-web-backend/entities"
)

type EventRepository struct {
	db *gorm.DB
}

func NewEventRepository(db *gorm.DB) *EventRepository {
	return &EventRepository{
		db: db,
	}
}
func (er *EventRepository) GetEvents(page int, size int, title string, etype string, status string, registed bool, userID string) ([]entities.Event, error) {
	var events []entities.Event
	query := er.db.Model(&entities.Event{})

	// Lọc theo title
	if title != "" {
		query = query.Where("title ILIKE ?", "%"+title+"%")
	}

	// Lọc theo type
	if etype != "" {
		query = query.Where("event_type = ?", etype)
	}

	// Lọc theo status
	if status != "" {
		query = query.Where("status ILIKE ?", status)

	}
	// Lọc theo sự kiện đã đăng ký
	if registed {
		userID, _ := uuid.Parse(userID)
		query = query.Joins("JOIN user_events ue ON ue.event_id = events.id").
			Where("ue.user_id = ?", userID)
	}

	// Phân trang
	offset := (page - 1) * size
	result := query.Offset(offset).Limit(size).Find(&events)
	if result.Error != nil {
		return nil, result.Error
	}
	return events, nil
}
func (er *EventRepository) GetRegistedEvent(page int, size int, userID string) ([]entities.Event, error) {
	result, err := er.GetEvents(page, size, "", "", "", true, userID)
	if err != nil {
		return nil, err
	}
	return result, nil
}
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
func (er *EventRepository) CreateEvent(event *entities.Event) (*entities.Event, error) {
	result := er.db.Create(&event)
	if result.Error != nil {
		return nil, result.Error
	}
	return event, nil
}

func (er *EventRepository) UpdateEvent(event *entities.Event) (*entities.Event, error) {
	result := er.db.Save(&event)
	if result.Error != nil {
		return nil, result.Error
	}
	return event, nil
}

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

func (er *EventRepository) SubscribeEvent(userID string, eventID string) error {
	UserEvent := entities.UserEvent{
		UserID:  uuid.MustParse(userID),
		EventID: uuid.MustParse(eventID),
	}
	result := er.db.Create(&UserEvent)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (er *EventRepository) EventAttendance(userID string, eventID string) error {
	EventAttendance := entities.EventAttendance{
		UserID:  uuid.MustParse(userID),
		EventID: uuid.MustParse(eventID),
	}
	result := er.db.Create(&EventAttendance)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (er *EventRepository) GetEventRegisted(page int, size int, eventID string) ([]entities.Users, error) {
	EventID, _ := uuid.Parse(eventID)
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
func (er *EventRepository) GetEventAttendance(page int, size int, eventID string) ([]entities.Users, error) {
	EventID, _ := uuid.Parse(eventID)
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
