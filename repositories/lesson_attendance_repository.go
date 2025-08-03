package repositories

import (
	"gorm.io/gorm"
)

type LessonAttendanceRepository struct {
	db *gorm.DB
}

func NewLessonAttendanceRepository(db *gorm.DB) *LessonAttendanceRepository {
	return &LessonAttendanceRepository{
		db: db,
	}
}