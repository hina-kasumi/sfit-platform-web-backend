package repositories

import (
	"fmt"
	"sfit-platform-web-backend/entities"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ModuleRepository struct {
	db *gorm.DB
}

func NewModuleRepository(db *gorm.DB) *ModuleRepository {
	return &ModuleRepository{
		db: db,
	}
}

func (r *ModuleRepository) AddModuleToCourse(courseID string, moduleTitle string) (uuid.UUID, time.Time, error) {
	courseUUID, err := uuid.Parse(courseID)
	if err != nil {
		return uuid.Nil, time.Time{}, fmt.Errorf("invalid course ID format: %w", err)
	}
	if courseUUID == uuid.Nil {
		return uuid.Nil, time.Time{}, fmt.Errorf("course ID cannot be nil")
	}

	// Check if the course exists
	var dummy int
	if err := r.db.Model(&entities.Course{}).Select("1").Where("id = ?", courseID).Limit(1).Scan(&dummy).Error; err != nil {
		return uuid.Nil, time.Time{}, fmt.Errorf("failed to check course existence: %w", err)
	}
	if dummy == 0 {
		return uuid.Nil, time.Time{}, fmt.Errorf("course with ID %s does not exist", courseID)
	}

	// Create the module
	module := entities.Module{
		ID:        uuid.New(),
		CourseID:  courseUUID,
		Title:     moduleTitle,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := r.db.Create(&module).Error; err != nil {
		return uuid.Nil, time.Time{}, err
	}
	return module.ID, module.CreatedAt, nil
}

func (r *ModuleRepository) DeleteModule(moduleId uuid.UUID) error {
	result := r.db.Delete(&entities.Module{}, moduleId)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	if err := r.db.Where("module_id = ?", moduleId).Delete(&entities.Lesson{}).Error; err != nil {
		return err
	}
	return nil
}
