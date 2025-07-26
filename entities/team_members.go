package entities

import (
	"time"

	"github.com/google/uuid"
)

type TeamRole string

const (
	RoleHeader     TeamRole = "HEADER"
	RoleViceHeader TeamRole = "VICE_HEADER"
	RoleMember     TeamRole = "MEMBER"
)

type TeamMembers struct {
	UserID    uuid.UUID `gorm:"type:uuid;primaryKey;column:user_id"`
	TeamID    uuid.UUID `gorm:"type:uuid;primaryKey;column:team_id"`
	Role      TeamRole  `gorm:"type:varchar(50);not null;check:role in ('HEADER', 'VICE_HEADER', 'MEMBER')"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
