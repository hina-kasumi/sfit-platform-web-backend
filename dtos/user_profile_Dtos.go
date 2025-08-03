package dtos

import (
	"github.com/google/uuid"
	"time"
)

type UpdateUserProfileRequest struct {
	FullName     string            `json:"full_name"`
	ClassName    string            `json:"class_name"`
	Khoa         string            `json:"khoa"`
	Phone        string            `json:"phone"`
	Introduction string            `json:"introduction"`
	SocialLink   map[string]string `json:"social_link"`
}

type GetUserProfileResponse struct {
	UserID          uuid.UUID         `json:"user_id"`
	FullName        string            `json:"full_name"`
	ClassName       string            `json:"class_name"`
	Khoa            string            `json:"khoa"`
	Phone           string            `json:"phone"`
	Introduction    string            `json:"introduction"`
	CompletedCourse int64             `json:"completed_course"`
	JoinedEvent     int64             `json:"joined_event"`
	CompletedTask   int64             `json:"completed_task"`
	SocialLink      map[string]string `json:"social_link"`
	CreatedAt       time.Time         `json:"created_at"`
	UpdatedAt       time.Time         `json:"updated_at"`
}
