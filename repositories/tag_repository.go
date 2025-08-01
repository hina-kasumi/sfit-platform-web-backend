package repositories

import (
	"sfit-platform-web-backend/entities"

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

func (r *TagRepository) FindByID(id string) (*entities.Tag, error) {
    var tag entities.Tag
    err := r.db.Where("id = ?", id).First(&tag).Error
    if err != nil {
        return nil, err
    }
    return &tag, nil
}

// Tạo tag mới
func (r *TagRepository) Create(tag *entities.Tag) error {
    return r.db.Create(tag).Error
}