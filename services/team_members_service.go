package services

import (
	"errors"
	"github.com/google/uuid"
	"sfit-platform-web-backend/dtos"
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

func (s *TeamMembersService) DeleteMember(userIDStr, teamIDStr string) error {
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return errors.New("invalid user_id format")
	}
	teamID, err := uuid.Parse(teamIDStr)
	if err != nil {
		return errors.New("invalid team_id format")
	}

	existing, err := s.repo.FindByUserIDAndTeamID(userID, teamID)
	if err != nil || existing == nil {
		return errors.New("user is not a member of the team")
	}

	err = s.repo.DeleteByUserIDAndTeamID(userID, teamID)
	if err != nil {
		return err
	}
	return nil
}

func (s *TeamMembersService) UpdateMemberRole(userIDStr, teamIDStr, roleStr string) error {
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return errors.New("invalid user_id format")
	}

	teamID, err := uuid.Parse(teamIDStr)
	if err != nil {
		return errors.New("invalid team_id format")
	}

	validRoles := map[string]bool{
		string(entities.RoleHeader):     true,
		string(entities.RoleViceHeader): true,
		string(entities.RoleMember):     true,
	}

	if !validRoles[roleStr] {
		return errors.New("invalid role value")
	}

	existingMember, err := s.repo.FindByUserIDAndTeamID(userID, teamID)
	if err != nil || existingMember == nil {
		return errors.New("user is not a member of this team")
	}

	err = s.repo.UpdateRole(userID, teamID, roleStr)
	if err != nil {
		return err
	}

	return nil
}

func (s *TeamMembersService) GetTeamsJoinedByUser(userIDStr string) ([]dtos.UserJoinedTeamResponse, error) {
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, errors.New("invalid user_id format")
	}
	return s.repo.FindTeamsByUserID(userID)
}

func (s *TeamMembersService) GetMembers(teamID string, page, pageSize int) (*dtos.TeamMembersResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	members, total, err := s.repo.FindMembersByTeamID(teamID, page, pageSize)
	if err != nil {
		return nil, err
	}

	return &dtos.TeamMembersResponse{
		Users:    members,
		Page:     page,
		PageSize: pageSize,
		Total:    total,
	}, nil
}

func (s *TeamMembersService) DeleteAllMemberInTeam(teamID string) error {
	teamUUID, err := uuid.Parse(teamID)
	if err != nil {
		return errors.New("invalid team_id format")
	}

	err = s.repo.DeleteAllMembersInTeam(teamUUID)
	if err != nil {
		return err
	}

	return nil
}
