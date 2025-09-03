package services

import (
	"errors"
	"fmt"
	"sfit-platform-web-backend/dtos"
	"sfit-platform-web-backend/entities"
	"sfit-platform-web-backend/repositories"
	"time"

	"github.com/google/uuid"
)

type EventService struct {
	eventRepo *repositories.EventRepository
}

func NewEventService(eventRepo *repositories.EventRepository) *EventService {
	return &EventService{eventRepo: eventRepo}
}

// ==== Event Querry ====
func (eventSer *EventService) GetEvents(page int, size int, title string, eventType string, status string, userEventStatus string, userID string) ([]entities.Event, int64, error) {
	return eventSer.eventRepo.GetEvents(page, size, title, eventType, status, userEventStatus, userID)
}

func (eventSer *EventService) GetEventByID(id string) (*entities.Event, error) {
	return eventSer.eventRepo.GetEventByID(id)
}

// ==== Event CRUD ====
func (eventSer *EventService) CreateEvent(eventReq *dtos.NewEventRequest) (*entities.Event, error) {
	if eventReq.EndAt.Before(eventReq.BeginAt) {
		return nil, errors.New("end time must be after begin time")
	}
	return eventSer.eventRepo.CreateEvent(eventReq)
}

func (eventSer *EventService) UpdateEvent(eventReq *dtos.UpdateEventRequest) (*entities.Event, error) {
	if eventReq.ID == uuid.Nil {
		return nil, errors.New("event id is required")
	}
	if eventReq.EndAt.Before(eventReq.BeginAt) {
		return nil, errors.New("end time must be after begin time")
	}
	return eventSer.eventRepo.UpdateEvent(eventReq)
}

func (eventSer *EventService) DeleteEvent(id string) error {
	return eventSer.eventRepo.DeleteEvent(id)
}

func (eventSer *EventService) UpdateStatusUserAttendance(userID, eventID, status string) error {
	var err error
	event, err := eventSer.eventRepo.GetEventByID(eventID)
	if err != nil {
		return err
	}

	switch status {
	case string(entities.Registered):
		// Kiểm tra xem user đã đăng ký event hay chưa
		isRegisted, err := eventSer.eventRepo.CheckRegisted(userID, *event)
		if err != nil {
			return err
		}
		if isRegisted {
			return errors.New("user is already registered for this event")
		}
		// Kiểm tra số người đăng ký đã vượt quá số lượng người đăng ký cho event hay chưa
		count, err := eventSer.eventRepo.CountRegistedEvent(eventID)
		if err != nil {
			return err
		}
		if int(count) >= event.MaxPeople {
			return errors.New("event is full")
		}
		return eventSer.eventRepo.UpdateStatusUserAttendance(
			userID, eventID, entities.Registered,
		)
	case string(entities.Attended):
		now := time.Now()
		if now.Before(event.BeginAt) || now.After(event.EndAt) {
			return fmt.Errorf("event is not in progress")
		}
		// Kiểm tra xem user đã điểm danh chưa
		isAttendance, err := eventSer.eventRepo.CheckAttendance(userID, *event)
		if err != nil {
			return fmt.Errorf("event not found")
		}
		if isAttendance {
			return fmt.Errorf("user has already attended this event")
		}
		return eventSer.eventRepo.UpdateStatusUserAttendance(userID, eventID, entities.Attended)
	case "UNREGISTERED":
		err = eventSer.eventRepo.DeleteEventAttendance(userID, eventID)
	default:
		err = fmt.Errorf("invalid status")
	}

	return err
}

// ==== Get User ====
func (eventSer *EventService) GetUsersInEvent(eventID string, page int, size int, status string) ([]entities.Users, int64, error) {
	return eventSer.eventRepo.GetUsersInEvent(page, size, eventID, status)
}

func (eventSer *EventService) CheckRegisted(userID string, event entities.Event) (bool, error) {
	return eventSer.eventRepo.CheckRegisted(userID, event)
}
