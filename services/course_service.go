package services

import (
	"fmt"
	"sfit-platform-web-backend/dtos"
	"sfit-platform-web-backend/entities"
	"sfit-platform-web-backend/repositories"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	minPage     = 1
	maxPageSize = 100
	defaultSize = 20
)

type CourseService struct {
	userRepo             *repositories.UserRepository
	courseRepo           *repositories.CourseRepository
	favorCourseRepo      *repositories.FavoriteCourseRepository
	lessonRepo           *repositories.LessonRepository
	tagTempRepo          *repositories.TagTempRepository
	userCourseRepo       *repositories.UserCourseRepository
	userRateRepo         *repositories.UserRateRepository
	lessonAttendanceRepo *repositories.LessonAttendanceRepository
	moduleRepo           *repositories.ModuleRepository
}

func NewCourseService(
	userRepo *repositories.UserRepository,
	courseRepo *repositories.CourseRepository,
	favorCourseRepo *repositories.FavoriteCourseRepository,
	lessonRepo *repositories.LessonRepository,
	tagTempRepo *repositories.TagTempRepository,
	userCourseRepo *repositories.UserCourseRepository,
	userRateRepo *repositories.UserRateRepository,
	lessonAttendanceRepo *repositories.LessonAttendanceRepository,
	moduleRepo *repositories.ModuleRepository,
) *CourseService {
	return &CourseService{
		userRepo:             userRepo,
		courseRepo:           courseRepo,
		favorCourseRepo:      favorCourseRepo,
		lessonRepo:           lessonRepo,
		tagTempRepo:          tagTempRepo,
		userCourseRepo:       userCourseRepo,
		userRateRepo:         userRateRepo,
		lessonAttendanceRepo: lessonAttendanceRepo,
		moduleRepo:           moduleRepo,
	}
}

func (s *CourseService) CreateCourse(
	title, description, courseType string,
	targets, requires, teachers []string,
	language string,
	certificate bool,
	level string,
) (uuid.UUID, time.Time, error) {
	course := entities.Course{
		ID:          uuid.New(),
		Title:       title,
		Description: description,
		Type:        courseType,
		Target:      targets,
		Require:     requires,
		Teachers:    teachers,
		Language:    language,
		Certificate: certificate,
		Level:       level,
		CreatedAt:   time.Now(),
	}
	if err := s.courseRepo.CreateNewCourse(&course); err != nil {
		return uuid.Nil, time.Time{}, err
	}
	return course.ID, course.CreatedAt, nil
}

// validatePagination ensures page and pageSize are in valid range
func validatePagination(page, pageSize int) (int, int) {
	if page < minPage {
		page = minPage
	}
	if pageSize < minPage || pageSize > maxPageSize {
		pageSize = defaultSize
	}
	return page, pageSize
}

func (s *CourseService) GetListCourse(
	userID string,
	page, pageSize int,
	title string,
	onlyRegistered bool,
	courseType, level string,
) (*dtos.CourseListResponse, error) {
	// Check user exists
	if _, err := s.userRepo.GetUserByID(userID); err != nil {
		return nil, fmt.Errorf("user not found or unauthorized: %w", err)
	}

	// Sanitize pagination
	page, pageSize = validatePagination(page, pageSize)
	offset := (page - 1) * pageSize

	// Build filter
	parsedUserID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}
	filter := dtos.CourseFilter{
		Title:        strings.TrimSpace(title),
		OnlyRegisted: onlyRegistered,
		CourseType:   strings.TrimSpace(courseType),
		Level:        strings.TrimSpace(level),
		UserID:       parsedUserID,
		Page:         page,
		PageSize:     pageSize,
	}

	// Query course list
	courses, total, err := s.courseRepo.GetCourses(filter, pageSize, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get courses: %w", err)
	}

	totalPages := (int(total) + pageSize - 1) / pageSize

	return &dtos.CourseListResponse{
		Courses: courses,
		Pagination: dtos.PaginationResponse{
			CurrentPage:  page,
			TotalPages:   totalPages,
			TotalCourses: total,
		},
	}, nil
}

