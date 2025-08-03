package dtos

import (
	"encoding/json"

	"github.com/google/uuid"
)

type CreateCourseRequest struct {
	Title       string   `json:"title" binding:"required"`
	Description string   `json:"description" binding:"required"`
	Type        string   `json:"type" binding:"required"`
	Target      []string `json:"targets"`
	Require     []string `json:"requires"`
	Teachers    []string `json:"teachers"`
	Language    string   `json:"language" binding:"required"`
	Certificate bool     `json:"certificate"`
	Level       string   `json:"level" binding:"required"`
	Tags        []string `json:"tags"`
}

type CreateCourseResponse struct {
	ID        string `json:"id"`
	CreatedAt string `json:"createdAt"`
}

type CourseInformationResponse struct {
	ID             string   `json:"id"`
	Title          string   `json:"title"`
	Description    string   `json:"description"`
	Type           string   `json:"type"`
	Teachers       []string `json:"teachers"`
	NumberLessons  int      `json:"numberLessons"`
	TimeLearn      int      `json:"timeLearning"`
	Rate           float64  `json:"rate"`
	Tags           []string `json:"tags"`
	LearnedLessons int      `json:"learnedLessons"`
	Registed       bool     `json:"registed"`
}

type CourseRaw struct {
	ID             string          `json:"id"`
	Title          string          `json:"title"`
	Description    string          `json:"description"`
	Type           string          `json:"type"`
	Teachers       json.RawMessage `json:"teachers"` // JSON array → parse sau
	NumberLessons  int             `json:"number_lessons"`
	TimeLearn      int             `json:"time_learn"`
	Rate           float64         `json:"rate"`
	Tags           json.RawMessage `json:"tags"` // JSON array of UUIDs → parse sau
	LearnedLessons int             `json:"learned_lessons"`
	Registed       bool            `json:"registed"`
}


type CourseFilter struct {
	Title        string
	OnlyRegisted bool
	CourseType   string
	Level        string
	UserID       uuid.UUID
	Page         int
	PageSize     int
}

type CourseListResponse struct {
	Courses    []CourseInformationResponse `json:"courses"`
	Pagination PaginationResponse          `json:"pagination"`
}

type PaginationResponse struct {
	CurrentPage  int `json:"currentPage"`
	TotalPages   int `json:"totalPages"`
	TotalCourses int `json:"totalCourses"`
}