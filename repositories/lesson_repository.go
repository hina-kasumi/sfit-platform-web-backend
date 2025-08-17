package repositories

import (
	"sfit-platform-web-backend/entities"

	"gorm.io/gorm"
)

type LessonRepository struct {
	db *gorm.DB
}

func NewLessonRepository(db *gorm.DB) *LessonRepository {
	return &LessonRepository{
		db: db,
	}
}

func (repo *LessonRepository) CreateLesson(lesson *entities.Lesson) (*entities.Lesson, error) {
	if err := repo.db.Create(lesson).Error; err != nil {
		return nil, err
	}
	return lesson, nil
}

func (repo *LessonRepository) GetLessonByID(id string) (*entities.Lesson, error) {
	var lesson entities.Lesson
	if err := repo.db.First(&lesson, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &lesson, nil
}

func (repo *LessonRepository) DeleteLessonByID(id string) error {
	return repo.db.Delete(&entities.Lesson{}, "id = ?", id).Error
}

func (repo *LessonRepository) UpdateLesson(lesson *entities.Lesson) error {
	return repo.db.Save(lesson).Error
}
