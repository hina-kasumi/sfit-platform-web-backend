package entities

import (
	"time"

	"github.com/google/uuid"
)

type Log struct {
	ID      uuid.UUID `gorm:"type:uuid;primaryKey"`
	At      time.Time `gorm:"autoCreateTime"`
	Type    string    `gorm:"type:varchar"`
	Message string    `gorm:"type:text"`
}
