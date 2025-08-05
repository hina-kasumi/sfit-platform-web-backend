package dtos

type CourseInfoResponse struct {
	ID            string   `json:"id"`
	Title         string   `json:"title"`
	Description   string   `json:"description"`
	Type          string   `json:"type"`
	Teachers      []string `json:"teachers"`
	TimeLearn     int      `json:"time_learn"`
	Rate          float64  `json:"rate"`
	Tags          []string `json:"tags"`
	TotalLesson   int      `json:"total_lesson"`
	LearnedLesson int      `json:"learned_lesson"`
	Registed      bool     `json:"registed"`
}

type CourseListResponse struct {
	Courses  []CourseInfoResponse `json:"courses"`
	Page     int                  `json:"page"`
	PageSize int                  `json:"pageSize"`
	Total    int64                `json:"total"`
}

type LessonInfo struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Learned   bool   `json:"learned"`
	StudyTime int    `json:"study_time"`
}

type ModuleInfo struct {
	ID          string       `json:"id"`
	ModuleTitle string       `json:"module_title"`
	TotalTime   int          `json:"total_time"`
	Lessons     []LessonInfo `json:"lessons"`
}

type CourseLessonsResponse []ModuleInfo

type CourseRateRequest struct {
	Course  string `json:"course" binding:"required"`
	Star    int    `json:"star" binding:"required,min=1,max=5"`
	Comment string `json:"comment"`
}

type RegisteredUserInfo struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type RegisteredUsersResponse struct {
	Users    []RegisteredUserInfo `json:"users"`
	Page     int                  `json:"page"`
	PageSize int                  `json:"pageSize"`
	Total    int64                `json:"total"`
}

type CourseRegisterRequest struct {
	CourseID string `json:"course_id" binding:"required"`
}
