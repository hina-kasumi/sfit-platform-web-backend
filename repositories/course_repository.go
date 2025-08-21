package repositories

import (
	"encoding/json"
	"fmt"
	"sfit-platform-web-backend/dtos"
	"sfit-platform-web-backend/entities"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CourseRepository struct {
	db *gorm.DB
}

func NewCourseRepository(db *gorm.DB) *CourseRepository {
	return &CourseRepository{
		db: db,
	}
}

// CreateNewCourse creates a new course in the database
func (r *CourseRepository) CreateNewCourse(course *entities.Course) error {
	return r.db.Exec(`
		INSERT INTO courses 
			(id, title, description, type, target, require, teachers, language, certificate, level, create_at, update_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`,
		course.ID,
		course.Title,
		course.Description,
		course.Type,
		stringArrayToPGArray(course.Target),
		stringArrayToPGArray(course.Require),
		stringArrayToPGArray(course.Teachers),
		course.Language,
		course.Certificate,
		course.Level,
		course.CreatedAt,
		course.UpdatedAt,
	).Error
}

// GetCourses retrieves courses with queries, pagination and user-specific data
func (r *CourseRepository) GetCourses(query dtos.CourseQuery) ([]dtos.CourseGeneralInformationResponse, int64, error) {
	courses, err := r.getCoursesList(query)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get courses: %w", err)
	}

	total, err := r.getCoursesCount(query)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count courses: %w", err)
	}

	return courses, total, nil
}

// GetCourseByID retrieves detailed course information by ID
func (r *CourseRepository) GetCourseByID(courseID, userID string) (*dtos.CourseDetailResponse, error) {
	courseDetail, err := r.getCourseBasicInfo(courseID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get course basic info: %w", err)
	}

	courseContent, err := r.getCourseContent(courseID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get course content: %w", err)
	}
	courseDetail.CourseContent = courseContent

	ratings, err := r.getCourseRatings(courseID)
	if err != nil {
		return nil, fmt.Errorf("failed to get course ratings: %w", err)
	}
	courseDetail.Rate = ratings

	return courseDetail, nil
}

// UpdateCourse updates an existing course in the database
func (r *CourseRepository) UpdateCourse(course *entities.Course) error {
	// check if course exists
	if err := r.EnsureCourseExists(course.ID); err != nil {
		return fmt.Errorf("course does not exist: %w", err)
	}

	return r.db.Exec(`
		UPDATE courses
		SET title = ?, description = ?, type = ?, target = ?, require = ?, teachers = ?, language = ?, certificate = ?, level = ?, update_at = ?
		WHERE id = ?
	`,
		course.Title,
		course.Description,
		course.Type,
		stringArrayToPGArray(course.Target),
		stringArrayToPGArray(course.Require),
		stringArrayToPGArray(course.Teachers),
		course.Language,
		course.Certificate,
		course.Level,
		course.UpdatedAt,
		course.ID,
	).Error
}

// Private helper methods

// stringArrayToPGArray converts Go string slice to PostgreSQL array format
func stringArrayToPGArray(arr []string) string {
	if len(arr) == 0 {
		return "{}"
	}
	escaped := make([]string, len(arr))
	for i, s := range arr {
		escaped[i] = `"` + strings.ReplaceAll(s, `"`, `\"`) + `"`
	}
	return "{" + strings.Join(escaped, ",") + "}"
}

// getCoursesList executes the main query to get courses with user-specific data
func (r *CourseRepository) getCoursesList(query dtos.CourseQuery) ([]dtos.CourseGeneralInformationResponse, error) {
	var rawCourses []dtos.CourseRaw

	mainQuery, args := r.buildCoursesQuery(query)

	err := r.db.Raw(mainQuery, args...).Scan(&rawCourses).Error
	if err != nil {
		return nil, err
	}

	return r.transformRawCourses(rawCourses)
}

