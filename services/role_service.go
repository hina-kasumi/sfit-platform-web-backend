package services

import (
	"sfit-platform-web-backend/entities"
	"sfit-platform-web-backend/repositories"
)

type RoleService struct {
	roleRepo *repositories.RoleRepository
}

func NewRoleService(roleRepo *repositories.RoleRepository) *RoleService {
	return &RoleService{roleRepo: roleRepo}
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

func (rs *RoleService) RemoveUserRole(userID string, roleIDs ...string) error {
	return rs.roleRepo.RemoveUserRole(userID, roleIDs...)
}
