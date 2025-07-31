package dtos

import (
	"github.com/google/uuid"
	"sfit-platform-web-backend/entities"
	"time"
)

type NewEventRequest struct {
	Title       string               `gorm:"type:varchar" json:"title"`
	Type        string               `gorm:"type:varchar" json:"type"`
	Description string               `gorm:"type:varchar" json:"description"`
	Priority    int                  `gorm:"type:int" json:"priority"`
	Location    string               `gorm:"type:varchar" json:"location"`
	MaxPeople   int                  `gorm:"column:max_people" json:"max_people"`
	EventType   string               `gorm:"column:event_type;type:varchar"`
	Agency      string               `gorm:"type:varchar"`
	Status      entities.EventStatus `gorm:"type:varchar;not null;check:status in ('DRAFT', 'UPCOMING', 'ONGOING', 'COMPLETED', 'CANCELLED')"`
	BeginAt     time.Time            `gorm:"column:begin_at" json:"begin_at"`
	EndAt       time.Time            `gorm:"column:end_at" json:"end_at"`
	CreatedAt   time.Time            `gorm:"column:create_at"`
	UpdatedAt   time.Time            `gorm:"column:update_at"`
}

type UpdateEventRequest struct {
	ID          uuid.UUID            `gorm:"type:uuid;primaryKey" json:"id"`
	Title       string               `gorm:"type:varchar" json:"title"`
	Type        string               `gorm:"type:varchar" json:"type"`
	Description string               `gorm:"type:varchar" json:"description"`
	Priority    int                  `gorm:"type:int" json:"priority"`
	Location    string               `gorm:"type:varchar" json:"location"`
	MaxPeople   int                  `gorm:"column:max_people" json:"max_people"`
	EventType   string               `gorm:"column:event_type;type:varchar"`
	Agency      string               `gorm:"type:varchar"`
	Status      entities.EventStatus `gorm:"type:varchar;not null;check:status in ('DRAFT', 'UPCOMING', 'ONGOING', 'COMPLETED', 'CANCELLED')"`
	BeginAt     time.Time            `gorm:"column:begin_at" json:"begin_at"`
	EndAt       time.Time            `gorm:"column:end_at" json:"end_at"`
}
