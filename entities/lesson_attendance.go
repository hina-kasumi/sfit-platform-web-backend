package entities

import (
	"time"

	"github.com/google/uuid"
)

type LessonAttendanceStatus string

const (
	Present         LessonAttendanceStatus = "present"
	AbsentExcused   LessonAttendanceStatus = "absent_excused"   // nghỉ có phép
	AbsentUnexcused LessonAttendanceStatus = "absent_unexcused" // nghỉ không phép
	Late            LessonAttendanceStatus = "late"
)

type LessonAttendance struct {
	UserID      uuid.UUID `gorm:"type:uuid;primaryKey;column:user_id"`
	LessonID    uuid.UUID `gorm:"type:uuid;primaryKey;column:lesson_id"`
	QuizPoint   int       `gorm:"column:quiz_point"`
	Duration    int       `gorm:"column:timestamp"`
	Status      string    `gorm:"type:varchar;check:status in ('present', 'absent', 'late')"`
	DeviceID    uuid.UUID `gorm:"type:uuid;column:device_id"`
	ModeratorID uuid.UUID `gorm:"type:uuid;column:moderator_id"`
	CreatedAt   time.Time `gorm:"column:create_at"`
	UpdatedAt   time.Time `gorm:"column:update_at"`
}
