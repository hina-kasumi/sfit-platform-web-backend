package repositories

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	entities "sfit-platform-web-backend/entities"
)

type TeamMembersRepository struct {
	db *gorm.DB
}

func NewTeamMembersRepository(db *gorm.DB) *TeamMembersRepository {
	return &TeamMembersRepository{db: db}
}


func (r *TeamMembersRepository) Create(teamMember *entities.TeamMembers) (*entities.TeamMembers, error) {
	err := r.db.Create(teamMember).Error
	return teamMember, err
}


func (r *TeamMembersRepository) FindByUserIDAndTeamID(userID, teamID uuid.UUID) (*entities.TeamMembers, error) {
	var tm entities.TeamMembers
	result := r.db.Where("user_id = ? AND team_id = ?", userID, teamID).First(&tm)
	if result.Error != nil {
		return nil, result.Error
	}
	return &tm, nil
}