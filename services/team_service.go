package services

import (
	"github.com/google/uuid"
	"sfit-platform-web-backend/entities"
	"sfit-platform-web-backend/repositories"
	"time"
)

type TeamService struct {
	teamRepo *repositories.TeamRepository
}

func NewTeamService(teamRepo *repositories.TeamRepository) *TeamService {
	return &TeamService{teamRepo: teamRepo}
}

func (s *TeamService) CreateTeam(id string, name string, description string) (*entities.Teams, error) {

	uuidID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	team := entities.Teams{
		ID:          uuidID,
		Name:        name,
		Description: description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	return s.teamRepo.Create(team)
}
