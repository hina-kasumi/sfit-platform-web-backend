package entities

import (
	"sfit-platform-web-backend/utils/converter"
	"time"

	"github.com/google/uuid"
)

type LessonType string

const (
	QuizLesson    LessonType = "Quiz"
	OnlineLesson  LessonType = "Online"
	OfflineLesson LessonType = "Offline"
	ReadingLesson LessonType = "Reading"
)

type Lesson struct {
	ID             uuid.UUID                             `gorm:"type:uuid;primaryKey"`
	Type           LessonType                            `gorm:"type:varchar;column:lesson_type"`
	Title          string                                `gorm:"type:varchar;column:title;not null"`
	ModuleID       uuid.UUID                             `gorm:"type:uuid;column:module_id;not null"`
	CourseID       uuid.UUID                             `gorm:"type:uuid;column:course_id;not null"`
	Description    string                                `gorm:"type:text;column:description;not null"`
	Duration       int                                   `gorm:"type:int;column:duration;not null"`
	QuizContent    converter.JSONB[[]Quiz]               `gorm:"type:jsonb;column:quiz_content;default:null"`
	OnlineContent  converter.JSONB[OnlineContentStruct]  `gorm:"type:jsonb;column:online_content;default:null"`
	OfflineContent converter.JSONB[OfflineContentStruct] `gorm:"type:jsonb;column:offline_content;default:null"`
	ReadingContent converter.JSONB[ReadingContentStruct] `gorm:"type:jsonb;column:reading_content;default:null"`
	CreatedAt      time.Time                             `gorm:"column:create_at"`
	UpdatedAt      time.Time                             `gorm:"column:update_at"`
	Position       float32                               `gorm:"type:real;column:position;default:0"`
}

type QuizContentStruct struct {
	Quiz []Quiz `json:"quiz"`
}

type Quiz struct {
	Question       string   `json:"questions"`
	Answers        []string `json:"answers"`
	CorrectAnswers []int    `json:"correct_answers"`
}

type OnlineContentStruct struct {
	VideoURL string `json:"video_url"`
}

type OfflineContentStruct struct {
	Location string    `json:"location"`
	Date     time.Time `json:"date"`
}

type ReadingContentStruct struct {
	Content string `json:"content"`
}

func newLesson(courseID, moduleID uuid.UUID, title, description string, position float32, duration int, lessonType LessonType) *Lesson {
	return &Lesson{
		ID:          uuid.New(),
		ModuleID:    moduleID,
		Title:       title,
		Description: description,
		CourseID:    courseID,
		Duration:    duration,
		Type:        lessonType,
		Position:    position,
	}
}

func NewQuizLesson(courseID, moduleID uuid.UUID, title, description string, position float32, duration int, quiz []Quiz) *Lesson {
	lesson := newLesson(courseID, moduleID, title, description, position, duration, QuizLesson)
	lesson.QuizContent = converter.JSONB[[]Quiz]{Data: quiz}
	return lesson
}

func NewOnlineLesson(courseID, moduleID uuid.UUID, title, description string, position float32, duration int, videoURL string) *Lesson {
	lesson := newLesson(courseID, moduleID, title, description, position, duration, OnlineLesson)
	lesson.OnlineContent = converter.JSONB[OnlineContentStruct]{Data: OnlineContentStruct{
		VideoURL: videoURL,
	}}
	return lesson
}

func NewOfflineLesson(courseID, moduleID uuid.UUID, title, description string, position float32, duration int, location string, date time.Time) *Lesson {
	lesson := newLesson(courseID, moduleID, title, description, position, duration, OfflineLesson)
	lesson.OfflineContent = converter.JSONB[OfflineContentStruct]{Data: OfflineContentStruct{
		Location: location,
		Date:     date,
	}}
	return lesson
}

func NewReadingLesson(courseID, moduleID uuid.UUID, title, description string, position float32, duration int, content string) *Lesson {
	lesson := newLesson(courseID, moduleID, title, description, position, duration, ReadingLesson)
	lesson.ReadingContent = converter.JSONB[ReadingContentStruct]{Data: ReadingContentStruct{
		Content: content,
	}}
	return lesson
}
