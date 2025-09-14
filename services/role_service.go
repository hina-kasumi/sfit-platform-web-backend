package services

import (
	"fmt"
	"sfit-platform-web-backend/entities"
	"sfit-platform-web-backend/repositories"
	"slices"
)

type RoleService struct {
	roleRepo *repositories.RoleRepository
	userRepo *repositories.UserRepository
}

func NewRoleService(roleRepo *repositories.RoleRepository, userRepo *repositories.UserRepository) *RoleService {
	return &RoleService{
		roleRepo: roleRepo,
		userRepo: userRepo,
	}
}

func (rs *RoleService) DeleteRole(id string) error {
	return rs.roleRepo.DeleteRole(id)
}

func (rs *RoleService) AddUserRole(userID string, roleIDs ...string) error {
	roleEnums := make([]entities.RoleEnum, len(roleIDs))
	for i, id := range roleIDs {
		roleEnums[i] = entities.RoleEnum(id)
	}
	return rs.roleRepo.AddUserRole(userID, roleEnums...)
}

func (rs *RoleService) RemoveUserRole(curUser, userID string, roleIDs ...string) error {
	user, err := rs.userRepo.GetUserByID(userID)
	if err != nil {
		return err
	}

	if curUser == user.ID.String() {
		if slices.Contains(roleIDs, string(entities.RoleEnumAdmin)) {
			return fmt.Errorf("Cannot remove admin role from yourself")
		}
	}
	return rs.roleRepo.RemoveUserRole(userID, roleIDs...)
}

func (rs *RoleService) SyncRoles(userID string) error {
	return rs.roleRepo.SyncRoles(userID)
}

func (rs *RoleService) GetUserRoles(userID string) ([]entities.RoleEnum, error) {
	return rs.roleRepo.GetUserRoles(userID)
}
