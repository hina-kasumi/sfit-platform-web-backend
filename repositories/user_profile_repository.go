package repositories

import (
	"sfit-platform-web-backend/entities"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserProfileRepository struct {
	db *gorm.DB
}

func NewUserProfileRepository(db *gorm.DB) *UserProfileRepository {
	return &UserProfileRepository{
		db: db,
	}
}

func (repo *UserProfileRepository) DeleteUser(userID uuid.UUID) error {
	tx := repo.db.Begin()

	if err := tx.Where("user_id = ?", userID).Delete(&entities.UserCourse{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Where("user_id = ?", userID).Delete(&entities.UserEvent{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Where("user_id = ?", userID).Delete(&entities.FavoriteCourse{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Where("user_id = ?", userID).Delete(&entities.LessonAttendance{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Where("user_id = ?", userID).Delete(&entities.EventAttendance{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Where("user_id = ?", userID).Delete(&entities.UserRate{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Where("user_id = ?", userID).Delete(&entities.UserProfile{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Where("id = ?", userID).Delete(&entities.Users{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (repo *UserProfileRepository) GetUserProfileByID(userID uuid.UUID) (*entities.UserProfile, error) {
	profile := entities.UserProfile{
		UserID: userID,
	}
	if err := repo.db.First(&profile).Error; err != nil {
		return nil, err
	}
	return &profile, nil
}

func (repo *UserProfileRepository) CreateUserProfile(profile *entities.UserProfile) (*entities.UserProfile, error) {
	if err := repo.db.Create(&profile).Error; err != nil {
		return nil, err
	}
	return profile, nil
}

func (repo *UserProfileRepository) UpdateUserProfile(profile *entities.UserProfile) (*entities.UserProfile, error) {
	if err := repo.db.Save(&profile).Error; err != nil {
		return nil, err
	}
	return profile, nil
}
