package services

import (
	"fmt"
	"sfit-platform-web-backend/internal/dtos"
	"time"
)

// lấy danh sách khóa học với các bộ lọc
func (s *CourseService) GetCourses(userID string, filter dtos.GetListCoursesForm) ([]dtos.GetListCoursesResponse, int64, error) {
	return s.courseRepo.GetCoursesV2(userID, filter)
}

func (s *CourseService) GetCourseDetailByIDV2(userID string, courseID string) (dtos.GetCourseDetailResponse, error) {
	// kiểm tra cache
	cachedCourseDetail, err := s.getCacheCourseDetailV2(userID, courseID)
	if err == nil {
		return cachedCourseDetail, nil
	}

	// nếu không có trong cache thì lấy từ db
	courseDetail, err := s.courseRepo.GetCourseDetailByIDV2(userID, courseID)
	if err != nil {
		return dtos.GetCourseDetailResponse{}, err
	}

	// lưu vào cache
	err = s.setCacheCourseDetailV2(userID, courseID, courseDetail)
	if err != nil {
		fmt.Println("Error setting cache for course detail:", err)
	}

	return courseDetail, nil
}

func (s *CourseService) buildCacheKeyCourseDetailV2(userID string, courseID string) string {
	return "course_detail_v2:" + userID + ":" + courseID
}

func (s *CourseService) setCacheCourseDetailV2(userID string, courseID string, courseDetail dtos.GetCourseDetailResponse) error {
	cacheKey := s.buildCacheKeyCourseDetailV2(userID, courseID)
	ttl := 60 * 60 * 24 // 1 ngày
	err := s.redisClient.Set(s.ctx, cacheKey, courseDetail, time.Second*time.Duration(ttl)).Err()
	if err != nil {
		return err
	}
	return nil
}

func (s *CourseService) getCacheCourseDetailV2(userID string, courseID string) (dtos.GetCourseDetailResponse, error) {
	cacheKey := s.buildCacheKeyCourseDetailV2(userID, courseID)
	var courseDetail dtos.GetCourseDetailResponse
	err := s.redisClient.Get(s.ctx, cacheKey).Scan(&courseDetail)
	if err != nil {
		return dtos.GetCourseDetailResponse{}, err
	}
	return courseDetail, nil
}
