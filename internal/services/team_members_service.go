package services

import (
	"context"
	"encoding/json"
	"errors"
	"sfit-platform-web-backend/internal/dtos"
	"sfit-platform-web-backend/internal/model"
	"sfit-platform-web-backend/internal/repositories"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type TeamMembersService struct {
	redis    *redis.Client
	ctx      context.Context
	repo     *repositories.TeamMembersRepository
	userRepo *repositories.UserRepository
	role     *RoleService
}

func NewTeamMembersService(repo *repositories.TeamMembersRepository, userRepo *repositories.UserRepository, role *RoleService, redisClient *redis.Client, ctx context.Context) *TeamMembersService {
	return &TeamMembersService{
		repo:     repo,
		userRepo: userRepo,
		role:     role,
		redis:    redisClient,
		ctx:      ctx,
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

	err = s.deleteTeamOfMembersCache(userIDStr)
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

	if roleStr == string(model.RoleEnumHead) || roleStr == string(model.RoleEnumVice) || roleStr == string(model.RoleEnumMember) {
		err = s.repo.SaveMember(userID, teamID, roleStr)
		if err != nil {
			return err
		}
	} else {
		return errors.New("invalid role")
	}
	err = s.deleteTeamOfMembersCache(userIDStr)
	if err != nil {
		return err
	}

	return s.role.SyncRoles(userIDStr)
}

func (s *TeamMembersService) GetTeamsJoinedByUser(userIDStr string) ([]dtos.UserJoinedTeamResponse, error) {
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, errors.New("invalid user_id format")
	}

	// Try to get from cache
	cachedMemberIDs, err := s.getTeamOfMembersFromCache(userID)
	if err == nil && len(cachedMemberIDs) > 0 {
		return cachedMemberIDs, nil
	}

	// If not found in cache, get from DB
	memberIDs, err := s.repo.FindTeamsByUserID(userID)
	if err != nil {
		return nil, err
	}

	// Set to cache
	err = s.setTeamOfMembersToCache(userIDStr, memberIDs)
	if err != nil {
		return nil, err
	}

	return memberIDs, nil
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

// =================== Cached Team Methods ===================
func (s *TeamMembersService) buildCacheKey(userID string) string {
	return "user:teams:" + userID
}

func (s *TeamMembersService) getTeamOfMembersFromCache(userID uuid.UUID) ([]dtos.UserJoinedTeamResponse, error) {
	key := s.buildCacheKey(userID.String())
	data, err := s.redis.Get(s.ctx, key).Result()
	if err != nil {
		return nil, err
	}

	var memberIDs []dtos.UserJoinedTeamResponse
	err = json.Unmarshal([]byte(data), &memberIDs)
	if err != nil {
		return nil, err
	}
	return memberIDs, nil
}

func (s *TeamMembersService) setTeamOfMembersToCache(userID string, memberIDs []dtos.UserJoinedTeamResponse) error {
	key := s.buildCacheKey(userID)
	data, err := json.Marshal(memberIDs)
	if err != nil {
		return err
	}
	return s.redis.Set(s.ctx, key, data, 24*time.Hour).Err()
}

func (s *TeamMembersService) deleteTeamOfMembersCache(userID string) error {
	key := s.buildCacheKey(userID)

	return s.redis.Del(s.ctx, key).Err()
}
