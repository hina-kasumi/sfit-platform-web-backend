package repositories

import (
	"sfit-platform-web-backend/dtos"
	"sfit-platform-web-backend/entities"
)

func (r *CourseRepository) GetCoursesV2(userID string, filter dtos.GetListCoursesForm) ([]dtos.GetListCoursesResponse, int64, error) {
	query := r.db.Model(&entities.Course{}).Preload("Tags")
	if filter.Title != "" {
		query = query.Where("title ILIKE ?", "%"+filter.Title+"%")
	}
	if filter.Type != "" {
		query = query.Where("type ILIKE ?", filter.Type)
	}
	if filter.Level != "" {
		query = query.Where("level ILIKE ?", filter.Level)
	}
	if filter.UserStatus != "" && userID != "" {
		var coursesIDs []string
		r.db.Model(&entities.UserCourse{}).Select("course_id").Where("user_id = ? AND status = ?", userID, filter.UserStatus).Scan(&coursesIDs)
		query = query.Where("id IN ?", coursesIDs)
	}

	// thục hiện phân trang
	var count int64
	query.Count(&count)
	var courses []entities.Course
	offset := (filter.Page - 1) * filter.PageSize
	result := query.Offset(offset).Limit(filter.PageSize).Find(&courses)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	// map courses sang response
	var response []dtos.GetListCoursesResponse
	for _, course := range courses {
		var learnedLessons int64 = 0
		status := ""
		if userID != "" {
			err := r.db.Model(&entities.LessonAttendance{}).
				Where("user_id = ? AND course_id = ? AND status = ?", userID, course.ID, string(entities.Present)).
				Count(&learnedLessons).Error
			if err != nil {
				return nil, 0, err
			}

			// lấy trạng thái của user với khóa học
			var userCourse entities.UserCourse
			err = r.db.Model(&entities.UserCourse{}).
				Where("user_id = ? AND course_id = ?", userID, course.ID).
				First(&userCourse).Error
			if err == nil {
				status = string(userCourse.Status)
			}
		}
		var tags []string
		for _, tag := range course.Tags {
			tags = append(tags, tag.TagID)
		}

		response = append(response, dtos.GetListCoursesResponse{
			ID:             course.ID,
			Title:          course.Title,
			Description:    course.Description,
			Type:           course.Type,
			Teachers:       course.Teachers,
			TimeLearn:      course.TotalTime,
			Rate:           course.Rate,
			Tags:           tags,
			TotalLessons:   course.TotalLessons,
			LearnedLessons: int(learnedLessons),
			UserStatus:     status, // cần bổ sung logic lấy trạng thái của user với khóa học
		})
	}

	return response, count, nil
}

// GetCourseDetailByIDV2 lấy chi tiết khóa học theo ID (phiên bản 2)
func (r *CourseRepository) GetCourseDetailByIDV2(userID string, courseID string) (dtos.GetCourseDetailResponse, error) {
	var course entities.Course
	result := r.db.Preload("Tags").Preload("Modules").Preload("Modules.Lessons").
		Where("id = ?", courseID).First(&course)
	if result.Error != nil {
		return dtos.GetCourseDetailResponse{}, result.Error
	}

	// kiểm tra user đã thích khóa học chưa
	var liked bool = false
	var status string = ""
	if userID != "" {
		var count int64
		r.db.Model(&entities.FavoriteCourse{}).Where("user_id = ? AND course_id = ?", userID, courseID).Count(&count)
		if count > 0 {
			liked = true
		}
		var userCourse entities.UserCourse
		err := r.db.Model(&entities.UserCourse{}).Where("user_id = ? AND course_id = ?", userID, courseID).First(&userCourse).Error
		if err == nil {
			status = string(userCourse.Status)
		}
	}

	// lấy đánh giá
	var rate []entities.UserRate
	err := r.db.Model(&entities.UserRate{}).Where("courses_id = ?", courseID).Find(&rate).Error
	if err != nil {
		return dtos.GetCourseDetailResponse{}, err
	}

	// map tags
	var tags []string
	for _, tag := range course.Tags {
		tags = append(tags, tag.TagID)
	}

	// map modules
	modules := make([]dtos.ModuleResponse, len(course.Modules))
	for i, module := range course.Modules {
		modules[i] = dtos.ModuleResponse{
			ID:           module.ID,
			Title:        module.Title,
			CreatedAt:    module.CreatedAt,
			UpdatedAt:    module.UpdatedAt,
			TotalTime:    module.TotalTime,
			TotalLessons: module.TotalLessons,
			Lessons:      []dtos.LessonResponse{},
		}
	}
	//map lessons
	for i, module := range course.Modules {
		for _, lesson := range module.Lessons {
			// kiểm tra xem user đã học bài học này chưa
			var lessonAttendance *entities.LessonAttendance
			err := r.db.Model(&entities.LessonAttendance{}).
				Where("user_id = ? AND lesson_id = ?", userID, lesson.ID).
				First(&lessonAttendance).Error
			if err != nil && err.Error() != "record not found" {
				return dtos.GetCourseDetailResponse{}, err
			}

			status := ""
			if lessonAttendance != nil {
				status = string(lessonAttendance.Status)
			}
			modules[i].Lessons = append(modules[i].Lessons, dtos.LessonResponse{
				ID:        lesson.ID.String(),
				Title:     lesson.Title,
				Type:      string(lesson.Type),
				StudyTime: lesson.Duration,
				Learned:   lessonAttendance != nil,
				Status:    status,
			})
		}
	}
	course.Modules = nil // giải phóng bộ nhớ

	return dtos.GetCourseDetailResponse{
		ID:           course.ID,
		Title:        course.Title,
		Liked:        liked,
		Description:  course.Description,
		Type:         course.Type,
		Level:        course.Level,
		Teachers:     course.Teachers,
		Star:         course.Rate,
		TotalLessons: course.TotalLessons,
		Tags:         tags,
		Target:       course.Target,
		Require:      course.Require,
		Language:     course.Language,
		TotalTime:    course.TotalTime,
		TotalLearned: len(course.UsersCourse),
		UpdatedAt:    course.UpdatedAt,
		CreatedAt:    course.CreatedAt,
		Certificate:  course.Certificate,
		Rate:         rate,
		Status:       status,
		Modules:      modules,
	}, nil
}
