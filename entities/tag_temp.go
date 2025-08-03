package entities

import "github.com/google/uuid"

type TagTemp struct {
	ID       uuid.UUID `gorm:"type:uuid;primaryKey"`
	TagID    string    `gorm:"type:varchar;column:tag_id"`
	EventID  uuid.UUID `gorm:"type:uuid;column:event_id"`
	CourseID uuid.UUID `gorm:"type:uuid;column:course_id"`
	Tag      Tag       `gorm:"foreignKey:TagID;references:ID;constraint:OnDelete:CASCADE"`
}
