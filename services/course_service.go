package services

import (
	"fmt"
	"sfit-platform-web-backend/dtos"
	"sfit-platform-web-backend/entities"
	"sfit-platform-web-backend/repositories"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
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

// Lấy danh sách khóa học đã đăng ký của người dùng với phân trang
func (cs *CourseService) GetRegisteredCourses(userID string, page, pageSize int) (dtos.CourseListResponse, error) {
    offset := (page - 1) * pageSize
    var total int64
    var result dtos.CourseListResponse

    userUUID, err := uuid.Parse(userID)
    if err != nil {
        return result, err
    }

    // Lấy danh sách khóa học đã đăng ký (theo DTO)
    courses, err := cs.courseRepo.GetCoursesByUserIDWithPagination(userID, offset, pageSize, &total)
    if err != nil {
        return result, err
    }

    // Bổ sung thông tin lessons, tags
    for i := range courses {
        courseID, _ := uuid.Parse(courses[i].ID)
        totalLessons, learnedLessons := cs.courseRepo.CountLessonProgress(userUUID, courseID)
        tags := cs.courseRepo.GetCourseTags(courseID)

        courses[i].NumberLessons = totalLessons
        courses[i].LearnedLessons = learnedLessons
        courses[i].Tags = tags
    }

    totalPages := int((total + int64(pageSize) - 1) / int64(pageSize))
    result = dtos.CourseListResponse{
        Courses: courses,
        Pagination: dtos.PaginationResponse{
            CurrentPage:  page,
            TotalPages:   totalPages,
            TotalCourses: int(total),
        },
    }
    return result, nil
}

// Lấy danh sách bài học trong khóa học
func (cs *CourseService) GetCourseLessons(courseID string, userID string) (dtos.CourseLessonsResponse, error) {
	// lấy courseID và userID từ chuỗi
	courseUUID, err := uuid.Parse(courseID)
	if err != nil {
		return nil, err
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}

	// Lấy danh sách các module và bài học trong khóa học
	modules, lessonsByModule, learnedLessons, err := cs.courseRepo.GetCourseModulesWithLessons(courseUUID, userUUID)
	if err != nil {
		return nil, err
	}

	var response dtos.CourseLessonsResponse

	for _, module := range modules {
		var lessons []dtos.LessonInfo
		totalTime := 0

		for _, lesson := range lessonsByModule[module.ID] {
			// Lấy tiêu đề và thời gian của bài học dựa trên loại bài học
			var title string
			var duration int

			switch lesson.Type {
			case "Quiz":
				title = lesson.QuizContent.Data.Description
				duration = lesson.QuizContent.Data.Duration
			case "Online":
				title = lesson.OnlineContent.Data.Title
				duration = lesson.OnlineContent.Data.Duration
			case "Offline":
				title = lesson.OfflineContent.Data.Location
				duration = lesson.OfflineContent.Data.Duration
			}

			totalTime += duration

			lessons = append(lessons, dtos.LessonInfo{
				ID:        lesson.ID.String(),
				Title:     title,
				Learned:   learnedLessons[lesson.ID],
				StudyTime: duration,
			})
		}

		moduleInfo := dtos.ModuleInfo{
			ID:          module.ID.String(),
			ModuleTitle: module.Title,
			TotalTime:   totalTime,
			Lessons:     lessons,
		}

		response = append(response, moduleInfo)
	}

	return response, nil
}

// Đánh giá khóa học của người dùng
func (cs *CourseService) RateCourse(userID string, courseID string, star int, comment string) error {

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return err
	}

	courseUUID, err := uuid.Parse(courseID)
	if err != nil {
		return err
	}

	// Đánh giá khóa học của người dùng
	return cs.courseRepo.CreateOrUpdateCourseRating(userUUID, courseUUID, star, comment)
}

// Xóa khóa học theo ID
func (cs *CourseService) DeleteCourse(courseID string) error {

	courseUUID, err := uuid.Parse(courseID)
	if err != nil {
		return err
	}

	return cs.courseRepo.DeleteCourse(courseUUID)
}

// Lấy danh sách người dùng đã đăng ký khóa học
func (cs *CourseService) GetRegisteredUsers(courseID string, page, pageSize int) (dtos.RegisteredUsersResponse, error) {
	offset := (page - 1) * pageSize
	var total int64
	var result dtos.RegisteredUsersResponse

	courseUUID, err := uuid.Parse(courseID)
	if err != nil {
		return result, err
	}

	// kiểm tra xem khóa học có tồn tại không
	users, err := cs.courseRepo.GetRegisteredUsersByCourseID(courseUUID, offset, pageSize, &total)
	if err != nil {
		return result, err
	}

	var usersResp []dtos.RegisteredUserInfo

	for _, user := range users {
		usersResp = append(usersResp, dtos.RegisteredUserInfo{
			ID:       user.ID.String(),
			Username: user.Username,
			Email:    user.Email,
		})
	}

	result = dtos.RegisteredUsersResponse{
		Users:    usersResp,
		Page:     page,
		PageSize: pageSize,
		Total:    total,
	}
	return result, nil
}

// Đăng ký người dùng vào khóa học
func (cs *CourseService) RegisterUserToCourse(userID string, courseID string) error {
	// Lấy userID và courseID từ chuỗi
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return err
	}

	courseUUID, err := uuid.Parse(courseID)
	if err != nil {
		return err
	}

	// Kiểm tra xem khóa học có tồn tại hay không
	exists, err := cs.courseRepo.CheckCourseExists(courseUUID)
	if err != nil {
		return err
	}
	if !exists {
		return gorm.ErrRecordNotFound
	}

	// Đăng ký người dùng vào khóa học
	return cs.courseRepo.RegisterUserToCourse(userUUID, courseUUID)
}
