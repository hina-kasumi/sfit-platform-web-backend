package entities

import (
	"time"

	"github.com/google/uuid"
)

type TaskAssignments struct {
	TaskID      uuid.UUID `gorm:"type:uuid;column:task_id;primaryKey"`
	UserID      uuid.UUID `gorm:"type:uuid;column:user_id;primaryKey"`
	IsCompleted bool      `gorm:"column:is_completed"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}
