package services

import (
	"errors"
	"github.com/google/uuid"
	"sfit-platform-web-backend/entities"
	"sfit-platform-web-backend/repositories"
	"time"
)

type TeamMembersService struct {
	repo     *repositories.TeamMembersRepository
	userRepo *repositories.UserRepository
}

func NewTeamMembersService(repo *repositories.TeamMembersRepository, userRepo *repositories.UserRepository) *TeamMembersService {
	return &TeamMembersService{
		repo:     repo,
		userRepo: userRepo,
	}
}

func (s *TeamMembersService) AddMember(userIDStr, teamIDStr, roleStr string) (*entities.TeamMembers, error) {
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, errors.New("invalid user_id format")
	}
	teamID, err := uuid.Parse(teamIDStr)
	if err != nil {
		return nil, errors.New("invalid team id format")
	}

	switch roleStr {
	case string(entities.RoleHeader),
		string(entities.RoleViceHeader),
		string(entities.RoleMember):
	default:
		return nil, errors.New("invalid role format")
	}

	exitsing, err := s.repo.FindByUserIDAndTeamID(userID, teamID)
	if err == nil && exitsing != nil {
		return nil, errors.New("user already in the team")
	}
	existsUser, err := s.userRepo.GetUserByID(userIDStr)
	if err != nil || existsUser == nil {
		return nil, errors.New("user_id does not exist")
	}
	teamMember := entities.TeamMembers{
		UserID:    userID,
		TeamID:    teamID,
		Role:      entities.TeamRole(roleStr),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return s.repo.Create(&teamMember)
}
