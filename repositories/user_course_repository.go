package repositories

import (
	"gorm.io/gorm"
)

type UserCourseRepository struct {
	db *gorm.DB
}

func NewUserCourseRepository(db *gorm.DB) *UserCourseRepository {
	return &UserCourseRepository{
		db: db,
	}
}