// buildCoursesQuery constructs the SQL query for courses with dynamic querys
func (r *CourseRepository) buildCoursesQuery(query dtos.CourseQuery) (string, []interface{}) {
	offset := (query.Page - 1) * query.PageSize
	baseQuery := `
		WITH course_lessons AS (
			SELECT 
				m.course_id, 
				COUNT(l.id) AS lesson_count
			FROM modules m
			INNER JOIN lessons l ON l.module_id = m.id
			GROUP BY m.course_id
		),
		user_progress AS (
			SELECT 
				m.course_id,
				COUNT(CASE WHEN la.status = 'present' THEN 1 END) AS learned,
				COALESCE(SUM(la.timestamp / 60), 0) AS time_learn
			FROM lesson_attendances la
			INNER JOIN lessons l ON l.id = la.lesson_id
			INNER JOIN modules m ON m.id = l.module_id
			WHERE la.user_id = $1
			GROUP BY m.course_id
		),
		course_ratings AS (
			SELECT 
				courses_id AS course_id, 
				ROUND(AVG(star)::numeric, 2) AS rate
			FROM user_rates
			GROUP BY courses_id
		),
		course_tags AS (
			SELECT 
				tt.course_id, 
				ARRAY_AGG(t.id ORDER BY t.id) AS tag_list
			FROM tag_temps tt
			INNER JOIN tags t ON t.id = tt.tag_id
			GROUP BY tt.course_id
		)
		SELECT
			c.id,
			c.title,
			c.description,
			c.type,
			to_json(c.teachers)::json AS teachers, 
			COALESCE(cl.lesson_count, 0) AS number_lessons,
			COALESCE(up.time_learn, 0) AS time_learn,
			COALESCE(cr.rate, 0.0) AS rate,
			array_to_json(ct.tag_list) AS tags,
			COALESCE(up.learned, 0) AS learned_lessons,
			CASE WHEN up.learned > 0 THEN true ELSE false END AS registed
		FROM courses c
		LEFT JOIN course_lessons cl ON cl.course_id = c.id
		LEFT JOIN user_progress up ON up.course_id = c.id
		LEFT JOIN course_ratings cr ON cr.course_id = c.id
		LEFT JOIN course_tags ct ON ct.course_id = c.id`

	whereConditions, args := r.buildWhereConditions(query)
	args = append([]interface{}{query.UserID}, args...) // UserID first for CTE

	if len(whereConditions) > 0 {
		baseQuery += " WHERE " + strings.Join(whereConditions, " AND ")
	}

	paramIndex := len(args) + 1
	baseQuery += fmt.Sprintf(" ORDER BY c.id LIMIT $%d OFFSET $%d", paramIndex, paramIndex+1)
	args = append(args, query.PageSize, offset)

	return baseQuery, args
}

// buildWhereConditions constructs WHERE clause conditions based on query
func (r *CourseRepository) buildWhereConditions(query dtos.CourseQuery) ([]string, []interface{}) {
	var conditions []string
	var args []interface{}
	paramIndex := 2 // Start from $2 since $1 is reserved for UserID in CTE

	if query.Title != "" {
		conditions = append(conditions, fmt.Sprintf("c.title ILIKE '%%' || $%d || '%%'", paramIndex))
		args = append(args, query.Title)
		paramIndex++
	}

	if query.OnlyRegisted {
		if query.UserID != uuid.Nil {
			conditions = append(conditions, fmt.Sprintf(`EXISTS (
				SELECT 1 FROM lesson_attendances la 
				INNER JOIN lessons l ON l.id = la.lesson_id 
				INNER JOIN modules m ON m.id = l.module_id 
				WHERE la.user_id = $%d AND m.course_id = c.id
			)`, paramIndex))
			args = append(args, query.UserID)
			paramIndex++
		} else {
			conditions = append(conditions, "false")
		}
	}

	if query.CourseType != "" {
		conditions = append(conditions, fmt.Sprintf("c.type = $%d", paramIndex))
		args = append(args, query.CourseType)
		paramIndex++
	}

	if query.Level != "" {
		conditions = append(conditions, fmt.Sprintf("c.level = $%d", paramIndex))
		args = append(args, query.Level)
		paramIndex++
	}

	return conditions, args
}

