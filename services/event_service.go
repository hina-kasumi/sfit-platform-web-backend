package services

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"sfit-platform-web-backend/dtos"
	"sfit-platform-web-backend/entities"
	"sfit-platform-web-backend/repositories"
	"time"
)

type EventService struct {
	eventRepo *repositories.EventRepository
}

func NewEventService(eventRepo *repositories.EventRepository) *EventService {
	return &EventService{eventRepo: eventRepo}
}

// ==== Event Querry ====
func (eventSer *EventService) GetEvents(page int, size int, title string, etype string, status string, registed bool, userID string) ([]entities.Event, error) {
	return eventSer.eventRepo.GetEvents(page, size, title, etype, status, registed, userID)
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

// ==== Registration Logic =====
func (eventSer *EventService) UnsubscribeEvent(userID string, eventID string) error {
	event, err := eventSer.eventRepo.GetEventByID(eventID)
	if err != nil {
		return err
	}
	// Kiểm tra xem user đã đăng ký event hay chưa
	isRegisted, err := eventSer.eventRepo.CheckRegisted(userID, *event)
	if err != nil {
		return err
	}
	if !isRegisted {
		return errors.New("user is not registered for this event")
	}
	return eventSer.eventRepo.UnsubscribeEvent(userID, eventID)
}

func (eventSer *EventService) SubscribeEvent(userID string, eventID string) error {
	event, err := eventSer.eventRepo.GetEventByID(eventID)
	if err != nil {
		return err
	}
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
	return eventSer.eventRepo.SubscribeEvent(userID, eventID)
}

// ==== Attendance Logic ====
func (eventSer *EventService) EventAttendance(userID string, eventID string) error {
	event, err := eventSer.eventRepo.GetEventByID(eventID)
	if err != nil {
		return fmt.Errorf("event not found")
	}

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
	return eventSer.eventRepo.EventAttendance(userID, eventID)
}

// ==== Get User ====
func (eventSer *EventService) GetEventAttendance(page int, size int, eventID string) ([]entities.Users, error) {
	return eventSer.eventRepo.GetEventAttendance(page, size, eventID)
}
func (eventSer *EventService) GetEventRegisted(page int, size int, eventID string) ([]entities.Users, error) {
	return eventSer.eventRepo.GetEventRegisted(page, size, eventID)
}

func (eventSer *EventService) CheckRegisted(userID string, event entities.Event) (bool, error) {
	return eventSer.eventRepo.CheckRegisted(userID, event)
}
