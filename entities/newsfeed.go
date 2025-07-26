package entities

import (
	"time"

	"github.com/google/uuid"
)

type Newsfeed struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	Title       string    `gorm:"type:varchar"`
	Description string    `gorm:"type:varchar"`
	Type        string    `gorm:"type:varchar"`
	CreatedAt   time.Time `gorm:"column:create_at"`
	UpdatedAt   time.Time `gorm:"column:update_at"`
	CreatorID   uuid.UUID `gorm:"type:uuid;not null"`
	Creator     Users     `gorm:"foreignKey:CreatorID;references:ID"`
}
