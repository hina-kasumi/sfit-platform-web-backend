package services

import (
	"errors"
	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"sfit-platform-web-backend/dtos"
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
			return time.Time{}, time.Time{}, errors.New("user profile not found")
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

func (profileSer *UserProfileService) GetUserProfile(userID uuid.UUID) (*dtos.GetUserProfileResponse, error) {
	var profile entities.UserProfile
	if err := profileSer.db.First(&profile, "user_id = ?", userID).Error; err != nil {
		return nil, err
	}

	var socialLink map[string]string
	_ = json.Unmarshal([]byte(profile.SocialLink), &socialLink)

	return &dtos.GetUserProfileResponse{
		UserID:          profile.UserID,
		FullName:        profile.FullName,
		ClassName:       profile.ClassName,
		Khoa:            profile.Khoa,
		Phone:           profile.Phone,
		Introduction:    profile.Introduction,
		CompletedCourse: 0,
		JoinedEvent:     0,
		CompletedTask:   0,
		SocialLink:      socialLink,
		CreatedAt:       profile.CreatedAt,
		UpdatedAt:       profile.UpdatedAt,
	}, nil

}

func (profileSer *UserProfileService) CreateUserProfile(profile *entities.UserProfile) error {
	var existing entities.UserProfile
	result := profileSer.db.First(&existing, "user_id = ?", profile.UserID)
	if result.Error == nil {
		return errors.New("profile already exists")
	}

	if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return result.Error
	}

	return profileSer.db.Create(profile).Error
}
