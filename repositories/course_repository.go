package repositories

import (
	"sfit-platform-web-backend/entities"
	"strings"

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
// func (r *CourseRepository) Create(course *entities.Course) error {
//     return r.db.Create(course).Error
// }
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

