package repositories

import (
	"github.com/google/uuid"
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

func (r *TeamRepository) Update(team *entities.Teams) error {
	result := r.db.Save(team)
	return result.Error
}

func (r *TeamRepository) FindByID(id uuid.UUID) (*entities.Teams, error) {
	var team entities.Teams
	result := r.db.First(&team, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &team, nil
}
