package repositories

import (
	"sfit-platform-web-backend/internal/model"

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

	if err := tx.Where("user_id = ?", userID).Delete(&model.UserCourse{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Where("user_id = ?", userID).Delete(&model.FavoriteCourse{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Where("user_id = ?", userID).Delete(&model.LessonAttendance{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Where("user_id = ?", userID).Delete(&model.EventAttendance{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Where("user_id = ?", userID).Delete(&model.UserRate{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Where("user_id = ?", userID).Delete(&model.UserProfile{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Where("id = ?", userID).Delete(&model.Users{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (repo *UserProfileRepository) GetUserProfileByID(userID uuid.UUID) (*model.UserProfile, error) {
	profile := model.UserProfile{
		UserID: userID,
	}
	if err := repo.db.First(&profile).Error; err != nil {
		return nil, err
	}
	return &profile, nil
}

func (repo *UserProfileRepository) CreateUserProfile(profile *model.UserProfile) (*model.UserProfile, error) {
	if err := repo.db.Create(&profile).Error; err != nil {
		return nil, err
	}
	return profile, nil
}

func (repo *UserProfileRepository) UpdateUserProfile(profile *model.UserProfile) (*model.UserProfile, error) {
	if err := repo.db.Save(&profile).Error; err != nil {
		return nil, err
	}
	return profile, nil
}

func (ur *UserProfileRepository) GetUsersByMsvs(msvs []string) ([]model.UserProfile, error) {
	var users []model.UserProfile
	result := ur.db.Where("msv IN ?", msvs).Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}

func (ur *UserProfileRepository) InitAdminProfile(userID uuid.UUID, email string) error {
	profile := &model.UserProfile{
		UserID:     userID,
		Email:      email,
		SocialLink: "{}",
		FullName:   email,
	}
	return ur.db.Create(profile).Error
}
