package entities

import "github.com/google/uuid"

type UserRole struct {
	RoleID    RoleEnum  `gorm:"type:varchar;column:role_id;primaryKey"`
	UserID    uuid.UUID `gorm:"type:uuid;column:user_id;primaryKey"`
	CreatedAt uuid.Time `gorm:"column:create_at"`
}
