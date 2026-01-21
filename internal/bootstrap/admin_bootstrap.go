package bootstrap

import (
	"sfit-platform-web-backend/internal/config"
	"sfit-platform-web-backend/internal/services"
)

func InitAdmin(cfg config.AdminConfig, userSer *services.UserService) error {
	return userSer.EnsureAdminExists(cfg)
}
