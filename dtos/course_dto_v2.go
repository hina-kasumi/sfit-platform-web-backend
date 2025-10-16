package dtos

import (
	"sfit-platform-web-backend/entities"
	"time"

	"github.com/google/uuid"
)

type GetListCoursesForm struct {
	*PageListQuery
	Title      string `form:"title"`
	UserStatus string `form:"status"`
	Type       string `form:"type"`
	Level      string `form:"level"`
}

type GetListCoursesResponse struct {
	ID             uuid.UUID `json:"id"`
	Title          string    `json:"title"`
	Description    string    `json:"description"`
	Type           string    `json:"type"`
	Teachers       []string  `json:"teachers"`
	TimeLearn      int       `json:"time_learn"`
	Rate           float32   `json:"rate"`
	Tags           []string  `json:"tags"`
	TotalLessons   int       `json:"total_lessons"`
	LearnedLessons int       `json:"learned_lessons"`
	UserStatus     string    `json:"status"`
}

type GetCourseDetailResponse struct {
	ID           uuid.UUID           `json:"id"`
	Title        string              `json:"title"`
	Liked        bool                `json:"liked"`
	Description  string              `json:"description"`
	Type         string              `json:"type"`
	Level        string              `json:"level"`
	Teachers     []string            `json:"teachers"`
	Star         float32             `json:"star"`
	TotalLessons int                 `json:"total_lessons"`
	Tags         []string            `json:"tags"`
	Target       []string            `json:"targets"`
	Require      []string            `json:"requires"`
	Language     string              `json:"language"`
	TotalTime    int                 `json:"total_time"`
	TotalLearned int                 `json:"total_registered"`
	Rate         []entities.UserRate `json:"rate"`
	UpdatedAt    time.Time           `json:"updated_at"`
	CreatedAt    time.Time           `json:"created_at"`
	Certificate  bool                `json:"certificate"`
	Status       string              `json:"status"`
	Modules      []ModuleResponse    `json:"modules"`
}

type ModuleResponse struct {
	ID           uuid.UUID        `json:"id"`
	Title        string           `json:"module_title"`
	CreatedAt    time.Time        `json:"created_at"`
	UpdatedAt    time.Time        `json:"updated_at"`
	TotalTime    int              `json:"total_time"`
	TotalLessons int              `json:"total_lessons"`
	Lessons      []LessonResponse `json:"lessons"`
}
