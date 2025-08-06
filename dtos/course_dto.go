package dtos

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

//
// REQUEST DTOs
//

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

type SetFavouriteCourseRequest struct {
	CourseID string `json:"course_id" binding:"required"`
}

type UpdateCourseRequest struct {
	ID          string   `json:"id" binding:"required"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Type        string   `json:"type"`
	Target      []string `json:"targets"`
	Require     []string `json:"requires"`
	Teachers    []string `json:"teachers"`
	Language    string   `json:"language"`
	Certificate bool     `json:"certificate"`
	Level       string   `json:"level"`
	Tags        []string `json:"tags"`
}

type AddModuleToCourseRequest struct {
	CourseID    string `json:"course_id" binding:"required"`
	ModuleTitle string `json:"module_title" binding:"required"`
}

//
// RESPONSE DTOs
//

type CreateCourseResponse struct {
	ID        string `json:"id"`
	CreatedAt string `json:"createdAt"`
}

type CourseListResponse struct {
	Courses    []CourseGeneralInformationResponse `json:"courses"`
	Pagination PaginationResponse                 `json:"pagination"`
}

type PaginationResponse struct {
	CurrentPage  int `json:"currentPage"`
	TotalPages   int `json:"totalPages"`
	TotalCourses int `json:"totalCourses"`
}

type CourseGeneralInformationResponse struct {
	ID             string   `json:"id"`
	Title          string   `json:"title"`
	Description    string   `json:"description"`
	Type           string   `json:"type"`
	Teachers       []string `json:"teachers"`
	NumberLessons  int      `json:"numberLessons"`
	TimeLearn      int      `json:"timeLearned"`
	Rate           float64  `json:"rate"`
	Tags           []string `json:"tags"`
	LearnedLessons int      `json:"learnedLessons"`
	Registed       bool     `json:"registed"`
}

type CourseDetailResponse struct {
	Title          string                  `json:"title"`
	Description    string                  `json:"description"`
	Like           bool                    `json:"like"`
	Type           string                  `json:"type"`
	Level          string                  `json:"level"`
	Teachers       []string                `json:"teachers"`
	Star           float64                 `json:"star"`
	TotalLessons   int                     `json:"total_lessons"`
	Tags           []string                `json:"tags"`
	Target         []string                `json:"target"`
	Require        []string                `json:"require"`
	TotalTime      int                     `json:"total_time"` // total time in seconds
	TotalRegitered int                     `json:"total_registered"`
	UpdatedAt      time.Time               `json:"updated_at"`
	Language       string                  `json:"language"`
	CourseContent  []CourseContentResponse `json:"course_content"`
	Rate           []RateResponse          `json:"rate"`
}

type CourseContentResponse struct {
	ID          string `json:"id"`
	ModuleTitle string `json:"module_title"`
	// TotalTime   int              `json:"total_time"`
	Lessons []LessonResponse `json:"lessons"`
}

type LessonResponse struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Learned   bool   `json:"learned"`
	StudyTime int    `json:"study_time"`
}

type RateResponse struct {
	Name      string  `json:"name"`
	Comment   string  `json:"comment"`
	Star      float64 `json:"star"`
	CreatedAt string  `json:"created_at"`
}

type UpdateCourseResponse struct {
	// ID string  `json:"id"`
	UpdatedAt string `json:"updated_at"`
}

type AddModuleToCourseResponse struct {
	ModuleID    string `json:"module_id"`
	CourseID    string `json:"course_id"`
	ModuleTitle string `json:"module_title"`
	CreatedAt   string `json:"created_at"`
}

//
// INTERNAL USE STRUCT
//

type CourseRaw struct {
	ID             string          `json:"id"`
	Title          string          `json:"title"`
	Description    string          `json:"description"`
	Type           string          `json:"type"`
	Teachers       json.RawMessage `json:"teachers"` // JSON array of strings
	NumberLessons  int             `json:"number_lessons"`
	TimeLearn      int             `json:"time_learn"`
	Rate           float64         `json:"rate"`
	Tags           json.RawMessage `json:"tags"` // JSON array of UUID strings
	LearnedLessons int             `json:"learned_lessons"`
	Registed       bool            `json:"registed"`
}

type CourseDetailRaw struct {
	Title           string          `json:"title"`
	Description     string          `json:"description"`
	Like            bool            `json:"like"`
	Type            string          `json:"type"`
	Level           string          `json:"level"`
	Teachers        json.RawMessage `json:"teachers"`
	Star            float64         `json:"star"`
	TotalLessons    int             `json:"total_lessons"`
	Tags            json.RawMessage `json:"tags"`
	Target          json.RawMessage `json:"target"`
	Require         json.RawMessage `json:"require"`
	TotalTime       int             `json:"total_time"`
	TotalRegistered int             `json:"total_registered"`
	UpdatedAt       time.Time       `json:"updated_at"`
	Language        string          `json:"language"`
}

//
// FILTER STRUCT
//

type CourseFilter struct {
	Title        string
	OnlyRegisted bool
	CourseType   string
	Level        string
	UserID       uuid.UUID
	CourseID     uuid.UUID
	Page         int
	PageSize     int
}
