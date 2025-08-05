package entities

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type CourseLevel string

const (
	Beginner     CourseLevel = "Beginner"
	Intermediate CourseLevel = "Intermediate"
	Advanced     CourseLevel = "Advanced"
)

type Course struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey"`
	Title       string         `gorm:"type:varchar"`
	Description string         `gorm:"type:text"`
	Target      pq.StringArray `gorm:"type:text[]"`
	Require     pq.StringArray `gorm:"type:text[]"`
	Teachers    pq.StringArray `gorm:"type:text[]"`
	Language    string         `gorm:"type:varchar"`
	Certificate bool           `gorm:"type:boolean"`
	Level       string         `gorm:"type:varchar;check:level in ('Beginner', 'Intermediate', 'Advanced')"`
	CreatedAt   time.Time      `gorm:"column:create_at"`
	UpdatedAt   time.Time      `gorm:"column:update_at"`
}
