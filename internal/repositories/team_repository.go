package repositories

import (
	"sfit-platform-web-backend/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TeamRepository struct {
	db *gorm.DB
}

func NewTeamRepository(db *gorm.DB) *TeamRepository {
	return &TeamRepository{db: db}
}

func (r *TeamRepository) FindAll() ([]model.Teams, error) {
	var teams []model.Teams
	result := r.db.Find(&teams)
	if result.Error != nil {
		return nil, result.Error
	}
	return teams, nil
}

func (r *TeamRepository) Create(team model.Teams) (*model.Teams, error) {
	err := r.db.Create(&team).Error
	return &team, err
}

func (r *TeamRepository) Update(team *model.Teams) error {
	result := r.db.Save(team)
	return result.Error
}

func (r *TeamRepository) FindByID(id uuid.UUID) (*model.Teams, error) {
	var team model.Teams
	result := r.db.First(&team, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &team, nil
}

func (r *TeamRepository) DeleteTeam(id uuid.UUID) error {
	team := model.Teams{ID: id}
	result := r.db.First(&team)
	if result.Error != nil {
		return result.Error
	}

	return r.db.Delete(&team).Error
}
