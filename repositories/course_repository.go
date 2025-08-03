package repositories

import (
	// "encoding/json"
	"encoding/json"
	"fmt"
	"sfit-platform-web-backend/dtos"
	"sfit-platform-web-backend/entities"
	"strings"

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


// func (r *CourseRepository) GetCoursesInfo(userID uuid.UUID) ([]dtos.CourseInformationResponse, error) {
// 	var courses []dtos.CourseInformationResponse
	
// 	// Optimized query with proper indexing hints and cleaner structure
// 	query := `
// 		WITH course_lessons AS (
// 			SELECT 
// 				m.course_id, 
// 				COUNT(l.id) AS lesson_count
// 			FROM modules m
// 			INNER JOIN lessons l ON l.module_id = m.id
// 			GROUP BY m.course_id
// 		),
// 		user_progress AS (
// 			SELECT 
// 				m.course_id,
// 				COUNT(CASE WHEN la.status = 'present' THEN 1 END) AS learned,
// 				COALESCE(SUM(EXTRACT(EPOCH FROM la.timestamp)::int / 60), 0) AS time_learn
// 			FROM lesson_attendances la
// 			INNER JOIN lessons l ON l.id = la.lesson_id
// 			INNER JOIN modules m ON m.id = l.module_id
// 			WHERE la.user_id = $1
// 			GROUP BY m.course_id
// 		),
// 		course_ratings AS (
// 			SELECT 
// 				courses_id AS course_id, 
// 				ROUND(AVG(star)::numeric, 2) AS rate
// 			FROM user_rates
// 			GROUP BY courses_id
// 		),
// 		course_tags AS (
// 			SELECT 
// 				tt.course_id, 
// 				ARRAY_AGG(t.id ORDER BY t.id) AS tag_list
// 			FROM tag_temps tt
// 			INNER JOIN tags t ON t.id = tt.tag_id
// 			GROUP BY tt.course_id
// 		)
// 		SELECT
// 			c.id,
// 			c.title,
// 			c.description,
// 			c.type,
// 			c.teachers,
// 			COALESCE(cl.lesson_count, 0) AS number_lessons,
// 			COALESCE(up.time_learn, 0) AS time_learn,
// 			COALESCE(cr.rate, 0.0) AS rate,
// 			COALESCE(ct.tag_list, '{}') AS tags,
// 			COALESCE(up.learned, 0) AS learned_lessons,
// 			CASE WHEN up.learned > 0 THEN true ELSE false END AS registed
// 		FROM courses c
// 		LEFT JOIN course_lessons cl ON cl.course_id = c.id
// 		LEFT JOIN user_progress up ON up.course_id = c.id
// 		LEFT JOIN course_ratings cr ON cr.course_id = c.id
// 		LEFT JOIN course_tags ct ON ct.course_id = c.id
// 		ORDER BY c.id`

// 	err := r.db.Raw(query, userID).Scan(&courses).Error
// 	return courses, err
// }

// func (r *CourseRepository) GetCourses(filter dtos.CourseFilter, pageSize, offset int) ([]dtos.CourseInformationResponse, int, error) {
// 	var courses []dtos.CourseInformationResponse
// 	var total int

// 	query := `
// 		WITH course_lessons AS (
// 			SELECT 
// 				m.course_id, 
// 				COUNT(l.id) AS lesson_count
// 			FROM modules m
// 			INNER JOIN lessons l ON l.module_id = m.id
// 			GROUP BY m.course_id
// 		),
// 		user_progress AS (
// 			SELECT 
// 				m.course_id,
// 				COUNT(CASE WHEN la.status = 'present' THEN 1 END) AS learned,
// 				COALESCE(SUM(EXTRACT(EPOCH FROM la.timestamp)::int / 60), 0) AS time_learn
// 			FROM lesson_attendances la
// 			INNER JOIN lessons l ON l.id = la.lesson_id
// 			INNER JOIN modules m ON m.id = l.module_id
// 			WHERE la.user_id = $1
// 			GROUP BY m.course_id
// 		),
// 		course_ratings AS (
// 			SELECT 
// 				courses_id AS course_id, 
// 				ROUND(AVG(star)::numeric, 2) AS rate
// 			FROM user_rates
// 			GROUP BY courses_id
// 		),
// 		course_tags AS (
// 			SELECT 
// 				tt.course_id, 
// 				ARRAY_AGG(t.id ORDER BY t.id) AS tag_list
// 			FROM tag_temps tt
// 			INNER JOIN tags t ON t.id = tt.tag_id
// 			GROUP BY tt.course_id
// 		)
// 		SELECT
// 			c.id,
// 			c.title,
// 			c.description,
// 			c.type,
// 			c.teachers,
// 			COALESCE(cl.lesson_count, 0) AS number_lessons,
// 			COALESCE(up.time_learn, 0) AS time_learn,
// 			COALESCE(cr.rate, 0.0) AS rate,
// 			COALESCE(ct.tag_list, '{}') AS tags,
// 			COALESCE(up.learned, 0) AS learned_lessons,
// 			CASE WHEN up.learned > 0 THEN true ELSE false END AS registed
// 		FROM courses c
// 		LEFT JOIN course_lessons cl ON cl.course_id = c.id
// 		LEFT JOIN user_progress up ON up.course_id = c.id
// 		LEFT JOIN course_ratings cr ON cr.course_id = c.id
// 		LEFT JOIN course_tags ct ON ct.course_id = c.id`

// 	if filter.Title != "" {
// 		query += " WHERE c.title ILIKE '%' || ? || '%'"
// 	}
// 	if filter.OnlyRegisted {
// 		if filter.UserID != uuid.Nil {
// 			query += " AND EXISTS (SELECT 1 FROM lesson_attendances la WHERE la.user_id = ? AND la.course_id = c.id)"
// 		} else {
// 			query += " AND false" // Không có userID, không thể lọc
// 		}
// 	}
// 	if filter.CourseType != "" {
// 		query += " AND c.type = ?"
// 	}
// 	if filter.Level != "" {
// 		query += " AND c.level = ?"
// 	}
// 	query += " ORDER BY c.id LIMIT ? OFFSET ?"
// 	args := []interface{}{filter.UserID}
// 	if filter.Title != "" {
// 		args = append(args, filter.Title)
// 	}
// 	if filter.OnlyRegisted {
// 		args = append(args, filter.UserID)
// 	}
// 	if filter.CourseType != "" {
// 		args = append(args, filter.CourseType)
// 	}
// 	if filter.Level != "" {
// 		args = append(args, filter.Level)
// 	}
// 	args = append(args, pageSize, offset)
// 	err := r.db.Raw(query, args...).Scan(&courses).Error
// 	if err != nil {
// 		return nil, 0, err
// 	}
// 	// Lấy tổng số khóa học
// 	countQuery := "SELECT COUNT(*) FROM courses c"
// 	if filter.Title != "" {
// 		countQuery += " WHERE c.title ILIKE '%' || ? || '%'"
// 	}
// 	if filter.OnlyRegisted {
// 		if filter.UserID != uuid.Nil {
// 			countQuery += " AND EXISTS (SELECT 1 FROM lesson_attendances la WHERE la.user_id = ? AND la.course_id = c.id)"
// 		} else {
// 			countQuery += " AND false" // Không có userID, không thể lọc
// 		}
// 	}
// 	if filter.CourseType != "" {
// 		countQuery += " AND c.type = ?"
// 	}
// 	if filter.Level != "" {
// 		countQuery += " AND c.level = ?"
// 	}
// 	err = r.db.Raw(countQuery, args[:len(args)-2]...).Scan(&total).Error
// 	if err != nil {
// 		return nil, 0, err
// 	}
// 	return courses, total, nil
// }

func (r *CourseRepository) GetCourses(filter dtos.CourseFilter, pageSize, offset int) ([]dtos.CourseInformationResponse, int, error) {
	// var courses []dtos.CourseInformationResponse
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
	var courses []dtos.CourseInformationResponse
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

		courses = append(courses, dtos.CourseInformationResponse{
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