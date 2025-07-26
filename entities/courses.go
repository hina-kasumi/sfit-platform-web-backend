package entities

import (
	"time"

	"github.com/google/uuid"
)

type CourseLevel string

const (
	Beginner     CourseLevel = "Beginner"
	Intermediate CourseLevel = "Intermediate"
	Advanced     CourseLevel = "Advanced"
)

type Course struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	Title       string    `gorm:"type:varchar"`
	Description string    `gorm:"type:text"`
	Target      []string  `gorm:"type:text[]"`
	Require     []string  `gorm:"type:text[]"`
	Teachers    []string  `gorm:"type:text[]"`
	Language    string    `gorm:"type:varchar"`
	Certificate bool      `gorm:"type:boolean"`
	Level       string    `gorm:"type:varchar;check:level in ('Beginner', 'Intermediate', 'Advanced')"`
	CreatedAt   time.Time `gorm:"column:create_at"`
	UpdatedAt   time.Time `gorm:"column:update_at"`
}
