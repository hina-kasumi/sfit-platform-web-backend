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


type CourseService struct {
    user_repo                *repositories.UserRepository
	course_repo              *repositories.CourseRepository
	lesson_repo              *repositories.LessonRepository
	tagTemp_repo             *repositories.TagTempRepository
	userCourse_repo          *repositories.UserCourseRepository
	userRate_repo            *repositories.UserRateRepository
	lessonAttendance_repo    *repositories.LessonAttendanceRepository
	module_repo              *repositories.ModuleRepository
}

func NewCourseService(
    user_repo                *repositories.UserRepository,
	course_repo              *repositories.CourseRepository,
	lesson_repo              *repositories.LessonRepository,
	tagTemp_repo             *repositories.TagTempRepository,
	userCourse_repo          *repositories.UserCourseRepository,
	userRate_repo            *repositories.UserRateRepository,
	lessonAttendance_repo    *repositories.LessonAttendanceRepository,
	module_repo              *repositories.ModuleRepository,
) *CourseService {
	return &CourseService{
        user_repo:               user_repo,
		course_repo:             course_repo,
		lesson_repo:             lesson_repo,
		tagTemp_repo:            tagTemp_repo,
		userCourse_repo:         userCourse_repo,
		userRate_repo:           userRate_repo,
		lessonAttendance_repo:   lessonAttendance_repo,
		module_repo:             module_repo,
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
    if err := s.course_repo.CreateNewCourse(&course); err != nil {
        return uuid.Nil, time.Time{}, err
    }
    return course.ID, course.CreatedAt, nil
}

const minPage = 1
const maxPageSize = 100;

func (s *CourseService) GetListCourse(userID string, page, pageSize int, title string, onlyRegisted bool, courseType, level string) (*dtos.CourseListResponse, error) {
	// Validate user exists - REQUIRED
	if _, err := s.user_repo.GetUserByID(userID); err != nil {
		return nil, fmt.Errorf("user not found or unauthorized: %w", err)
	}

	// Validate pagination parameters
	if page < minPage {
		page = 1
	}
	if pageSize < minPage || pageSize > maxPageSize {
		pageSize = 20 // Default page size
	}
	offset := (page - 1) * pageSize

	// Build filter với required userID
    id, err := uuid.Parse(userID)
	filter := dtos.CourseFilter{
		Title:        strings.TrimSpace(title),
		OnlyRegisted: onlyRegisted,
		CourseType:   strings.TrimSpace(courseType),
		Level:        strings.TrimSpace(level),
		UserID:       id,
		Page:         page,
		PageSize:     pageSize,
	}

	// Get courses với user-specific data
	courses, total, err := s.course_repo.GetCourses(filter, pageSize, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get courses: %w", err)
	}

	// Calculate pagination info
	// totalPages := (int(total) + pageSize - 1) / pageSize


	return &dtos.CourseListResponse{
		Courses: courses,
		Pagination: dtos.PaginationResponse{
			CurrentPage:  page,
			TotalPages:   pageSize,
            TotalCourses: total,
		},
	}, nil
}


func (s *CourseService) GetMyCourses(userID string, page, pageSize int, title, courseType, level string) (*dtos.CourseListResponse, error) {
	// Validate user exists
	if _, err := s.user_repo.GetUserByID(userID); err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Validate pagination
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	// Force onlyRegisted = true for "my courses"
	return s.GetListCourse(userID, page, pageSize, title, true, courseType, level)
}