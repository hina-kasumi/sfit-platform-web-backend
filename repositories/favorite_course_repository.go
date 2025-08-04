package repositories

import (
	"gorm.io/gorm"
)

type FavoriteCourseRepository struct {
	db *gorm.DB
}

func NewTeamMemberRepository(db *gorm.DB) *FavoriteCourseRepository {
	return &FavoriteCourseRepository{
		db: db,
	}
}