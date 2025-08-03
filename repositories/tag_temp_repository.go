package repositories

import (
	"fmt"
	"sfit-platform-web-backend/entities"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TagTempRepository struct {
	db *gorm.DB
}

func NewTagTempRepository(db *gorm.DB) *TagTempRepository {
	return &TagTempRepository{
		db: db,
	}
}

func (r *TagTempRepository) CreateNewTagTemp(tagID string, courseID uuid.UUID) (*entities.TagTemp, error) {
	tagTemp := &entities.TagTemp{
		ID:       uuid.New(),
		TagID:    tagID,
		CourseID: courseID,
	}

	// Lưu TagTemp vào database
	if err := r.db.Create(tagTemp).Error; err != nil {
		return nil, fmt.Errorf("failed to create tag temp: %w", err)
	}

	return tagTemp, nil
}