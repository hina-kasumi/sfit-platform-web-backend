package entities

import (
	"time"

	"github.com/google/uuid"
)

type Module struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey"`
	CourseID     uuid.UUID `gorm:"type:uuid;column:course_id;references:ID;constraint:OnDelete:CASCADE"`
	Title        string    `gorm:"type:varchar;column:module_title"`
	CreatedAt    time.Time `gorm:"column:create_at"`
	UpdatedAt    time.Time `gorm:"column:update_at"`
	TotalTime    int       `gorm:"column:total_time;default:0"` // đơn vị: giây hoặc phút
	TotalLessons int       `gorm:"column:total_lessons;default:0"`

	Lessons []Lesson `gorm:"foreignKey:ModuleID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
