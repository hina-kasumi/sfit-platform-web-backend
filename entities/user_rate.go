package entities

import (
	"time"

	"github.com/google/uuid"
)

type UserRate struct {
	UserID    uuid.UUID `gorm:"type:uuid;primaryKey;column:user_id"`
	CourseID  uuid.UUID `gorm:"type:uuid;primaryKey;column:courses_id"`
	Star      int       `gorm:"type:int;check:star >= 1 and star <= 5"`
	Comment   string    `gorm:"type:text"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:update_at"`
}
