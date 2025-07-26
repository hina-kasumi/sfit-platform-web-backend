package entities

import (
	"time"

	"github.com/google/uuid"
)

type UserCourse struct {
	UserID    uuid.UUID `gorm:"type:uuid;primaryKey;column:user_id"`
	CourseID  uuid.UUID `gorm:"type:uuid;primaryKey;column:course_id"`
	CreatedAt time.Time `gorm:"column:create_at"`
	UpdatedAt time.Time `gorm:"column:update_at"`
}
