package repositories

import (
	"gorm.io/gorm"
)

type UserRateRepository struct {
	db *gorm.DB
}

func NewUserRateRepository(db *gorm.DB) *UserRateRepository {
	return &UserRateRepository{
		db: db,
	}
}