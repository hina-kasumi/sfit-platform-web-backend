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
	ID           uuid.UUID      `gorm:"type:uuid;primaryKey"`
	Title        string         `gorm:"type:varchar"`
	Description  string         `gorm:"type:text"`
	Type         string         `gorm:"type:varchar"`
	Target       pq.StringArray `gorm:"type:text[]"`
	Require      pq.StringArray `gorm:"type:text[]"`
	Teachers     pq.StringArray `gorm:"type:text[]"`
	Language     string         `gorm:"type:varchar"`
	TotalTime    int            `gorm:"column:total_time;default:0"` // đơn vị: giây hoặc phút
	Rate         float32        `gorm:"column:rate;default:0"`
	TotalLessons int            `gorm:"column:total_lessons;default:0"`
	Certificate  bool           `gorm:"type:boolean"`
	Level        string         `gorm:"type:varchar;check:level in ('Beginner', 'Intermediate', 'Advanced')"`
	CreatedAt    time.Time      `gorm:"column:create_at"`
	UpdatedAt    time.Time      `gorm:"column:update_at"`

	Modules     []Module     `gorm:"foreignKey:CourseID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	UsersCourse []UserCourse `gorm:"foreignKey:CourseID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Tags        []TagTemp    `gorm:"foreignKey:CourseID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
