package entities

import (
	"time"

	"github.com/google/uuid"
)

type Task struct {
	ID              uuid.UUID `gorm:"type:uuid;primaryKey"`
	EventID         uuid.UUID `gorm:"type:uuid;column:event_id"`
	Name            string    `gorm:"type:varchar"`
	Description     string    `gorm:"type:text"`
	StartDate       time.Time `gorm:"column:start_date"`
	Deadline        time.Time `gorm:"column:dead_line"`
	Assignee        string    `gorm:"type:varchar"`
	PercentComplete int       `gorm:"column:percent_complete"`
	CreatedAt       time.Time `gorm:"column:create_at"`
	UpdatedAt       time.Time `gorm:"column:update_at"`
}
