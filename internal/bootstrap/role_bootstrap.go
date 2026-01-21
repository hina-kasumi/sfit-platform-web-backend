package bootstrap

import "sfit-platform-web-backend/internal/services"

func InitRoles(roleService *services.RoleService) error {
	return roleService.EnsureRolesExists()
}
