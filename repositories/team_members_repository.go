package repositories

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"sfit-platform-web-backend/dtos"
	entities "sfit-platform-web-backend/entities"
	"time"
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

func (r *TeamMembersRepository) DeleteByUserIDAndTeamID(userID, teamID uuid.UUID) error {
	result := r.db.Where("user_id = ? AND team_id = ?", userID, teamID).Delete(&entities.TeamMembers{})
	return result.Error
}

func (r *TeamMembersRepository) UpdateRole(userID, teamID uuid.UUID, role string) error {
	updateData := map[string]interface{}{
		"role":       role,
		"updated_at": time.Now(),
	}
	result := r.db.Model(&entities.TeamMembers{}).
		Where("user_id = ? AND team_id = ?", userID, teamID).
		Updates(updateData)
	return result.Error
}

func (r *TeamMembersRepository) FindTeamsByUserID(userID uuid.UUID) ([]dtos.UserJoinedTeamResponse, error) {
	var results []dtos.UserJoinedTeamResponse

	err := r.db.Table("team_members tm").
		Select("tm.team_id, t.name, tm.role").
		Joins("JOIN teams t ON tm.team_id = t.id").
		Where("tm.user_id = ?", userID).
		Scan(&results).Error

	if err != nil {
		return nil, err
	}
	return results, nil
}
