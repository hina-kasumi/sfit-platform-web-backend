package entities

import (
	"time"

	"github.com/google/uuid"
)

type UserProfile struct {
	UserID       uuid.UUID `gorm:"type:uuid;primaryKey;column:user_id"`
	FullName     string    `gorm:"type:varchar(50)"`
	MSV          string    `gorm:"type:varchar(50)"`
	Avatar       string    `gorm:"type:varchar(255)"`
	CoverImage   string    `gorm:"type:varchar(255);column:cover_image"`
	ClassName    string    `gorm:"type:varchar(50);column:class_name"`
	Khoa         string    `gorm:"type:varchar(50)"`
	Phone        string    `gorm:"type:varchar(50)"`
	CreatedAt    time.Time `gorm:"column:create_at"`
	UpdatedAt    time.Time `gorm:"column:update_at"`
	Introduction string    `gorm:"type:varchar(255)"`
	Email        string    `gorm:"type:varchar(50);unique"`
	Location     string    `gorm:"type:varchar(50)"`
	SocialLink   string    `gorm:"type:json"`
}