// transformRawCourses converts raw database results to response DTOs
func (r *CourseRepository) transformRawCourses(rawCourses []dtos.CourseRaw) ([]dtos.CourseGeneralInformationResponse, error) {
	var courses []dtos.CourseGeneralInformationResponse

	for _, raw := range rawCourses {
		var teachers []string
		var tags []string

		if err := json.Unmarshal(raw.Teachers, &teachers); err != nil {
			return nil, fmt.Errorf("failed to parse teachers: %w", err)
		}

		if err := json.Unmarshal(raw.Tags, &tags); err != nil {
			return nil, fmt.Errorf("failed to parse tags: %w", err)
		}

		courses = append(courses, dtos.CourseGeneralInformationResponse{
			ID:             raw.ID,
			Title:          raw.Title,
			Description:    raw.Description,
			Type:           raw.Type,
			Teachers:       teachers,
			NumberLessons:  raw.NumberLessons,
			TimeLearn:      raw.TimeLearn,
			Rate:           raw.Rate,
			Tags:           tags,
			LearnedLessons: raw.LearnedLessons,
			Registed:       raw.Registed,
		})
	}

	return courses, nil
}

// getCoursesCount gets the total count of courses matching the query
func (r *CourseRepository) getCoursesCount(query dtos.CourseQuery) (int64, error) {
	var total int64

	countQuery := "SELECT COUNT(*) FROM courses c"
	conditions, args := r.buildCountWhereConditions(query)

	if len(conditions) > 0 {
		countQuery += " WHERE " + strings.Join(conditions, " AND ")
	}

	err := r.db.Raw(countQuery, args...).Scan(&total).Error
	return total, err
}

// buildCountWhereConditions builds WHERE conditions for count query
func (r *CourseRepository) buildCountWhereConditions(query dtos.CourseQuery) ([]string, []interface{}) {
	var conditions []string
	var args []interface{}
	paramIndex := 1

	if query.Title != "" {
		conditions = append(conditions, fmt.Sprintf("c.title ILIKE '%%' || $%d || '%%'", paramIndex))
		args = append(args, query.Title)
		paramIndex++
	}

	if query.OnlyRegisted {
		if query.UserID != uuid.Nil {
			conditions = append(conditions, fmt.Sprintf(`EXISTS (
				SELECT 1 FROM lesson_attendances la 
				INNER JOIN lessons l ON l.id = la.lesson_id 
				INNER JOIN modules m ON m.id = l.module_id 
				WHERE la.user_id = $%d AND m.course_id = c.id
			)`, paramIndex))
			args = append(args, query.UserID)
			paramIndex++
		} else {
			conditions = append(conditions, "false")
		}
	}

	if query.CourseType != "" {
		conditions = append(conditions, fmt.Sprintf("c.type = $%d", paramIndex))
		args = append(args, query.CourseType)
		paramIndex++
	}

	if query.Level != "" {
		conditions = append(conditions, fmt.Sprintf("c.level = $%d", paramIndex))
		args = append(args, query.Level)
		paramIndex++
	}

	return conditions, args
}

