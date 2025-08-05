package repositories

import (
	"sfit-platform-web-backend/entities"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FavoriteCourseRepository struct {
	db *gorm.DB
}

func NewFavoriteCourseRepository(db *gorm.DB) *FavoriteCourseRepository {
	return &FavoriteCourseRepository{
		db: db,
	}
}

// MarkCourseAsFavourite marks a course as favourite for a user
func (r *FavoriteCourseRepository) MarkCourseAsFavourite(userID, courseID string) error {
	favorCourse := entities.FavoriteCourse{
		UserID:    uuid.MustParse(userID),
		CourseID:  uuid.MustParse(courseID),
		CreatedAt: time.Now(),
	}

	return r.db.Create(&favorCourse).Error
}

// UnmarkCourseAsFavourite removes a course from the user's favourites
func (r *FavoriteCourseRepository) UnmarkCourseAsFavourite(userID, courseID string) error {
	return r.db.Where("user_id = ? AND course_id = ?", userID, courseID).
		Delete(&entities.FavoriteCourse{}).Error
}