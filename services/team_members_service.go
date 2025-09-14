package services

import (
	"errors"
	"sfit-platform-web-backend/dtos"
	"sfit-platform-web-backend/entities"
	"sfit-platform-web-backend/repositories"

	"github.com/google/uuid"
)

type TeamMembersService struct {
	repo     *repositories.TeamMembersRepository
	userRepo *repositories.UserRepository
	role     *RoleService
}

func NewTeamMembersService(repo *repositories.TeamMembersRepository, userRepo *repositories.UserRepository, role *RoleService) *TeamMembersService {
	return &TeamMembersService{
		repo:     repo,
		userRepo: userRepo,
		role:     role,
	}
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
	return s.role.SyncRoles(userIDStr)
}

func (s *TeamMembersService) SaveMember(userIDStr, teamIDStr, roleStr string) error {
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return errors.New("invalid user_id format")
	}

	teamID, err := uuid.Parse(teamIDStr)
	if err != nil {
		return errors.New("invalid team_id format")
	}

	if roleStr == string(entities.RoleEnumHead) || roleStr == string(entities.RoleEnumVice) || roleStr == string(entities.RoleEnumMember) {
		err = s.repo.SaveMember(userID, teamID, roleStr)
		if err != nil {
			return err
		}
	} else {
		return errors.New("invalid role")
	}

	return s.role.SyncRoles(userIDStr)
}

func (s *TeamMembersService) GetTeamsJoinedByUser(userIDStr string) ([]dtos.UserJoinedTeamResponse, error) {
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, errors.New("invalid user_id format")
	}
	return s.repo.FindTeamsByUserID(userID)
}

func (s *TeamMembersService) GetMembers(teamID string, page, pageSize int, role string) (*dtos.PageListResp, error) {
	if page < 1 {
		page = 1
	}
	members, total, err := s.repo.FindMembersByTeamID(teamID, page, pageSize, role)
	if err != nil {
		return nil, err
	}

	return &dtos.PageListResp{
		Items:      members,
		Page:       page,
		PageSize:   pageSize,
		TotalCount: total,
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

func (s *TeamMembersService) GetRoleUserInTeam(userIDStr, teamIDStr string) (string, error) {
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return "", errors.New("invalid user_id format")
	}
	teamID, err := uuid.Parse(teamIDStr)
	if err != nil {
		return "", errors.New("invalid team_id format")
	}

	role, err := s.repo.FindRoleByUserIDAndTeamID(userID, teamID)
	if err != nil {
		return "", err
	}
	return role, nil
}
