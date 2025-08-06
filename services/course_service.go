package services

import (
	"sfit-platform-web-backend/dtos"
	"sfit-platform-web-backend/entities"
	"sfit-platform-web-backend/repositories"

	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type CourseService struct {
	repo *repositories.CourseRepository
	db   *gorm.DB
}

func NewCourseService(repo *repositories.CourseRepository, db *gorm.DB) *CourseService {
	return &CourseService{repo: repo, db: db}
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

	// Lấy danh sách khóa học đã đăng ký của người dùng
	courses, err := cs.repo.GetCoursesByUserIDWithPagination(userID, offset, pageSize, &total)
	if err != nil {
		return result, err
	}

	var coursesResp []dtos.CourseInfoResponse

	for _, course := range courses {
		// Lấy số bài học và bài học đã học
		totalLessons, learnedLessons := cs.repo.CountLessonProgress(userUUID, course.ID)

		// Lấy tags
		tags := cs.repo.GetCourseTags(course.ID)

		coursesResp = append(coursesResp, dtos.CourseInfoResponse{
			ID:            course.ID.String(),
			Title:         course.Title,
			Description:   course.Description,
			Type:          course.Level,
			Teachers:      []string(course.Teachers),
			TimeLearn:     0,
			Rate:          5,
			Tags:          tags,
			TotalLesson:   totalLessons,
			LearnedLesson: learnedLessons,
			Registed:      true,
		})
	}

	result = dtos.CourseListResponse{
		Courses:  coursesResp,
		Page:     page,
		PageSize: pageSize,
		Total:    total,
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
	modules, lessonsByModule, learnedLessons, err := cs.repo.GetCourseModulesWithLessons(courseUUID, userUUID)
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
	return cs.repo.CreateOrUpdateCourseRating(userUUID, courseUUID, star, comment)
}

// Xóa khóa học theo ID
func (cs *CourseService) DeleteCourse(courseID string) error {

	courseUUID, err := uuid.Parse(courseID)
	if err != nil {
		return err
	}

	return cs.repo.DeleteCourse(courseUUID)
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
	users, err := cs.repo.GetRegisteredUsersByCourseID(courseUUID, offset, pageSize, &total)
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
	exists, err := cs.repo.CheckCourseExists(courseUUID)
	if err != nil {
		return err
	}
	if !exists {
		return gorm.ErrRecordNotFound
	}

	// Đăng ký người dùng vào khóa học
	return cs.repo.RegisterUserToCourse(userUUID, courseUUID)
}

// Tạo khóa học mới
func (cs *CourseService) CreateCourse(req dtos.CreateCourseRequest) (dtos.CreateCourseResponse, error) {
	var response dtos.CreateCourseResponse

	id := uuid.New()
	now := time.Now()

	// Tạo course entity với pq.StringArray
	course := entities.Course{
		ID:          id,
		Title:       req.Title,
		Description: req.Description,
		Target:      pq.StringArray(req.Target),
		Require:     pq.StringArray(req.Require),
		Teachers:    pq.StringArray(req.Teachers),
		Language:    req.Language,
		Certificate: req.Certificate,
		Level:       req.Level,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	// Sử dụng repository CreateCourse
	err := cs.repo.CreateCourse(&course)
	if err != nil {
		return response, err
	}

	// Tạo response
	response = dtos.CreateCourseResponse{
		ID:          id.String(),
		Title:       req.Title,
		Description: req.Description,
		Target:      req.Target,
		Require:     req.Require,
		Teachers:    req.Teachers,
		Language:    req.Language,
		Certificate: req.Certificate,
		Level:       req.Level,
		CreatedAt:   now.Format("2006-01-02T15:04:05Z07:00"),
	}

	return response, nil
}
