package entities

import (
	"time"

	"github.com/google/uuid"
)

type EventStatus string

const (
	StatusDraft     EventStatus = "DRAFT"
	StatusUpcoming  EventStatus = "UPCOMING"
	StatusOngoing   EventStatus = "ONGOING"
	StatusCompleted EventStatus = "COMPLETED"
	StatusCancelled EventStatus = "CANCELLED"
)

type Event struct {
	ID          uuid.UUID   `gorm:"type:uuid;primaryKey"`
	Title       string      `gorm:"type:varchar"`
	Type        string      `gorm:"type:varchar"`
	Description string      `gorm:"type:varchar"`
	Priority    int         `gorm:"type:int"`
	Location    string      `gorm:"type:varchar"`
	MaxPeople   int         `gorm:"column:max_people"`
	EventType   string      `gorm:"column:event_type;type:varchar"`
	Agency      string      `gorm:"type:varchar"`
	Status      EventStatus `gorm:"type:varchar;not null;check:status in ('DRAFT', 'UPCOMING', 'ONGOING', 'COMPLETED', 'CANCELLED')"`
	BeginAt     time.Time   `gorm:"column:begin_at"`
	EndAt       time.Time   `gorm:"column:end_at"`
	CreatedAt   time.Time   `gorm:"column:create_at"`
	UpdatedAt   time.Time   `gorm:"column:update_at"`
}