// getCourseBasicInfo retrieves basic course information
func (r *CourseRepository) getCourseBasicInfo(courseID, userID string) (*dtos.CourseDetailResponse, error) {
	var courseRaw dtos.CourseDetailRaw
	// COALESCE(SUM(
	// 	CASE
	// 		WHEN l.lesson_type = 'Quiz' THEN (l.quiz_content->>'duration')::int
	// 		WHEN l.lesson_type = 'Online' THEN (l.online_content->>'duration')::int
	// 		WHEN l.lesson_type = 'Offline' THEN (l.offline_content->>'duration')::int
	// 		ELSE 0
	// 	END
	// ), 0) as total_time
	query := `
		WITH course_stats AS (
			SELECT 
				c.id,
				c.title,
				c.description,
				c.type,
				c.level,
				array_to_json(c.teachers) as teachers,
				array_to_json(c.target) as target,
				array_to_json(c.require) as require,
				c.total_time,
				c.language,
				c.update_at,
				CASE WHEN fc.user_id IS NOT NULL THEN true ELSE false END as like,
				COALESCE(AVG(ur.star), 0) as star,
				COUNT(DISTINCT uc.user_id) as total_registered,
				COUNT(DISTINCT l.id) as total_lessons
			FROM courses c
			LEFT JOIN favorite_courses fc ON c.id = fc.course_id AND fc.user_id = $2::uuid
			LEFT JOIN user_rates ur ON c.id = ur.courses_id
			LEFT JOIN user_courses uc ON c.id = uc.course_id
			LEFT JOIN modules m ON c.id = m.course_id
			LEFT JOIN lessons l ON m.id = l.module_id
			WHERE c.id = $1::uuid
			GROUP BY c.id, c.title, c.description, c.type, c.level, c.teachers, 
					c.target, c.require, c.language, c.update_at, fc.user_id
		),
		course_tags AS (
			SELECT 
				tt.course_id,
				json_agg(t.id) as tags
			FROM tag_temps tt
			JOIN tags t ON tt.tag_id = t.id
			WHERE tt.course_id = $1::uuid
			GROUP BY tt.course_id
		)
		SELECT 
			cs.title,
			cs.description,
			cs.like,
			cs.type,
			cs.level,
			cs.teachers,
			cs.star,
			cs.total_lessons,
			COALESCE(ct.tags, '[]'::json) as tags,
			cs.target,
			cs.require,
			cs.total_time,
			cs.total_registered,
			cs.update_at as updated_at,
			cs.language
		FROM course_stats cs
		LEFT JOIN course_tags ct ON cs.id = ct.course_id`

	err := r.db.Raw(query, courseID, userID).Scan(&courseRaw).Error
	if err != nil {
		return nil, err
	}

	return r.transformCourseDetail(courseRaw)
}

// transformCourseDetail converts raw course data to response DTO
func (r *CourseRepository) transformCourseDetail(raw dtos.CourseDetailRaw) (*dtos.CourseDetailResponse, error) {
	course := &dtos.CourseDetailResponse{
		Title:          raw.Title,
		Description:    raw.Description,
		Like:           raw.Like,
		Type:           raw.Type,
		Level:          raw.Level,
		Star:           raw.Star,
		TotalTime:      raw.TotalTime,
		TotalLessons:   raw.TotalLessons,
		TotalRegitered: raw.TotalRegistered,
		UpdatedAt:      raw.UpdatedAt,
		Language:       raw.Language,
	}

	// Unmarshal JSON arrays
	if len(raw.Teachers) > 0 {
		if err := json.Unmarshal(raw.Teachers, &course.Teachers); err != nil {
			return nil, fmt.Errorf("failed to parse teachers: %w", err)
		}
	}

	if len(raw.Tags) > 0 {
		if err := json.Unmarshal(raw.Tags, &course.Tags); err != nil {
			return nil, fmt.Errorf("failed to parse tags: %w", err)
		}
	}

	if len(raw.Target) > 0 {
		if err := json.Unmarshal(raw.Target, &course.Target); err != nil {
			return nil, fmt.Errorf("failed to parse target: %w", err)
		}
	}

	if len(raw.Require) > 0 {
		if err := json.Unmarshal(raw.Require, &course.Require); err != nil {
			return nil, fmt.Errorf("failed to parse require: %w", err)
		}
	}

	return course, nil
}

