package services

import (
	"sfit-platform-web-backend/entities"
	"sfit-platform-web-backend/repositories"
	"time"

	"github.com/google/uuid"
)

type TeamService struct {
	teamRepo       *repositories.TeamRepository
	teamMembersSer *TeamMembersService
}

func NewTeamService(teamRepo *repositories.TeamRepository, teamMembersSer *TeamMembersService) *TeamService {
	return &TeamService{
		teamRepo:       teamRepo,
		teamMembersSer: teamMembersSer,
	}
}

func (s *TeamService) GetTeamList() ([]entities.Teams, error) {
	return s.teamRepo.FindAll()
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

func (s *TeamService) DeleteTeam(id string) error {
	teamID, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	err = s.teamRepo.DeleteTeam(teamID)
	if err != nil {
		return err
	}
	err = s.teamMembersSer.DeleteAllMemberInTeam(id)
	if err != nil {
		return err
	}

	return nil
}
