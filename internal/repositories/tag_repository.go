package repositories

import (
	"sfit-platform-web-backend/internal/model"

	"gorm.io/gorm"
)

type TagRepository struct {
	db *gorm.DB
}

func NewTagRepository(db *gorm.DB) *TagRepository {
	return &TagRepository{
		db: db,
	}
}

func (r *TagRepository) FindByID(id string) (*model.Tag, error) {
	var tag model.Tag
	err := r.db.Where("id = ?", id).First(&tag).Error
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

// Tạo tag mới
func (r *TagRepository) CreateNewTag(tag *model.Tag) error {
	return r.db.Create(tag).Error
}

func (r *TagRepository) FindAll() ([]model.Tag, error) {
	var tags []model.Tag
	result := r.db.Find(&tags)
	return tags, result.Error
}
