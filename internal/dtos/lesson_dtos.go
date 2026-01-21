package dtos

import (
	"sfit-platform-web-backend/internal/model"
	"time"

	"github.com/google/uuid"
)

type LessonRequest struct {
	Title          string           `json:"title" binding:"required"`
	Description    string           `json:"description"`
	Duration       int              `json:"duration"`
	Type           model.LessonType `json:"type" binding:"required"`
	QuizContent    []model.Quiz     `json:"quiz_content"`
	VideoURL       string           `json:"video_url"`
	Location       string           `json:"location"`
	Date           time.Time        `json:"date"`
	ReadingContent string           `json:"reading_content"`
	Position       float32          `json:"position"`
}

type UpdateStatusLessonAttendanceReq struct {
	Status   model.LessonAttendanceStatus `json:"status"`
	DeviceID string                       `json:"device_id"`
	Duration int                          `json:"duration"`
	Answer   [][]int                      `json:"answer"`
}

type GetUserAttendanceLessonReq struct {
	PageListQuery
	Status model.LessonAttendanceStatus `form:"status"`
}
type GetUserAttendanceLessonRp struct {
	UserID      uuid.UUID  `json:"id" gorm:"column:user_id"`
	Username    string     `json:"username"`
	Email       string     `json:"email"`
	Status      string     `json:"status"`
	QuizPoint   *int       `json:"quiz_point"`
	Duration    *int       `json:"duration"`
	DeviceID    *uuid.UUID `json:"device_id"`
	ModeratorID uuid.UUID  `json:"moderator_id"`
}
