package repositories

import (
	"gorm.io/gorm"
	"sfit-platform-web-backend/entities"
)

type TeamRepository struct {
	db *gorm.DB
}

func NewTeamRepository(db *gorm.DB) *TeamRepository {
	return &TeamRepository{db: db}
}
func (r *TeamRepository) Create(team entities.Teams) (*entities.Teams, error) {
	err := r.db.Create(&team).Error
	return &team, err
}
