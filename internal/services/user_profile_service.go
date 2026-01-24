package services

import (
	"context"
	"errors"
	"log"
	"sfit-platform-web-backend/internal/dtos"
	"sfit-platform-web-backend/internal/model"
	"sfit-platform-web-backend/internal/repositories"
	"time"

	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type UserProfileService struct {
	redisClient     *redis.Client
	ctx             context.Context
	userSer         *UserService
	userProfileRepo *repositories.UserProfileRepository
	eventSer        *EventService
	courseSer       *CourseService
	taskSer         *TaskService
}

func NewUserProfileService(userProfileRepo *repositories.UserProfileRepository,
	userSer *UserService,
	eventSer *EventService,
	courseSer *CourseService,
	taskSer *TaskService,
	redisClient *redis.Client,
	ctx context.Context,
) *UserProfileService {
	return &UserProfileService{
		redisClient:     redisClient,
		ctx:             ctx,
		userSer:         userSer,
		eventSer:        eventSer,
		courseSer:       courseSer,
		taskSer:         taskSer,
		userProfileRepo: userProfileRepo,
	}
}

func (profileSer *UserProfileService) UpdateUserProfile(profile *model.UserProfile) (createAt, updateAt time.Time, err error) {
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
	existing.Avatar = profile.Avatar
	existing.CoverImage = profile.CoverImage
	existing.MSV = profile.MSV
	existing.Location = profile.Location
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
	// try to get from cache
	cachedProfile, err := profileSer.getCachedUserProfile(userID)
	if err == nil && cachedProfile != nil {
		return cachedProfile, nil
	}
	// if not found in cache, get from database

	profile, err := profileSer.userProfileRepo.GetUserProfileByID(userID)
	if err != nil {
		return nil, err
	}

	var socialLink map[string]string
	_ = json.Unmarshal([]byte(profile.SocialLink), &socialLink)

	_, joinedEvents, _ := profileSer.eventSer.GetEvents(1, 10, "", "", "", string(model.Attended), profile.UserID.String())

	isCompleted := true
	_, completedTasks, _ := profileSer.taskSer.ListTasksByUserID(profile.UserID.String(), 1, 10, &isCompleted)

	listCourse, _ := profileSer.courseSer.GetCourseUserCompletion(profile.UserID)
	lenCourse := len(listCourse)

	result := &dtos.GetUserProfileResponse{
		UserID:          profile.UserID,
		Avatar:          profile.Avatar,
		CoverImage:      profile.CoverImage,
		FullName:        profile.FullName,
		ClassName:       profile.ClassName,
		Khoa:            profile.Khoa,
		Phone:           profile.Phone,
		Introduction:    profile.Introduction,
		Email:           profile.Email,
		CompletedCourse: int64(lenCourse),
		JoinedEvent:     joinedEvents,
		CompletedTask:   completedTasks,
		SocialLink:      socialLink,
		Location:        profile.Location,
		Msv:             profile.MSV,
		CreatedAt:       profile.CreatedAt,
		UpdatedAt:       profile.UpdatedAt,
	}
	if err := profileSer.setCachedUserProfile(profile.UserID, result); err != nil {
		log.Printf("Failed to cache user profile: %v", err)
	}
	return result, nil
}

func (profileSer *UserProfileService) buildCacheKey(userID uuid.UUID) string {
	return "user_profile:" + userID.String()
}

// cache user profile
func (profileSer *UserProfileService) getCachedUserProfile(userID uuid.UUID) (*dtos.GetUserProfileResponse, error) {
	cacheKey := profileSer.buildCacheKey(userID)
	cachedData, err := profileSer.redisClient.Get(profileSer.ctx, cacheKey).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Cache miss
		}
		return nil, err
	}

	var profile dtos.GetUserProfileResponse
	if err := json.Unmarshal([]byte(cachedData), &profile); err != nil {
		return nil, err
	}

	return &profile, nil
}

// set user profile to cache
func (profileSer *UserProfileService) setCachedUserProfile(userID uuid.UUID, profile *dtos.GetUserProfileResponse) error {
	cacheKey := profileSer.buildCacheKey(userID)
	data, err := json.Marshal(profile)
	if err != nil {
		return err
	}

	// Set cache with an expiration time of 10 minutes
	return profileSer.redisClient.Set(profileSer.ctx, cacheKey, data, 10*time.Minute).Err()
}

func (profileSer *UserProfileService) CreateUserProfile(profile *model.UserProfile) error {
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
