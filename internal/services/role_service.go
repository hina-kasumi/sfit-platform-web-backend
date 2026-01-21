package services

import (
	"fmt"
	"sfit-platform-web-backend/internal/model"
	"sfit-platform-web-backend/internal/repositories"
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

func (rs *RoleService) EnsureRolesExists() error {
	return rs.roleRepo.EnsureRolesExists()
}

func (rs *RoleService) DeleteRole(id string) error {
	return rs.roleRepo.DeleteRole(id)
}

func (rs *RoleService) AddUserRole(userID string, roleIDs ...string) error {
	roleEnums := make([]model.RoleEnum, len(roleIDs))
	for i, id := range roleIDs {
		roleEnums[i] = model.RoleEnum(id)
	}
	return rs.roleRepo.AddUserRole(userID, roleEnums...)
}

func (rs *RoleService) RemoveUserRole(curUser, userID string, roleIDs ...string) error {
	user, err := rs.userRepo.GetUserByID(userID)
	if err != nil {
		return err
	}

	if curUser == user.ID.String() {
		if slices.Contains(roleIDs, string(model.RoleEnumAdmin)) {
			return fmt.Errorf("cannot remove admin role from yourself")
		}
	}
	return rs.roleRepo.RemoveUserRole(userID, roleIDs...)
}

func (rs *RoleService) SyncRoles(userID string) error {
	return rs.roleRepo.SyncRoles(userID)
}

func (rs *RoleService) GetUserRoles(userID string) ([]model.RoleEnum, error) {
	return rs.roleRepo.GetUserRoles(userID)
}
