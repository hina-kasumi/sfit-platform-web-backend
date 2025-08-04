package repositories

import (
	// "encoding/json"
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

// Tạo course mới
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

func (r *CourseRepository) CreateNewCourse(course *entities.Course) error {
	// Trick lỏ :>
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

func (r *CourseRepository) GetCourses(filter dtos.CourseFilter, pageSize, offset int) ([]dtos.CourseGeneralInformationResponse, int, error) {
	// var courses []dtos.CourseGeneralInformationResponse
	var rawCourses []dtos.CourseRaw

	// Base query với CTE
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
				COALESCE(SUM(EXTRACT(EPOCH FROM la.timestamp)::int / 60), 0) AS time_learn
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

	// Build WHERE conditions và arguments
	var whereConditions []string
	var args []interface{}

	// Luôn thêm userID làm tham số đầu tiên cho CTE
	args = append(args, filter.UserID)
	paramIndex := 2 // Bắt đầu từ $2

	// Build WHERE conditions
	if filter.Title != "" {
		whereConditions = append(whereConditions, fmt.Sprintf("c.title ILIKE '%%' || $%d || '%%'", paramIndex))
		args = append(args, filter.Title)
		paramIndex++
	}

	if filter.OnlyRegisted {
		if filter.UserID != uuid.Nil {
			whereConditions = append(whereConditions, fmt.Sprintf(`EXISTS (
				SELECT 1 FROM lesson_attendances la 
				INNER JOIN lessons l ON l.id = la.lesson_id 
				INNER JOIN modules m ON m.id = l.module_id 
				WHERE la.user_id = $%d AND m.course_id = c.id
			)`, paramIndex))
			args = append(args, filter.UserID)
			paramIndex++
		} else {
			whereConditions = append(whereConditions, "false") // Không có userID
		}
	}

	if filter.CourseType != "" {
		whereConditions = append(whereConditions, fmt.Sprintf("c.type = $%d", paramIndex))
		args = append(args, filter.CourseType)
		paramIndex++
	}

	if filter.Level != "" {
		whereConditions = append(whereConditions, fmt.Sprintf("c.level = $%d", paramIndex))
		args = append(args, filter.Level)
		paramIndex++
	}

	// Tạo main query
	mainQuery := baseQuery
	if len(whereConditions) > 0 {
		mainQuery += " WHERE " + strings.Join(whereConditions, " AND ")
	}
	mainQuery += fmt.Sprintf(" ORDER BY c.id LIMIT $%d OFFSET $%d", paramIndex, paramIndex+1)

	// Thêm pageSize và offset vào args
	mainArgs := append(args, pageSize, offset)

	// Execute main query
	// err := r.db.Raw(mainQuery, mainArgs...).Scan(&courses).Error
	err := r.db.Raw(mainQuery, mainArgs...).Scan(&rawCourses).Error
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get courses: %w", err)
	}

	// Chuyển đổi rawCourses sang courses
	var courses []dtos.CourseGeneralInformationResponse
	for _, raw := range rawCourses {
		var teachers []string
		var tags []string

		// Unmarshal teachers JSON
		if err := json.Unmarshal(raw.Teachers, &teachers); err != nil {
			return nil, 0, fmt.Errorf("failed to parse teachers: %w", err)
		}

		// Unmarshal tags JSON
		if err := json.Unmarshal(raw.Tags, &tags); err != nil {
			return nil, 0, fmt.Errorf("failed to parse tags: %w", err)
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

	// Tính tổng số khóa học
	var total int
	countQuery := "SELECT COUNT(*) FROM courses c"
	var countArgs []interface{}
	var countConditions []string
	countParamIndex := 1

	// Build count conditions (giống main query nhưng không có userID cho CTE)
	if filter.Title != "" {
		countConditions = append(countConditions, fmt.Sprintf("c.title ILIKE '%%' || $%d || '%%'", countParamIndex))
		countArgs = append(countArgs, filter.Title)
		countParamIndex++
	}

	if filter.OnlyRegisted {
		if filter.UserID != uuid.Nil {
			countConditions = append(countConditions, fmt.Sprintf(`EXISTS (
				SELECT 1 FROM lesson_attendances la 
				INNER JOIN lessons l ON l.id = la.lesson_id 
				INNER JOIN modules m ON m.id = l.module_id 
				WHERE la.user_id = $%d AND m.course_id = c.id
			)`, countParamIndex))
			countArgs = append(countArgs, filter.UserID)
			countParamIndex++
		} else {
			countConditions = append(countConditions, "false")
		}
	}

	if filter.CourseType != "" {
		countConditions = append(countConditions, fmt.Sprintf("c.type = $%d", countParamIndex))
		countArgs = append(countArgs, filter.CourseType)
		countParamIndex++
	}

	if filter.Level != "" {
		countConditions = append(countConditions, fmt.Sprintf("c.level = $%d", countParamIndex))
		countArgs = append(countArgs, filter.Level)
		countParamIndex++
	}

	if len(countConditions) > 0 {
		countQuery += " WHERE " + strings.Join(countConditions, " AND ")
	}

	// Execute count query
	err = r.db.Raw(countQuery, countArgs...).Scan(&total).Error
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count courses: %w", err)
	}

	return courses, total, nil
}

//Mới ------------------------------------------------------------

/*
// Lấy chi tiết khóa học theo ID
func (r *CourseRepository) GetCourseByID(courseID, userID string) (*dtos.CourseDetailResponse, error) {
	var course dtos.CourseDetailResponse

	// Query chi tiết khóa học
	// Title           string                  `json:"title"`
    // Description     string                  `json:"description"`
    // Like            bool                    `json:"like"`
    // Type            string                  `json:"type"`
    // Level           string                  `json:"level"`
    // Teachers        []string                `json:"teachers"`
    // Star            float64                 `json:"star"`
    // TotalLessons    int                     `json:"total_lessons"`
    // Tags            []string                `json:"tags"`
    // Target          []string                `json:"target"`
    // Require         []string                `json:"require"`
    // TotalTime       int                     `json:"total_time"`
    // TotalRegitered  int                     `json:"total_registered"`
    // UpdatedAt       time.Time               `json:"updated_at"`
    // Language        string                  `json:"language"`
    // CourseContent   []CourseContentResponse `json:"course_content"`
    // Rate            []RateResponse          `json:"rate"`
	query := ``

	err := r.db.Raw(query, courseID).Scan(&course).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get course detail: %w", err)
	}

	return &course, nil
}
*/

func (r *CourseRepository) GetCourseByID(courseID, userID string) (*dtos.CourseDetailResponse, error) {
	// Struct để nhận dữ liệu raw từ database
	type CourseDetailRaw struct {
		Title           string          `json:"title"`
		Description     string          `json:"description"`
		Like            bool            `json:"like"`
		Type            string          `json:"type"`
		Level           string          `json:"level"`
		Teachers        json.RawMessage `json:"teachers"`
		Star            float64         `json:"star"`
		TotalLessons    int             `json:"total_lessons"`
		Tags            json.RawMessage `json:"tags"`
		Target          json.RawMessage `json:"target"`
		Require         json.RawMessage `json:"require"`
		TotalTime       int             `json:"total_time"`
		TotalRegistered int             `json:"total_registered"`
		UpdatedAt       time.Time       `json:"updated_at"`
		Language        string          `json:"language"`
	}

	var courseRaw CourseDetailRaw

	// Query chính để lấy thông tin course
	mainQuery := `
		WITH course_stats AS (
		SELECT 
			c.id,
			c.title,
			c.description,
			c.type,
			c.level,
			array_to_json(c.teachers) as teachers,    -- Convert array to JSON
			array_to_json(c.target) as target,       -- Convert array to JSON  
			array_to_json(c.require) as require,     -- Convert array to JSON
			c.language,
			c.update_at,
			CASE WHEN fc.user_id IS NOT NULL THEN true ELSE false END as like,
			COALESCE(AVG(ur.star), 0) as star,
			COUNT(DISTINCT uc.user_id) as total_registered,
			COUNT(DISTINCT l.id) as total_lessons,
			COALESCE(SUM(
				CASE 
					WHEN l.lesson_type = 'Quiz' THEN (l.quiz_content->>'duration')::int
					WHEN l.lesson_type = 'Online' THEN (l.online_content->>'duration')::int  
					WHEN l.lesson_type = 'Offline' THEN (l.offline_content->>'duration')::int
					ELSE 0
				END
			), 0) as total_time
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

	err := r.db.Raw(mainQuery, courseID, userID).Scan(&courseRaw).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get course detail: %w", err)
	}

	// Parse JSON fields
	var course dtos.CourseDetailResponse
	course.Title = courseRaw.Title
	course.Description = courseRaw.Description
	course.Like = courseRaw.Like
	course.Type = courseRaw.Type
	course.Level = courseRaw.Level
	course.Star = courseRaw.Star
	course.TotalLessons = courseRaw.TotalLessons
	course.TotalTime = courseRaw.TotalTime
	course.TotalRegitered = courseRaw.TotalRegistered
	course.UpdatedAt = courseRaw.UpdatedAt
	course.Language = courseRaw.Language

	// Unmarshal JSON arrays
	if len(courseRaw.Teachers) > 0 {
		if err := json.Unmarshal(courseRaw.Teachers, &course.Teachers); err != nil {
			return nil, fmt.Errorf("failed to parse teachers: %w", err)
		}
	}

	if len(courseRaw.Tags) > 0 {
		if err := json.Unmarshal(courseRaw.Tags, &course.Tags); err != nil {
			return nil, fmt.Errorf("failed to parse tags: %w", err)
		}
	}

	if len(courseRaw.Target) > 0 {
		if err := json.Unmarshal(courseRaw.Target, &course.Target); err != nil {
			return nil, fmt.Errorf("failed to parse target: %w", err)
		}
	}

	if len(courseRaw.Require) > 0 {
		if err := json.Unmarshal(courseRaw.Require, &course.Require); err != nil {
			return nil, fmt.Errorf("failed to parse require: %w", err)
		}
	}

	// Lấy course content (modules và lessons)
	courseContent, err := r.getCourseContent(courseID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get course content: %w", err)
	}
	course.CourseContent = courseContent

	// Lấy ratings
	ratings, err := r.getCourseRatings(courseID)
	if err != nil {
		return nil, fmt.Errorf("failed to get course ratings: %w", err)
	}
	course.Rate = ratings

	return &course, nil
}

// Helper function để lấy course content
func (r *CourseRepository) getCourseContent(courseID, userID string) ([]dtos.CourseContentResponse, error) {
	type ModuleRaw struct {
		ID          string `json:"id"`
		ModuleTitle string `json:"module_title"`
		TotalTime   int    `json:"total_time"`
	}

	var modules []ModuleRaw

moduleQuery := `
    SELECT 
        m.id,
        m.module_title,
        COALESCE(SUM(
            CASE 
                WHEN l.lesson_type = 'Quiz' THEN (l.quiz_content->>'duration')::int
                WHEN l.lesson_type = 'Online' THEN (l.online_content->>'duration')::int
                WHEN l.lesson_type = 'Offline' THEN (l.offline_content->>'duration')::int
                ELSE 0
            END
        ), 0) as total_time
    FROM modules m
    LEFT JOIN lessons l ON m.id = l.module_id
    WHERE m.course_id = $1::uuid  -- ← Cần courseID, không phải moduleID
    GROUP BY m.id, m.module_title
    ORDER BY m.create_at ASC`

	err := r.db.Raw(moduleQuery, courseID).Scan(&modules).Error  // ← Chỉ cần courseID

	// err := r.db.Raw(moduleQuery, modulesID, userID).Scan(&modules).Error
	if err != nil {
		return nil, err
	}

	var courseContent []dtos.CourseContentResponse

	for _, module := range modules {
		// Lấy lessons cho mỗi module
		lessons, err := r.getModuleLessons(module.ID, userID)
		if err != nil {
			return nil, err
		}

		courseContent = append(courseContent, dtos.CourseContentResponse{
			ID:          module.ID,
			ModuleTitle: module.ModuleTitle,
			TotalTime:   module.TotalTime,
			Lessons:     lessons,
		})
	}

	return courseContent, nil
}

// Helper function để lấy lessons của module
func (r *CourseRepository) getModuleLessons(moduleID, userID string) ([]dtos.LessonResponse, error) {
	type LessonRaw struct {
		ID        string `json:"id"`
		Title     string `json:"title"`
		Learned   bool   `json:"learned"`
		StudyTime int    `json:"study_time"`
	}

	var lessons []LessonRaw

	// Query để lấy lessons và check attendance
	lessonQuery := `
		SELECT 
		l.id,
		CASE 
			WHEN l.lesson_type = 'Quiz' THEN 'Quiz Lesson'  -- Hardcode vì Quiz content có questions array 
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

// Helper function để lấy ratings
func (r *CourseRepository) getCourseRatings(courseID string) ([]dtos.RateResponse, error) {
	type RatingRaw struct {
		Name      string    `json:"name"`
		Comment   string    `json:"comment"`
		Star      float64   `json:"star"`
		CreatedAt time.Time `json:"created_at"`
	}

	var ratings []RatingRaw

	// Query để lấy ratings với thông tin user
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
