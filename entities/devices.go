package entities

import "github.com/google/uuid"

type Device struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name         string    `gorm:"type:varchar"`
	AccessToken  string    `gorm:"type:text"`
	WifiName     string    `gorm:"type:text"`
	WifiPassword string    `gorm:"type:text"`
}
