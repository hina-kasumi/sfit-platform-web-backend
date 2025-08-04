package entities

import (
	"time"

	"github.com/google/uuid"
)

type UserProfile struct {
	UserID       uuid.UUID `gorm:"type:uuid;primaryKey;column:user_id"`
	FullName     string    `gorm:"type:varchar"`
	MSV          string    `gorm:"type:varchar"`
	ClassName    string    `gorm:"type:varchar;column:class_name"`
	Khoa         string    `gorm:"type:varchar"`
	Phone        string    `gorm:"type:varchar"`
	CreatedAt    time.Time `gorm:"column:create_at"`
	UpdatedAt    time.Time `gorm:"column:update_at"`
	Introduction string    `gorm:"type:varchar"`
	Email        string    `gorm:"type:varchar;unique"`
	Location     string    `gorm:"type:varchar"`
	SocialLink   string    `gorm:"type:json"`
}
