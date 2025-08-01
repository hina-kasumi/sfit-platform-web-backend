package services

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"sfit-platform-web-backend/entities"
	"time"
)

type UserProfileService struct {
	db *gorm.DB
}

func NewUserProfileService(db *gorm.DB) *UserProfileService {
	return &UserProfileService{db: db}
}

func (profileSer *UserProfileService) UpdateUserProfile(profile *entities.UserProfile) (createAt, updateAt time.Time, err error) {
	var existing entities.UserProfile
	result := profileSer.db.First(&existing, "user_id = ?", profile.UserID)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			profile.CreatedAt = time.Now()
			profile.UpdatedAt = time.Now()
			if err := profileSer.db.Create(profile).Error; err != nil {
				return time.Time{}, time.Time{}, err
			}
			return profile.CreatedAt, profile.UpdatedAt, nil
			//nếu chưa có bản ghi cũ, tạo mới
		}
		return time.Time{}, time.Time{}, result.Error
	}

	existing.FullName = profile.FullName
	existing.ClassName = profile.ClassName
	existing.Khoa = profile.Khoa
	existing.Phone = profile.Phone
	existing.Introduction = profile.Introduction
	existing.SocialLink = profile.SocialLink
	existing.UpdatedAt = time.Now()
	//gán các trường từ request vào bản cũ
	if err := profileSer.db.Save(&existing).Error; err != nil {
		return time.Time{}, time.Time{}, err
	}

	return existing.CreatedAt, existing.UpdatedAt, nil
}

func (profileSer *UserProfileService) DeleteUser(userID uuid.UUID) error {
	tx := profileSer.db.Begin()

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
