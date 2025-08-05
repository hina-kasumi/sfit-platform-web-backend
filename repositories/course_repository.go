package repositories

import (
	"sfit-platform-web-backend/entities"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CourseRepository struct {
	db *gorm.DB
}

func NewCourseRepository(db *gorm.DB) *CourseRepository {
	return &CourseRepository{db: db}
}

// Lấy danh sách khóa học đã đăng ký của người dùng với phân trang
func (cr *CourseRepository) GetCoursesByUserIDWithPagination(userID string, offset, limit int, total *int64) ([]entities.Course, error) {
	var courses []entities.Course
	err := cr.db.
		Table("courses").
		Joins("JOIN user_courses ON user_courses.course_id = courses.id").
		Where("user_courses.user_id = ?", userID).
		Count(total).
		Offset(offset).
		Limit(limit).
		Find(&courses).Error

	return courses, err
}

// Lấy số bài học và bài học đã học của người dùng trong khóa học
func (cr *CourseRepository) CountLessonProgress(userID uuid.UUID, courseID uuid.UUID) (int, int) {
	var total, learned int64

	cr.db.
		Table("lessons").
		Joins("JOIN modules ON lessons.module_id = modules.id").
		Where("modules.course_id = ?", courseID).
		Count(&total)

	cr.db.
		Table("lesson_attendances").
		Joins("JOIN lessons ON lesson_attendances.lesson_id = lessons.id").
		Joins("JOIN modules ON lessons.module_id = modules.id").
		Where("lesson_attendances.user_id = ? AND modules.course_id = ?", userID, courseID).
		Count(&learned)

	return int(total), int(learned)
}

// Lấy tags của khóa học
func (cr *CourseRepository) GetCourseTags(courseID uuid.UUID) []string {
	var tags []string
	cr.db.
		Table("tags").
		Joins("JOIN tag_temps ON tags.id = tag_temps.tag_id").
		Where("tag_temps.course_id = ?", courseID).
		Pluck("tags.id", &tags)

	return tags
}

// Lấy danh sách các module và bài học trong khóa học
func (cr *CourseRepository) GetCourseModulesWithLessons(courseID uuid.UUID, userID uuid.UUID) ([]entities.Module, map[uuid.UUID][]entities.Lesson, map[uuid.UUID]bool, error) {
	var modules []entities.Module

	// Lấy tất cả các module của khóa học
	err := cr.db.Where("course_id = ?", courseID).Find(&modules).Error
	if err != nil {
		return nil, nil, nil, err
	}

	// Lấy tất cả các bài học cho mỗi module
	lessonsByModule := make(map[uuid.UUID][]entities.Lesson)
	learnedLessons := make(map[uuid.UUID]bool)

	for _, module := range modules {
		var lessons []entities.Lesson
		err := cr.db.Where("module_id = ?", module.ID).Find(&lessons).Error
		if err != nil {
			return nil, nil, nil, err
		}
		lessonsByModule[module.ID] = lessons

		// Kiểm tra bài học nào đã được người dùng học
		for _, lesson := range lessons {
			var count int64
			cr.db.Model(&entities.LessonAttendance{}).
				Where("user_id = ? AND lesson_id = ?", userID, lesson.ID).
				Count(&count)
			learnedLessons[lesson.ID] = count > 0
		}
	}

	return modules, lessonsByModule, learnedLessons, nil
}

// Tạo hoặc cập nhật đánh giá khóa học của người dùng
func (cr *CourseRepository) CreateOrUpdateCourseRating(userID uuid.UUID, courseID uuid.UUID, star int, comment string) error {
	rating := entities.UserRate{
		UserID:   userID,
		CourseID: courseID,
		Star:     star,
		Comment:  comment,
	}

	// Tạo hoặc cập nhật đánh giá
	result := cr.db.Where("user_id = ? AND course_id = ?", userID, courseID).
		Assign(entities.UserRate{Star: star, Comment: comment}).
		FirstOrCreate(&rating)

	return result.Error
}

// Xóa khóa học
// Xóa khóa học sẽ tự động xóa các bản ghi liên quan trong bảng user_courses
func (cr *CourseRepository) DeleteCourse(courseID uuid.UUID) error {
	
	result := cr.db.Delete(&entities.Course{}, courseID)
	if result.Error != nil {
		return result.Error
	}


	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// Dùng để lấy danh sách người dùng đã đăng ký khóa học với phân trang
func (cr *CourseRepository) GetRegisteredUsersByCourseID(courseID uuid.UUID, offset, limit int, total *int64) ([]entities.Users, error) {
	var users []entities.Users

	err := cr.db.
		Table("users").
		Joins("JOIN user_courses ON user_courses.user_id = users.id").
		Where("user_courses.course_id = ?", courseID).
		Count(total).
		Offset(offset).
		Limit(limit).
		Find(&users).Error

	return users, err
}

// Đăng ký người dùng vào khóa học
func (cr *CourseRepository) RegisterUserToCourse(userID uuid.UUID, courseID uuid.UUID) error {
	var count int64
	err := cr.db.Model(&entities.UserCourse{}).
		Where("user_id = ? AND course_id = ?", userID, courseID).
		Count(&count).Error
	if err != nil {
		return err
	}

	if count > 0 {
		return gorm.ErrDuplicatedKey
	}


	userCourse := entities.UserCourse{
		UserID:   userID,
		CourseID: courseID,
	}

	return cr.db.Create(&userCourse).Error
}

// Kiểm tra xem khóa học có tồn tại hay không
func (cr *CourseRepository) CheckCourseExists(courseID uuid.UUID) (bool, error) {
	var count int64
	err := cr.db.Model(&entities.Course{}).
		Where("id = ?", courseID).
		Count(&count).Error

	return count > 0, err
}
