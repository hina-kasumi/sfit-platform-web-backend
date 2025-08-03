package services

import (
	"errors"
	"sfit-platform-web-backend/dtos"
	"sfit-platform-web-backend/entities"
	"sfit-platform-web-backend/repositories"
	"time"

	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserProfileService struct {
	userSer         *UserService
	userProfileRepo *repositories.UserProfileRepository
}

func NewUserProfileService(userProfileRepo *repositories.UserProfileRepository, userSer *UserService) *UserProfileService {
	return &UserProfileService{
		userSer:         userSer,
		userProfileRepo: userProfileRepo,
	}
}

func (profileSer *UserProfileService) UpdateUserProfile(profile *entities.UserProfile) (createAt, updateAt time.Time, err error) {
	existing, err := profileSer.userProfileRepo.GetUserProfileByID(profile.UserID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return time.Time{}, time.Time{}, errors.New("user profile not found")
		}
		return time.Time{}, time.Time{}, err
	}

	user, err := profileSer.userSer.GetUserByID(profile.UserID.String())
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	existing.FullName = profile.FullName
	existing.ClassName = profile.ClassName
	existing.Khoa = profile.Khoa
	existing.Phone = profile.Phone
	existing.Introduction = profile.Introduction
	existing.SocialLink = profile.SocialLink
	existing.UpdatedAt = time.Now()

	// Cập nhật email trong bảng Users
	if profile.Email != "" {
		existing.Email = profile.Email
		user.Email = profile.Email
		if _, err = profileSer.userSer.UpdateUser(user); err != nil {
			return time.Time{}, time.Time{}, err
		}
	}

	//gán các trường từ request vào bản cũ
	if _, err := profileSer.userProfileRepo.UpdateUserProfile(existing); err != nil {
		return time.Time{}, time.Time{}, err
	}

	return existing.CreatedAt, existing.UpdatedAt, nil
}

func (profileSer *UserProfileService) DeleteUser(userID uuid.UUID) error {
	return profileSer.userProfileRepo.DeleteUser(userID)
}

func (profileSer *UserProfileService) GetUserProfile(userID uuid.UUID) (*dtos.GetUserProfileResponse, error) {
	profile, err := profileSer.userProfileRepo.GetUserProfileByID(userID)
	if err != nil {
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
		Email:           profile.Email,
		CompletedCourse: 0,
		JoinedEvent:     0,
		CompletedTask:   0,
		SocialLink:      socialLink,
		CreatedAt:       profile.CreatedAt,
		UpdatedAt:       profile.UpdatedAt,
	}, nil

}

func (profileSer *UserProfileService) CreateUserProfile(profile *entities.UserProfile) error {
	_, result := profileSer.userProfileRepo.GetUserProfileByID(profile.UserID)

	if result == nil {
		return errors.New("profile already exists")
	}
	if !errors.Is(result, gorm.ErrRecordNotFound) {
		return result
	}

	user, err := profileSer.userSer.GetUserByID(profile.UserID.String())
	if err != nil {
		return err
	}
	profile.Email = user.Email

	_, err = profileSer.userProfileRepo.CreateUserProfile(profile)
	return err
}
