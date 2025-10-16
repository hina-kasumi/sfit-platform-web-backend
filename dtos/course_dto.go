package dtos

import (
	"encoding/json"
	"sfit-platform-web-backend/entities"
	"time"

	"github.com/google/uuid"
)

// ===================== REQUEST DTOs =====================
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

type UpdateCourseRequest struct {
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
	ModuleTitle string `json:"module_title" binding:"required"`
}

type GetUserProgressInCourseRequest struct {
	CourseID string `json:"course_id" binding:"required"`
	UserID   string `json:"user_id" binding:"required"`
}

type SetFavouriteCourseRequest struct {
	CourseID string `json:"course_id" binding:"required"`
}

type CourseRegisterRequest struct {
	CourseID string                    `json:"course_id" binding:"required"`
	UserIDs  []string                  `json:"user_ids"`
	Msvs     []string                  `json:"msvs"`
	Status   entities.UserCourseStatus `json:"status" binding:"required"`
}

type CourseRateRequest struct {
	Course  string `json:"course" binding:"required"`
	Star    int    `json:"star" binding:"required,min=1,max=5"`
	Comment string `json:"comment"`
}

type GetCourseUserRegisted struct {
	PageListQuery
	Status *entities.UserCourseStatus `form:"status"`
}

type GetUsersInCourse struct {
	PageListQuery
	Status *entities.UserCourseStatus `form:"status"`
}

// ===================== RESPONSE DTOs =====================
type CreateCourseResponse struct {
	ID        string `json:"id"`
	CreatedAt string `json:"createdAt"`
}

type UpdateCourseResponse struct {
	UpdatedAt string `json:"updated_at"`
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
	TotalTime      int                     `json:"total_time"`
	TotalRegitered int                     `json:"total_registered"`
	UpdatedAt      time.Time               `json:"updated_at"`
	Language       string                  `json:"language"`
	CourseContent  []CourseContentResponse `json:"course_content"`
	Rate           []RateResponse          `json:"rate"`
}

type CourseContentResponse struct {
	ID          string           `json:"id"`
	ModuleTitle string           `json:"module_title"`
	Lessons     []LessonResponse `json:"lessons"`
}

type AddModuleToCourseResponse struct {
	ModuleID    string `json:"module_id"`
	CourseID    string `json:"course_id"`
	ModuleTitle string `json:"module_title"`
	CreatedAt   string `json:"created_at"`
}

type GetUserProgressInCourseResponse struct {
	Learned      int `json:"learned"`
	TotalLessons int `json:"total_lesson"`
}

type LessonResponse struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Learned   bool   `json:"learned"`
	StudyTime int    `json:"study_time"`
	Type      string `json:"type"`
	Status    string `json:"status"`
}

type RateResponse struct {
	Name      string  `json:"name"`
	Comment   string  `json:"comment"`
	Star      float64 `json:"star"`
	CreatedAt string  `json:"created_at"`
}

type CourseLessonsResponse []ModuleInfo

type RegisteredUsersResponse struct {
	Users []RegisteredUserInfo `json:"users"`
	// Page     int                  `json:"page"`
	// PageSize int                  `json:"pageSize"`
	// Total    int64                `json:"total"`
	PageListResp
}

type RegisteredUserInfo struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type ModuleInfo struct {
	ID          string       `json:"id"`
	ModuleTitle string       `json:"module_title"`
	TotalTime   int          `json:"total_time"`
	Lessons     []LessonInfo `json:"lessons"`
}

type LessonInfo struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Learned   bool   `json:"learned"`
	StudyTime int    `json:"study_time"`
	Type      string `json:"type"`
	Status    string `json:"status"`
}

type CourseGeneralInformationResponse struct {
	ID             string   `json:"id"`
	Title          string   `json:"title"`
	Description    string   `json:"description"`
	Type           string   `json:"type"`
	NumberLessons  int      `json:"number_lessons"`
	Teachers       []string `json:"teachers"`
	TimeLearn      int      `json:"time_learn"`
	Rate           float64  `json:"rate"`
	Tags           []string `json:"tags"`
	LearnedLessons int      `json:"learned_lessons"`
	Registed       bool     `json:"registed"`
	Status         string   `json:"status"`
}

// type PaginationResponse struct {
// 	CurrentPage  int `json:"current_page"`
// 	TotalPages   int `json:"total_pages"`
// 	TotalCourses int `json:"total_courses"`
// }

// ===================== INTERNAL USE STRUCT =====================
type CourseRaw struct {
	ID             string          `json:"id"`
	Title          string          `json:"title"`
	Description    string          `json:"description"`
	Type           string          `json:"type"`
	Teachers       json.RawMessage `json:"teachers"`
	NumberLessons  int             `json:"number_lessons"`
	TimeLearn      int             `json:"time_learn"`
	Rate           float64         `json:"rate"`
	Tags           json.RawMessage `json:"tags"`
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

// ===================== FILTER STRUCT =====================
type CourseQuery struct {
	Title        string `form:"title"`
	OnlyRegisted bool   `form:"only_registed"`
	CourseType   string `form:"type"`
	Level        string `form:"level"`
	UserID       uuid.UUID
	CourseID     uuid.UUID
	PageListQuery
}
