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

// GetCourses retrieves courses with filters, pagination and user-specific data
func (r *CourseRepository) GetCourses(filter dtos.CourseFilter, pageSize, offset int) ([]dtos.CourseGeneralInformationResponse, int, error) {
	courses, err := r.getCoursesList(filter, pageSize, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get courses: %w", err)
	}

	total, err := r.getCoursesCount(filter)
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
func (r *CourseRepository) getCoursesList(filter dtos.CourseFilter, pageSize, offset int) ([]dtos.CourseGeneralInformationResponse, error) {
	var rawCourses []dtos.CourseRaw

	query, args := r.buildCoursesQuery(filter, pageSize, offset)

	err := r.db.Raw(query, args...).Scan(&rawCourses).Error
	if err != nil {
		return nil, err
	}

	return r.transformRawCourses(rawCourses)
}

// buildCoursesQuery constructs the SQL query for courses with dynamic filters
func (r *CourseRepository) buildCoursesQuery(filter dtos.CourseFilter, pageSize, offset int) (string, []interface{}) {
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

	whereConditions, args := r.buildWhereConditions(filter)
	args = append([]interface{}{filter.UserID}, args...) // UserID first for CTE

	if len(whereConditions) > 0 {
		baseQuery += " WHERE " + strings.Join(whereConditions, " AND ")
	}

	paramIndex := len(args) + 1
	baseQuery += fmt.Sprintf(" ORDER BY c.id LIMIT $%d OFFSET $%d", paramIndex, paramIndex+1)
	args = append(args, pageSize, offset)

	return baseQuery, args
}

// buildWhereConditions constructs WHERE clause conditions based on filter
func (r *CourseRepository) buildWhereConditions(filter dtos.CourseFilter) ([]string, []interface{}) {
	var conditions []string
	var args []interface{}
	paramIndex := 2 // Start from $2 since $1 is reserved for UserID in CTE

	if filter.Title != "" {
		conditions = append(conditions, fmt.Sprintf("c.title ILIKE '%%' || $%d || '%%'", paramIndex))
		args = append(args, filter.Title)
		paramIndex++
	}

	if filter.OnlyRegisted {
		if filter.UserID != uuid.Nil {
			conditions = append(conditions, fmt.Sprintf(`EXISTS (
				SELECT 1 FROM lesson_attendances la 
				INNER JOIN lessons l ON l.id = la.lesson_id 
				INNER JOIN modules m ON m.id = l.module_id 
				WHERE la.user_id = $%d AND m.course_id = c.id
			)`, paramIndex))
			args = append(args, filter.UserID)
			paramIndex++
		} else {
			conditions = append(conditions, "false")
		}
	}

	if filter.CourseType != "" {
		conditions = append(conditions, fmt.Sprintf("c.type = $%d", paramIndex))
		args = append(args, filter.CourseType)
		paramIndex++
	}

	if filter.Level != "" {
		conditions = append(conditions, fmt.Sprintf("c.level = $%d", paramIndex))
		args = append(args, filter.Level)
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

// getCoursesCount gets the total count of courses matching the filter
func (r *CourseRepository) getCoursesCount(filter dtos.CourseFilter) (int, error) {
	var total int

	countQuery := "SELECT COUNT(*) FROM courses c"
	conditions, args := r.buildCountWhereConditions(filter)

	if len(conditions) > 0 {
		countQuery += " WHERE " + strings.Join(conditions, " AND ")
	}

	err := r.db.Raw(countQuery, args...).Scan(&total).Error
	return total, err
}

// buildCountWhereConditions builds WHERE conditions for count query
func (r *CourseRepository) buildCountWhereConditions(filter dtos.CourseFilter) ([]string, []interface{}) {
	var conditions []string
	var args []interface{}
	paramIndex := 1

	if filter.Title != "" {
		conditions = append(conditions, fmt.Sprintf("c.title ILIKE '%%' || $%d || '%%'", paramIndex))
		args = append(args, filter.Title)
		paramIndex++
	}

	if filter.OnlyRegisted {
		if filter.UserID != uuid.Nil {
			conditions = append(conditions, fmt.Sprintf(`EXISTS (
				SELECT 1 FROM lesson_attendances la 
				INNER JOIN lessons l ON l.id = la.lesson_id 
				INNER JOIN modules m ON m.id = l.module_id 
				WHERE la.user_id = $%d AND m.course_id = c.id
			)`, paramIndex))
			args = append(args, filter.UserID)
			paramIndex++
		} else {
			conditions = append(conditions, "false")
		}
	}

	if filter.CourseType != "" {
		conditions = append(conditions, fmt.Sprintf("c.type = $%d", paramIndex))
		args = append(args, filter.CourseType)
		paramIndex++
	}

	if filter.Level != "" {
		conditions = append(conditions, fmt.Sprintf("c.level = $%d", paramIndex))
		args = append(args, filter.Level)
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
			CASE 
				WHEN l.lesson_type = 'Quiz' THEN 'Quiz Lesson'
				WHEN l.lesson_type = 'Online' THEN l.online_content->>'title'
				WHEN l.lesson_type = 'Offline' THEN 'Offline Lesson'
				ELSE 'Unknown'
			END as title,
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

func (r *CourseRepository) GetListUserCompleteCourses(filter dtos.CourseFilter, pageSize, offset int) ([]dtos.UserListItem, int64, error) {
	// Ensure course exists
	if err := r.EnsureCourseExists(filter.CourseID); err != nil {
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
	if err := r.db.Raw(totalLessonQuery, filter.CourseID).Scan(&totalLessons).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to get total lessons: %w", err)
	}

	// If no lessons, return empty result
	if totalLessons == 0 {
		return []dtos.UserListItem{}, 0, nil
	}

	// Query users who completed all lessons (status = 'present')
	var users []dtos.UserListItem
	query := `
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

	if err := r.db.Raw(query, filter.CourseID, totalLessons, pageSize, offset).Scan(&users).Error; err != nil {
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

	if err := r.db.Raw(countQuery, filter.CourseID, totalLessons).Scan(&totalUser).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count users: %w", err)
	}

	return users, totalUser, nil
}
