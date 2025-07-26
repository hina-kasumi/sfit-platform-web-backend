package entities

import (
	"sfit-platform-web-backend/utits/converter"
	"time"

	"github.com/google/uuid"
)

type LessonType string

const (
	QuizLesson    LessonType = "Quiz"
	OnlineLesson  LessonType = "Online"
	OfflineLesson LessonType = "Offline"
)

type Lesson struct {
	ID             uuid.UUID                             `gorm:"type:uuid;primaryKey"`
	Type           LessonType                            `gorm:"type:varchar;column:lesson_type;check:lesson_type in ('Quiz', 'Online', 'Offline')"`
	ModuleID       uuid.UUID                             `gorm:"type:uuid;column:module_id;not null"`
	QuizContent    converter.JSONB[QuizContentStruct]    `gorm:"type:jsonb;column:quiz_content"`
	OnlineContent  converter.JSONB[OnlineContentStruct]  `gorm:"type:jsonb;column:online_content"`
	OfflineContent converter.JSONB[OfflineContentStruct] `gorm:"type:jsonb;column:offline_content"`
	CreatedAt      time.Time                             `gorm:"column:create_at"`
	UpdatedAt      time.Time                             `gorm:"column:update_at"`
}

type QuizContentStruct struct {
	Title       []string `json:"questions"`
	Description string   `json:"description"`
	Quiz        Quiz     `json:"quiz"`
	Duration    int      `json:"duration"`
}

type Quiz struct {
	Questions []string `json:"questions"`
	Answers   []string `json:"answers"`
}

type OnlineContentStruct struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	VideoURL    string `json:"video_url"`
	Duration    int    `json:"duration"`
}

type OfflineContentStruct struct {
	Location string `json:"location"`
	Date     string `json:"date"`
	Duration int    `json:"duration"`
}
