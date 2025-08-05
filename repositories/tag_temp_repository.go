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

func (r *TagTempRepository) UpdateTagTemp(courseID string, tags []entities.Tag) error {
	cid, err := uuid.Parse(courseID)
	if err != nil {
		return fmt.Errorf("invalid course ID: %w", err)
	}

	// 1. Xoá toàn bộ TagTemp cũ của khoá học này
	if err := r.db.Where("course_id = ?", cid).Delete(&entities.TagTemp{}).Error; err != nil {
		return fmt.Errorf("failed to delete old tag temps: %w", err)
	}

	// 2. Tạo danh sách TagTemp mới từ []Tag
	var newTagTemps []entities.TagTemp
	for _, tag := range tags {
		newTagTemps = append(newTagTemps, entities.TagTemp{
			ID:       uuid.New(),
			TagID:    tag.ID,
			CourseID: cid,
		})
	}

	// 3. Thêm mới vào DB (nếu có tag)
	if len(newTagTemps) > 0 {
		if err := r.db.Create(&newTagTemps).Error; err != nil {
			return fmt.Errorf("failed to create new tag temps: %w", err)
		}
	}

	return nil
}

