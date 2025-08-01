package repositories

import (
	"gorm.io/gorm"
	"sfit-platform-web-backend/entities"
)

type TagRepository struct {
	db *gorm.DB
}

func NewTagRepository(db *gorm.DB) *TagRepository {
	return &TagRepository{db: db}
}

func (r *TagRepository) FindAll() ([]entities.Tag, error) {
	var tags []entities.Tag
	result := r.db.Find(&tags)
	return tags, result.Error
}