// getCourseContent retrieves course modules and lessons
func (r *CourseRepository) getCourseContent(courseID, userID string) ([]dtos.CourseContentResponse, error) {
	type ModuleRaw struct {
		ID          string `json:"id"`
		ModuleTitle string `json:"module_title"`
		// TotalTime   int    `json:"total_time"`
	}

	var modules []ModuleRaw

	moduleQuery := `
		SELECT 
			m.id,
			m.module_title
		FROM modules m
		WHERE m.course_id = $1::uuid
		ORDER BY m.create_at ASC`

	// moduleQuery := `
	// 	SELECT
	// 		m.id,
	// 		m.module_title,
	// 		COALESCE(SUM(
	// 			CASE
	// 				WHEN l.lesson_type = 'Quiz' THEN (l.quiz_content->>'duration')::int
	// 				WHEN l.lesson_type = 'Online' THEN (l.online_content->>'duration')::int
	// 				WHEN l.lesson_type = 'Offline' THEN (l.offline_content->>'duration')::int
	// 				ELSE 0
	// 			END
	// 		), 0) as total_time
	// 	FROM modules m
	// 	LEFT JOIN lessons l ON m.id = l.module_id
	// 	WHERE m.course_id = $1::uuid
	// 	GROUP BY m.id, m.module_title
	// 	ORDER BY m.create_at ASC`

	err := r.db.Raw(moduleQuery, courseID).Scan(&modules).Error
	if err != nil {
		return nil, err
	}

	var courseContent []dtos.CourseContentResponse

	for _, module := range modules {
		lessons, err := r.getModuleLessons(module.ID, userID)
		if err != nil {
			return nil, err
		}

		courseContent = append(courseContent, dtos.CourseContentResponse{
			ID:          module.ID,
			ModuleTitle: module.ModuleTitle,
			// TotalTime:   module.TotalTime,
			Lessons: lessons,
		})
	}

	return courseContent, nil
}

// getModuleLessons retrieves lessons for a specific module
func (r *CourseRepository) getModuleLessons(moduleID, userID string) ([]dtos.LessonResponse, error) {
	type LessonRaw struct {
		ID        string `json:"id"`
		Title     string `json:"title"`
		Learned   bool   `json:"learned"`
		StudyTime int    `json:"study_time"`
	}

	var lessons []LessonRaw

	lessonQuery := `
		SELECT 
			l.id,
			l.title AS title,

			CASE WHEN la.user_id IS NOT NULL THEN true ELSE false END as learned,
			CASE 
				WHEN l.lesson_type = 'Quiz' THEN (l.quiz_content->>'duration')::int
				WHEN l.lesson_type = 'Online' THEN (l.online_content->>'duration')::int
				WHEN l.lesson_type = 'Offline' THEN (l.offline_content->>'duration')::int
				ELSE 0
			END as study_time
		FROM lessons l
		LEFT JOIN lesson_attendances la ON l.id = la.lesson_id AND la.user_id = $2::uuid
		WHERE l.module_id = $1::uuid
		ORDER BY l.create_at ASC`

	err := r.db.Raw(lessonQuery, moduleID, userID).Scan(&lessons).Error
	if err != nil {
		return nil, err
	}

	var result []dtos.LessonResponse
	for _, lesson := range lessons {
		result = append(result, dtos.LessonResponse{
			ID:        lesson.ID,
			Title:     lesson.Title,
			Learned:   lesson.Learned,
			StudyTime: lesson.StudyTime,
		})
	}

	return result, nil
}