func (s *CourseService) GetMyCourses(userID string, page, pageSize int, title, courseType, level string) (*dtos.CourseListResponse, error) {
	// Check user exists
	if _, err := s.userRepo.GetUserByID(userID); err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	// Only registered courses
	return s.GetListCourse(userID, page, pageSize, title, true, courseType, level)
}

func (s *CourseService) GetCourseDetailByID(courseID string, userID string) (*dtos.CourseDetailResponse, error) {
	// Check user exists
	if _, err := s.userRepo.GetUserByID(userID); err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	course, err := s.courseRepo.GetCourseByID(courseID, userID)
	if err != nil {
		return nil, fmt.Errorf("course not found: %w", err)
	}

	return course, nil
}

func (s *CourseService) MarkCourseAsFavourite(userID, courseID string) error {
	// Check user exists
	if _, err := s.userRepo.GetUserByID(userID); err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	// Check course exists
	if _, err := s.courseRepo.GetCourseByID(courseID, userID); err != nil {
		return fmt.Errorf("course not found: %w", err)
	}

	if err := s.favorCourseRepo.MarkCourseAsFavourite(userID, courseID); err != nil {
		return fmt.Errorf("failed to mark course as favourite: %w", err)
	}

	return nil
}

func (s *CourseService) UnmarkCourseAsFavourite(userID, courseID string) error {
	// Check user exists
	if _, err := s.userRepo.GetUserByID(userID); err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	// Check course exists
	if _, err := s.courseRepo.GetCourseByID(courseID, userID); err != nil {
		return fmt.Errorf("course not found: %w", err)
	}

	if err := s.favorCourseRepo.UnmarkCourseAsFavourite(userID, courseID); err != nil {
		return fmt.Errorf("failed to unmark course as favourite: %w", err)
	}

	return nil
}

func (s *CourseService) UpdateCourse(
	courseID, title, description, courseType string,
	targets, requires, teachers []string,
	language string,
	certificate bool,
	level string,
) (time.Time, error) {
	// Validate course exists
	// userID := middlewares.GetPrincipal()
	// if userID == "" {
	// 	return fmt.Errorf("unauthorized")
	// }

	// if _, err := s.courseRepo.GetCourseByID(courseID, userID); err != nil {
	// 	return fmt.Errorf("course not found: %w", err)
	// }

	// Prepare updated course
	updatedCourse := entities.Course{
		ID:          uuid.MustParse(courseID),
		Title:       title,
		Description: description,
		Type:        courseType,
		Target:      targets,
		Require:     requires,
		Teachers:    teachers,
		Language:    language,
		Certificate: certificate,
		Level:       level,
		UpdatedAt:   time.Now(),
	}

	// Update in repository
	if err := s.courseRepo.UpdateCourse(&updatedCourse); err != nil {
		return time.Time{}, fmt.Errorf("failed to update course: %w", err)
	}

	return updatedCourse.UpdatedAt, nil
}

func (s *CourseService) GetListUserCompleteCourse(
	courseID string,
	page, pageSize int,
) (*dtos.UserListResponse, error) {
	// Sanitize pagination
	page, pageSize = validatePagination(page, pageSize)
	offset := (page - 1) * pageSize

	// Join lesson_attendance and user_course to get users who completed the course
	// Build filter
	parsedCourseID, err := uuid.Parse(courseID)
	if err != nil {
		return nil, fmt.Errorf("invalid course ID: %w", err)
	}
	filter := dtos.CourseFilter{
		CourseID: parsedCourseID,
		Page:     page,
		PageSize: pageSize,
	}

	// Query course list
	users, total, err := s.courseRepo.GetListUserCompleteCourses(filter, pageSize, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	totalPages := (int(total) + pageSize - 1) / pageSize

	return &dtos.UserListResponse{
		Users:    users,
		Page:     page,
		PageSize: pageSize,
		Total:    int64(totalPages),
	}, nil
}
