package entities

import (
	"time"

	"github.com/google/uuid"
)

type UserEvent struct {
	UserID    uuid.UUID `gorm:"type:uuid;primaryKey;column:user_id"`
	EventID   uuid.UUID `gorm:"type:uuid;primaryKey;column:event_id"`
	CreatedAt time.Time `gorm:"column:create_at"`
	UpdatedAt time.Time `gorm:"column:update_at"`
}