// getCourseRatings retrieves course ratings and reviews
func (r *CourseRepository) getCourseRatings(courseID string) ([]dtos.RateResponse, error) {
	type RatingRaw struct {
		Name      string    `json:"name"`
		Comment   string    `json:"comment"`
		Star      float64   `json:"star"`
		CreatedAt time.Time `json:"created_at"`
	}

	var ratings []RatingRaw

	ratingQuery := `
		SELECT 
			u.username as name,
			ur.comment,
			ur.star,
			ur.created_at
		FROM user_rates ur
		JOIN users u ON ur.user_id = u.id
		WHERE ur.courses_id = $1::uuid
		ORDER BY ur.created_at DESC`

	err := r.db.Raw(ratingQuery, courseID).Scan(&ratings).Error
	if err != nil {
		return nil, err
	}

	var result []dtos.RateResponse
	for _, rating := range ratings {
		result = append(result, dtos.RateResponse{
			Name:      rating.Name,
			Comment:   rating.Comment,
			Star:      rating.Star,
			CreatedAt: rating.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return result, nil
}

func (r *CourseRepository) EnsureCourseExists(courseID uuid.UUID) error {
	var dummy int
	err := r.db.
		Model(&entities.Course{}).
		Select("1").
		Where("id = ?", courseID).
		Limit(1).
		Scan(&dummy).Error

	if err != nil {
		return fmt.Errorf("failed to check course existence: %w", err)
	}
	if dummy == 0 {
		return fmt.Errorf("course with ID %s does not exist", courseID)
	}
	return nil
}

func (r *CourseRepository) GetListUserCompleteCourses(query dtos.CourseQuery) ([]dtos.UserListItem, int64, error) {
	// Ensure course exists
	if err := r.EnsureCourseExists(query.CourseID); err != nil {
		return nil, 0, fmt.Errorf("course does not exist: %w", err)
	}

	// Get total lessons for the course
	var totalLessons int
	totalLessonQuery := `
		SELECT COUNT(DISTINCT l.id) AS total_lessons
		FROM courses c
		LEFT JOIN modules m ON m.course_id = c.id
		LEFT JOIN lessons l ON l.module_id = m.id
		WHERE c.id = ?
	`
	if err := r.db.Raw(totalLessonQuery, query.CourseID).Scan(&totalLessons).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to get total lessons: %w", err)
	}

	// If no lessons, return empty result
	if totalLessons == 0 {
		return []dtos.UserListItem{}, 0, nil
	}

	// Query users who completed all lessons (status = 'present')
	var users []dtos.UserListItem
	mainQuery := `
		SELECT 
			u.id,
			u.username,
			u.email
		FROM users u
		JOIN lesson_attendances la ON u.id = la.user_id
		JOIN lessons l ON l.id = la.lesson_id
		JOIN modules m ON m.id = l.module_id
		WHERE m.course_id = ? AND la.status = 'present'
		GROUP BY u.id, u.username, u.email
		HAVING COUNT(DISTINCT la.lesson_id) = ?
		LIMIT ? OFFSET ?
	`

	offset := (query.Page - 1) * query.PageSize
	if err := r.db.Raw(mainQuery, query.CourseID, totalLessons, query.PageSize, offset).Scan(&users).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to get users: %w", err)
	}

	// Query total count for pagination
	var totalUser int64
	countQuery := `
		SELECT COUNT(*) FROM (
			SELECT u.id
			FROM users u
			JOIN lesson_attendances la ON u.id = la.user_id
			JOIN lessons l ON l.id = la.lesson_id
			JOIN modules m ON m.id = l.module_id
			WHERE m.course_id = ? AND la.status = 'present'
			GROUP BY u.id
			HAVING COUNT(DISTINCT la.lesson_id) = ?
		) AS completed_users
	`

	if err := r.db.Raw(countQuery, query.CourseID, totalLessons).Scan(&totalUser).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count users: %w", err)
	}

	return users, totalUser, nil
}

func (r *UserCourseRepository) GetUserProgressInCourse(courseID, userID string) (int, int, error) {
	var res dtos.GetUserProgressInCourseResponse

	query := `
        SELECT 
            (SELECT COALESCE(SUM(la.timestamp / 60), 0)
             FROM lesson_attendances la
             JOIN lessons l ON l.id = la.lesson_id
             JOIN modules m ON m.id = l.module_id
             WHERE m.course_id = ? AND la.user_id = ?) AS learned,
            (SELECT COUNT(l.id)
             FROM lessons l
             JOIN modules m ON m.id = l.module_id
             WHERE m.course_id = ?) AS total_lessons
    `

	err := r.db.Raw(query, courseID, userID, courseID).Scan(&res).Error
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get progress: %w", err)
	}

	return res.Learned, res.TotalLessons, nil
}

