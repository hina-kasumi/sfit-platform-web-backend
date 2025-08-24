package repositories

import (
	"sfit-platform-web-backend/dtos"
	"sfit-platform-web-backend/entities"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LessonRepository struct {
	db *gorm.DB
}

func NewLessonRepository(db *gorm.DB) *LessonRepository {
	return &LessonRepository{
		db: db,
	}
}

func (repo *LessonRepository) CreateLesson(lesson *entities.Lesson) (*entities.Lesson, error) {
	if err := repo.db.Create(lesson).Error; err != nil {
		return nil, err
	}
	return lesson, nil
}

func (repo *LessonRepository) GetLessonByID(id string) (*entities.Lesson, error) {
	var lesson entities.Lesson
	if err := repo.db.First(&lesson, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &lesson, nil
}

func (repo *LessonRepository) DeleteLessonByID(id string) error {
	return repo.db.Delete(&entities.Lesson{}, "id = ?", id).Error
}

func (repo *LessonRepository) UpdateLesson(lesson *entities.Lesson) error {
	return repo.db.Save(lesson).Error
}

func (repo *LessonRepository) UpdateStatusLessonAttendance(
	userID, lessonID, courseID uuid.UUID,
	status entities.LessonAttendanceStatus,
	deviceID string, quizPoint int,
	currentUserID string,
	duration int,
) error {
	var dvID *uuid.UUID
	if deviceID == "" {
		dvID = nil
	} else {
		parsedID := uuid.MustParse(deviceID)
		dvID = &parsedID
	}
	lessonAttendance := entities.LessonAttendance{
		UserID:      userID,
		LessonID:    lessonID,
		CourseID:    courseID,
		Status:      status,
		DeviceID:    dvID,
		QuizPoint:   &quizPoint,
		ModeratorID: uuid.MustParse(currentUserID),
		Duration:    &duration,
	}
	// Try to update, if no rows affected, create new
	result := repo.db.
		Model(&entities.LessonAttendance{}).
		Where("user_id = ? AND lesson_id = ?", userID, lessonID).
		Updates(&lessonAttendance)
	if result.RowsAffected == 0 {
		return repo.db.Create(&lessonAttendance).Error
	}
	return result.Error
}

func (repo *LessonRepository) GetUsersByLessonID(lessonID uuid.UUID, req dtos.GetUserAttendanceLessonReq) ([]dtos.GetUserAttendanceLessonRp, int64, error) {
	var resp []dtos.GetUserAttendanceLessonRp
	query := repo.db.Model(&entities.LessonAttendance{}).
		Select("users.id as user_id, users.username, users.email, status, quiz_point, duration, device_id, moderator_id").
		Where("lesson_id = ?", lessonID).
		Joins("JOIN users ON users.id = lesson_attendances.user_id")
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}
	var total int64
	query = query.Count(&total)
	query = query.Offset((req.Page - 1) * req.PageSize).Limit(req.PageSize)

	if err := query.Scan(&resp).Error; err != nil {
		return nil, 0, err
	}

	if len(resp) == 0 {
		return []dtos.GetUserAttendanceLessonRp{}, 0, nil
	}
	return resp, total, nil
}
