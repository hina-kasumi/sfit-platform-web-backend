package entities

import (
	"time"

	"github.com/google/uuid"
)

type EventAttendanceStatus string

const (
	Registered EventAttendanceStatus = "REGISTERED"
	Attended   EventAttendanceStatus = "ATTENDED"
)

type EventAttendance struct {
	UserID    uuid.UUID             `gorm:"type:uuid;primaryKey;column:user_id"`
	EventID   uuid.UUID             `gorm:"type:uuid;primaryKey;column:event_id"`
	Status    EventAttendanceStatus `gorm:"type:varchar;column:status"`
	CreatedAt time.Time             `gorm:"column:create_at"`
	UpdatedAt time.Time             `gorm:"column:update_at"`
}
