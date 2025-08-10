package dtos

import (
	"sfit-platform-web-backend/entities"
	"time"

	"github.com/google/uuid"
)

type ListEventReq struct {
	Page     int    `form:"page" binding:"required"`
	PageSize int    `form:"pageSize" binding:"required"`
	Title    string `form:"title"`
	Type     string `form:"type"`
	Status   string `form:"status"`
}

type QueryUsersInEvent struct {
	Page     int    `form:"page" binding:"required"`
	PageSize int    `form:"pageSize" binding:"required"`
	Status   string `form:"status" binding:"required"`
}

type UpdateUserAttendanceReq struct {
	Status string `json:"status" binding:"required"`
}

type NewEventRequest struct {
	Title       string               `json:"title" binding:"required"`
	Type        string               `json:"type" binding:"required"`
	Description string               `json:"description" binding:"required"`
	Priority    int                  `json:"priority" binding:"required"`
	Location    string               `json:"location" binding:"required"`
	MaxPeople   int                  `json:"max_people" binding:"required"`
	Agency      string               `json:"agency" binding:"required"`
	Status      entities.EventStatus `json:"status" binding:"required"`
	BeginAt     time.Time            `json:"begin_at" binding:"required"`
	EndAt       time.Time            `json:"end_at" binding:"required"`
	Tags        []string             `json:"tags"`
}

type UpdateEventRequest struct {
	ID          uuid.UUID            `json:"id"`
	Title       string               `json:"title"`
	Type        string               `json:"type"`
	Description string               `json:"description"`
	Priority    int                  `json:"priority"`
	Location    string               `json:"location"`
	MaxPeople   int                  `json:"max_people"`
	Agency      string               `json:"agency"`
	Status      entities.EventStatus `json:"status"`
	BeginAt     time.Time            `json:"begin_at"`
	EndAt       time.Time            `json:"end_at"`
}

type EventDetailRp struct {
	entities.Event
	Tags []string `json:"tags"`
}