// Lấy danh sách khóa học đã đăng ký của người dùng với phân trang
func (cr *CourseRepository) GetCoursesByUserIDWithPagination(
	userID string,
	offset, limit int,
	total *int64,
) ([]dtos.CourseGeneralInformationResponse, error) {

	var courses []dtos.CourseGeneralInformationResponse

	// Đếm tổng
	if err := cr.db.
		Table("courses").
		Joins("JOIN user_courses ON user_courses.course_id = courses.id").
		Where("user_courses.user_id = ?", userID).
		Count(total).Error; err != nil {
		return nil, err
	}

	// Lấy dữ liệu: convert text[] thành JSON để map an toàn vào []string
	query := `
        SELECT 
            courses.id,
            courses.title,
            courses.description,
            courses.level AS type,
            COALESCE(ARRAY(SELECT unnest(teachers)), '{}') AS teachers,
            0 AS number_lessons,
            0 AS time_learn,
            5 AS rate,
            '{}'::text[] AS tags,
            0 AS learned_lessons,
            TRUE AS registed
        FROM courses
        JOIN user_courses ON user_courses.course_id = courses.id
        WHERE user_courses.user_id = ?
        OFFSET ? LIMIT ?
    `

	if err := cr.db.Raw(query, userID, offset, limit).Scan(&courses).Error; err != nil {
		return nil, err
	}

	return courses, nil
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
	result := cr.db.Where("user_id = ? AND courses_id = ?", userID, courseID).
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

// Helper: chuyển string JSON sang []string
func ParseStringArray(data string) []string {
	var arr []string
	if err := json.Unmarshal([]byte(data), &arr); err != nil {
		return []string{}
	}
	return arr
}

func (cr *CourseRepository) UpdateTotalTime(moduleID uuid.UUID, time int) error {
	var courseID string
	err := cr.db.Model(&entities.Module{}).Where("id = ?", moduleID).Select("course_id").Scan(&courseID).Error
	if err != nil {
		return err
	}

	if time < 0 {
		time = 0
	}
	return cr.db.Model(&entities.Course{}).
		Where("id = ?", courseID).
		Update("total_time", gorm.Expr("GREATEST(total_time + ?, 0)", time)).Error
}

func (cr *CourseRepository) UpdateTotalLesson(courseID uuid.UUID, lessonCount int) (*entities.Course, error) {
	var course entities.Course
	err := cr.db.Model(&course).
		Where("id = ?", courseID).
		Update("total_lessons", gorm.Expr("GREATEST(total_lessons + ?, 0)", lessonCount)).Error
	if err != nil {
		return nil, err
	}
	return &course, nil

}

func (cr *CourseRepository) GetModuleByID(moduleID uuid.UUID) (*entities.Module, error) {
	var module entities.Module
	err := cr.db.Where("id = ?", moduleID).First(&module).Error
	if err != nil {
		return nil, err
	}
	return &module, nil
}

func (cr *CourseRepository) GetCourseUserCompletion(userID uuid.UUID) ([]string, error) {
	var courseIDs []string

	err := cr.db.Table("courses").
		Select("id").
		Where("total_lessons <= (?)",
			cr.db.Table("lesson_attendances").
				Select("COUNT(lesson_id)").
				Where("course_id = courses.id AND user_id = ?", userID),
		).Scan(&courseIDs).Error

	if err != nil {
		return nil, err
	}
	return courseIDs, nil
}

// func (cr *CourseRepository) GetNumberOfCompletedCourses(userID uuid.UUID) (int64) {
// 	var count int64
// 	err := cr.db.Model(&entities.UserCourse{}).
// 		Where("user_id = ? AND status = ?", userID, entities.).
// 		Count(&count).Error

// 	if err != nil {
// 		return 0
// 	}

// 	return count
// }
