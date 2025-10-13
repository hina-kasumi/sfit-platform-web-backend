package entities

import (
	"time"

	"github.com/google/uuid"
)

type UserCourseStatus string

const (
	UserCourseStatusRequest UserCourseStatus = "REQUEST"
	UserCourseStatusLearn   UserCourseStatus = "LEARN"
	UserCourseStatusBlocked UserCourseStatus = "BLOCKED"
)

type UserCourse struct {
	UserID    uuid.UUID        `gorm:"type:uuid;primaryKey;column:user_id"`
	CourseID  uuid.UUID        `gorm:"type:uuid;primaryKey;column:course_id"`
	Status    UserCourseStatus `gorm:"column:status"`
	CreatedAt time.Time        `gorm:"column:create_at"`
	UpdatedAt time.Time        `gorm:"column:update_at"`
}
