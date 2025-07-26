package entities

import (
	"time"

	"github.com/google/uuid"
)

type FavoriteCourse struct {
	UserID    uuid.UUID `gorm:"type:uuid;primaryKey;column:user_id"`
	CourseID  uuid.UUID `gorm:"type:uuid;primaryKey;column:course_id"`
	CreatedAt time.Time `gorm:"column:create_at"`
}
