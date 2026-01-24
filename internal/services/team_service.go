package services

import (
	"context"
	"sfit-platform-web-backend/internal/model"
	"sfit-platform-web-backend/internal/repositories"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type TeamService struct {
	redisClient    *redis.Client
	ctx            context.Context
	teamRepo       *repositories.TeamRepository
	teamMembersSer *TeamMembersService
}

func NewTeamService(teamRepo *repositories.TeamRepository, teamMembersSer *TeamMembersService, redisClient *redis.Client, ctx context.Context) *TeamService {
	return &TeamService{
		redisClient:    redisClient,
		ctx:            ctx,
		teamRepo:       teamRepo,
		teamMembersSer: teamMembersSer,
	}
}

func (s *TeamService) GetTeamList() ([]model.Teams, error) {
	return s.teamRepo.FindAll()
}

func (s *TeamService) CreateTeam(name string, description string) (*model.Teams, error) {

	team := model.Teams{
		ID:          uuid.New(),
		Name:        name,
		Description: description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	return s.teamRepo.Create(team)
}

func (s *TeamService) UpdateTeam(id uuid.UUID, name string, description string) (*model.Teams, error) {
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
