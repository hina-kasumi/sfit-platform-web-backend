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

func (s *TeamService) CreateTeam(name string, description string) (*entities.Teams, error) {

	team := entities.Teams{
		ID:          uuid.New(),
		Name:        name,
		Description: description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	return s.teamRepo.Create(team)
}

func (s *TeamService) UpdateTeam(id uuid.UUID, name string, description string) (*entities.Teams, error) {
	team, err := s.teamRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	team.Name = name
	team.Description = description
	team.UpdatedAt = time.Now()

	err = s.teamRepo.Update(team)
	if err != nil {
		return nil, err
	}

	return team, nil
}
