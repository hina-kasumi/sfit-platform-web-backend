package entities

import (
	"time"

	"github.com/google/uuid"
)

type TeamMembers struct {
	UserID    uuid.UUID `gorm:"type:uuid;primaryKey;column:user_id"`
	TeamID    uuid.UUID `gorm:"type:uuid;primaryKey;column:team_id"`
	Role      RoleEnum  `gorm:"type:varchar(50);not null;"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
