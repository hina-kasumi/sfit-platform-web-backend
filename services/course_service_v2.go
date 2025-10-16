package services

import "sfit-platform-web-backend/dtos"

// lấy danh sách khóa học với các bộ lọc
func (s *CourseService) GetCourses(userID string, filter dtos.GetListCoursesForm) ([]dtos.GetListCoursesResponse, int64, error) {
	return s.courseRepo.GetCoursesV2(userID, filter)
}

func (s *CourseService) GetCourseDetailByIDV2(userID string, courseID string) (dtos.GetCourseDetailResponse, error) {
	return s.courseRepo.GetCourseDetailByIDV2(userID, courseID)
}
