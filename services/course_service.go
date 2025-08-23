package services

import (
	"fmt"
	"sfit-platform-web-backend/dtos"
	"sfit-platform-web-backend/entities"
	"sfit-platform-web-backend/repositories"
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

func (s *CourseService) CreateCourse(req dtos.CreateCourseRequest) (uuid.UUID, time.Time, error) {
	course := entities.Course{
		ID:          uuid.New(),
		Title:       req.Title,
		Description: req.Description,
		Type:        req.Type,
		Target:      req.Target,
		Require:     req.Require,
		Teachers:    req.Teachers,
		Language:    req.Language,
		Certificate: req.Certificate,
		Level:       req.Level,
		CreatedAt:   time.Now(),
	}
	if err := s.courseRepo.CreateNewCourse(&course); err != nil {
		return uuid.Nil, time.Time{}, err
	}
	return course.ID, course.CreatedAt, nil
}

// validatePagination ensures page and pageSize are in valid range
// func validatePagination(page, pageSize int) (int, int) {
// 	if page < minPage {
// 		page = minPage
// 	}
// 	if pageSize < minPage || pageSize > maxPageSize {
// 		pageSize = defaultSize
// 	}
// 	return page, pageSize
// }

func (s *CourseService) GetListCourse(req dtos.CourseQuery) (*dtos.CourseListResponse, error) {
	// Check user exists
	if _, err := s.userRepo.GetUserByID(req.UserID.String()); err != nil {
		return nil, fmt.Errorf("user not found or unauthorized: %w", err)
	}

	// Sanitize pagination
	// page, pageSize = validatePagination(page, pageSize)
	// offset := (req.Page - 1) * req.PageSize

	// Build filter
	// parsedUserID, err := uuid.Parse(userID)
	// if err != nil {
	// 	return nil, fmt.Errorf("invalid user ID: %w", err)
	// }
	// filter := dtos.CourseFilter{
	// 	Title:        strings.TrimSpace(title),
	// 	OnlyRegisted: onlyRegistered,
	// 	CourseType:   strings.TrimSpace(courseType),
	// 	Level:        strings.TrimSpace(level),
	// 	UserID:       parsedUserID,
	// 	Page:         page,
	// 	PageSize:     pageSize,
	// }

	// Query course list
	// courses, total, err := s.courseRepo.GetCourses(filter, pageSize, offset)
	courses, total, err := s.courseRepo.GetCourses(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get courses: %w", err)
	}

	// totalPages := (int(total) + pageSize - 1) / pageSize

	return &dtos.CourseListResponse{
		Courses: courses,
		PageListResp: dtos.PageListResp{
			TotalCount: total,
			Page:       req.Page,
			PageSize:   req.PageSize,
		},
	}, nil
}

// func (s *CourseService) GetMyCourses(userID string, page, pageSize int, title, courseType, level string) (*dtos.CourseListResponse, error) {
// 	// Check user exists
// 	if _, err := s.userRepo.GetUserByID(userID); err != nil {
// 		return nil, fmt.Errorf("user not found: %w", err)
// 	}
// 	// Only registered courses
// 	return s.GetListCourse(userID, page, pageSize, title, true, courseType, level)
// }

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
) (*dtos.PageListResp, error) {
	// Sanitize pagination
	// page, pageSize = validatePagination(page, pageSize)
	// offset := (page - 1) * pageSize

	// Join lesson_attendance and user_course to get users who completed the course
	// Build filter
	parsedCourseID, err := uuid.Parse(courseID)
	if err != nil {
		return nil, fmt.Errorf("invalid course ID: %w", err)
	}
	filter := dtos.CourseQuery{
		CourseID: parsedCourseID,
		PageListQuery: dtos.PageListQuery{
			Page:     page,
			PageSize: pageSize,
		},
	}

	// Query course list
	users, total, err := s.courseRepo.GetListUserCompleteCourses(filter)
	// users, total, err := s.courseRepo.GetListUserCompleteCourses(filter, pageSize, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	totalPages := (int(total) + pageSize - 1) / pageSize

	return &dtos.PageListResp{
		Items:      users,
		Page:       page,
		PageSize:   pageSize,
		TotalCount: int64(totalPages),
	}, nil
}

