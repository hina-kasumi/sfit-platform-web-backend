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
func (r *TagRepository) CreateNewTag(tag *entities.Tag) error {
    return r.db.Create(tag).Error
}package repositories

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
