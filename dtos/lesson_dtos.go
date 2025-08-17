package dtos

import (
	"sfit-platform-web-backend/entities"
	"time"
)

type LessonRequest struct {
	Title          string          `json:"title" binding:"required"`
	Description    string          `json:"description"`
	Duration       int             `json:"duration"`
	Type           string          `json:"type" binding:"required,oneof=Quiz Online Offline Reading"`
	QuizContent    []entities.Quiz `json:"quiz_content"`
	VideoURL       string          `json:"video_url"`
	Location       string          `json:"location"`
	Date           time.Time       `json:"date"`
	ReadingContent string          `json:"reading_content"`
}