func (s *CourseService) AddModuleToCourse(courseID, moduleTitle string) (uuid.UUID, time.Time, error) {
	moduleID, create_at, err := s.moduleRepo.AddModuleToCourse(courseID, moduleTitle)
	if err != nil {
		return uuid.Nil, time.Time{}, fmt.Errorf("failed to add module to course: %w", err)
	}
	return moduleID, create_at, nil
}

func (s *CourseService) GetUserProgressInCourse(courseID, userID string) (int, int, error) {
	// Check user exists
	if _, err := s.userRepo.GetUserByID(userID); err != nil {
		return 0, 0, fmt.Errorf("user not found: %w", err)
	}

	// Check course exists
	if _, err := s.courseRepo.GetCourseByID(courseID, userID); err != nil {
		return 0, 0, fmt.Errorf("course not found: %w", err)
	}

	Learned, TotalLessons, err := s.userCourseRepo.GetUserProgressInCourse(courseID, userID)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get user progress: %w", err)
	}

	return Learned, TotalLessons, nil
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

	// Lấy danh sách course từ repo
	courses, err := cs.courseRepo.GetCoursesByUserIDWithPagination(userID, offset, pageSize, &total)
	if err != nil {
		return result, err
	}

	// Bổ sung dữ liệu động: lessons, tags, time_learn, rate
	for i := range courses {
		courseID, _ := uuid.Parse(courses[i].ID)

		// Đếm tổng số bài học và bài đã học
		totalLessons, learnedLessons, err := cs.courseRepo.CountLessonProgress(userUUID, courseID)
		if err != nil {
			return result, err
		}
		courses[i].NumberLessons = totalLessons
		courses[i].LearnedLessons = learnedLessons

		// Lấy tags
		tags := cs.courseRepo.GetCourseTags(courseID)
		courses[i].Tags = tags

		// Lấy time_learn (có sẵn ở bảng courses)
		timeLearn, _ := cs.courseRepo.GetCourseTotalTime(courseID)
		courses[i].TimeLearn = timeLearn

		// Lấy rate trung bình
		rate, _ := cs.courseRepo.GetCourseAverageRate(courseID)
		courses[i].Rate = rate

		courses[i].Registed = true
	}

	result = dtos.CourseListResponse{
		Courses: courses,
		PageListResp: dtos.PageListResp{
			Page:       page,
			PageSize:   pageSize,
			TotalCount: total,
			Items:      nil, // có thể để null
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
				title = lesson.Description
				duration = lesson.Duration
			case "Online":
				title = lesson.Title
				duration = lesson.Duration
			case "Offline":
				title = lesson.OfflineContent.Data.Location
				duration = lesson.Duration
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
		Users: usersResp,
		PageListResp: dtos.PageListResp{
			Page:       page,
			PageSize:   pageSize,
			TotalCount: total,
		},
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

func (s *CourseService) UpdateTotalTime(module_id string, time int) error {
	moduleUUID, err := uuid.Parse(module_id)
	if err != nil {
		return err
	}

	return s.courseRepo.UpdateTotalTime(moduleUUID, time)
}

func (s *CourseService) UpdateTotalLessons(courseID string, lessonCount int) error {
	courseUUID, err := uuid.Parse(courseID)
	if err != nil {
		return err
	}

	_, err = s.courseRepo.UpdateTotalLesson(courseUUID, lessonCount)
	return err
}

func (s *CourseService) GetModuleByID(moduleID string) (*entities.Module, error) {
	moduleUUID, err := uuid.Parse(moduleID)
	if err != nil {
		return nil, err
	}
	return s.courseRepo.GetModuleByID(moduleUUID)
}

func (s *CourseService) GetCourseUserCompletion(userID uuid.UUID) ([]string, error) {
	return s.courseRepo.GetCourseUserCompletion(userID)
}
