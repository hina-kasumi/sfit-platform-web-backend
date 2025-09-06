package dtos

import (
	"time"

	"github.com/google/uuid"
)

type UpdateUserProfileRequest struct {
	Avatar       string            `json:"avatar"`
	CoverImage   string            `json:"cover_image"`
	FullName     string            `json:"full_name"`
	ClassName    string            `json:"class_name"`
	Khoa         string            `json:"khoa"`
	Phone        string            `json:"phone"`
	Email        string            `json:"email"`
	Introduction string            `json:"introduction"`
	SocialLink   map[string]string `json:"social_link"`
}

type GetUserProfileResponse struct {
	UserID          uuid.UUID         `json:"user_id"`
	Avatar          string            `json:"avatar"`
	CoverImage      string            `json:"cover_image"`
	FullName        string            `json:"full_name"`
	ClassName       string            `json:"class_name"`
	Khoa            string            `json:"khoa"`
	Phone           string            `json:"phone"`
	Email           string            `json:"email"`
	Introduction    string            `json:"introduction"`
	CompletedCourse int64             `json:"completed_course"`
	JoinedEvent     int64             `json:"joined_event"`
	CompletedTask   int64             `json:"completed_task"`
	SocialLink      map[string]string `json:"social_link"`
	CreatedAt       time.Time         `json:"created_at"`
	UpdatedAt       time.Time         `json:"updated_at"`
}

type CreateUserProfileRequest struct {
	FullName        string            `json:"full_name"`
	ClassName       string            `json:"class_name"`
	Avatar          string            `json:"avatar"`
	CoverImage      string            `json:"cover_image"`
	Khoa            string            `json:"khoa"`
	Phone           string            `json:"phone"`
	Introduction    string            `json:"introduction"`
	CompletedCourse int64             `json:"completed_course"`
	JoinedEvent     int64             `json:"joined_event"`
	CompletedTask   int64             `json:"completed_task"`
	SocialLink      map[string]string `json:"social_link"`
}
