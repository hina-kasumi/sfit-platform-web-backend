package entities

import (
	"time"

	"github.com/google/uuid"
)

type LessonAttendanceStatus string

const (
	Present         LessonAttendanceStatus = "present"          // đi học hoặc làm thành công
	AbsentExcused   LessonAttendanceStatus = "absent_excused"   // nghỉ có phép hoặc không học
	AbsentUnexcused LessonAttendanceStatus = "absent_unexcused" // nghỉ không phép hoặc không làm có phép
	Late            LessonAttendanceStatus = "late"             // làm muộn hoặc đi học muộn
)

type LessonAttendance struct {
	UserID      uuid.UUID              `gorm:"type:uuid;primaryKey;column:user_id"`
	LessonID    uuid.UUID              `gorm:"type:uuid;primaryKey;column:lesson_id"`
	QuizPoint   *int                   `gorm:"column:quiz_point"`
	Duration    *int                   `gorm:"column:duration"`
	Status      LessonAttendanceStatus `gorm:"type:varchar"`
	DeviceID    *uuid.UUID             `gorm:"type:uuid;column:device_id"`
	ModeratorID uuid.UUID              `gorm:"type:uuid;column:moderator_id"` // ID của người điều hành
	CreatedAt   time.Time              `gorm:"column:create_at"`
	UpdatedAt   time.Time              `gorm:"column:update_at"`
}
