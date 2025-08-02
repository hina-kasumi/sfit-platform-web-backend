package services

import (
	"sfit-platform-web-backend/dtos"
	"sfit-platform-web-backend/entities"
	"sfit-platform-web-backend/repositories"
)

type EventService struct {
	eventRepo *repositories.EventRepository
}

func NewEventService(eventRepo *repositories.EventRepository) *EventService {
	return &EventService{eventRepo: eventRepo}
}
func (eventSer *EventService) GetEvents(page int, size int, title string, etype string, status string, registed bool, userID string) ([]entities.Event, error) {
	return eventSer.eventRepo.GetEvents(page, size, title, etype, status, registed, userID)
}
func (eventSer *EventService) GetRegistedEvent(page int, size int, userID string) ([]entities.Event, error) {
	return eventSer.eventRepo.GetRegistedEvent(page, size, userID)
}
func (eventSer *EventService) GetEventByID(id string) (*entities.Event, error) {
	return eventSer.eventRepo.GetEventByID(id)
}
func (eventSer *EventService) CreateEvent(eventReq *dtos.NewEventRequest) (*entities.Event, error) {
	return eventSer.eventRepo.CreateEvent(eventReq)
}

func (eventSer *EventService) UpdateEvent(eventReq *dtos.UpdateEventRequest) (*entities.Event, error) {
	return eventSer.eventRepo.UpdateEvent(eventReq)
}

func (eventSer *EventService) DeleteEvent(id string) error {
	return eventSer.eventRepo.DeleteEvent(id)
}

func (eventSer *EventService) UnsubscribeEvent(userID string, eventID string) error {
	return eventSer.eventRepo.UnsubscribeEvent(userID, eventID)
}

func (eventSer *EventService) SubscribeEvent(userID string, eventID string) error {
	return eventSer.eventRepo.SubscribeEvent(userID, eventID)
}
func (eventSer *EventService) EventAttendance(userID string, eventID string) error {
	return eventSer.eventRepo.EventAttendance(userID, eventID)
}

func (eventSer *EventService) GetEventAttendance(page int, size int, eventID string) ([]entities.Users, error) {
	return eventSer.eventRepo.GetEventAttendance(page, size, eventID)
}
func (eventSer *EventService) GetEventRegisted(page int, size int, eventID string) ([]entities.Users, error) {
	return eventSer.eventRepo.GetEventRegisted(page, size, eventID)
}
