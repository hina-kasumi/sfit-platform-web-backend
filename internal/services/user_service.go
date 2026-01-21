package services

import (
	"sfit-platform-web-backend/internal/config"
	"sfit-platform-web-backend/internal/dtos"
	"sfit-platform-web-backend/internal/model"
	"sfit-platform-web-backend/internal/repositories"
)

type UserService struct {
	user_repo         *repositories.UserRepository
	user_profile_repo *repositories.UserProfileRepository
	role_repo         *RoleService
}

func NewUserService(user_repo *repositories.UserRepository, user_profile_repo *repositories.UserProfileRepository, role_repo *RoleService) *UserService {
	return &UserService{
		user_repo:         user_repo,
		user_profile_repo: user_profile_repo,
		role_repo:         role_repo,
	}
}

func (user_ser *UserService) EnsureAdminExists(cfg config.AdminConfig) error {
	admin, err := user_ser.user_repo.FindAdmin()
	if err != nil {
		return err
	}

	if admin == nil {
		defaultAdminUsername := cfg.Username
		defaultAdminEmail := cfg.Email
		defaultAdminPassword := cfg.Password

		// Create default admin user
		adminUser, err := user_ser.user_repo.CreateUser(defaultAdminUsername, defaultAdminEmail, defaultAdminPassword)
		if err != nil {
			return err
		}

		err = user_ser.role_repo.AddUserRole(adminUser.ID.String(), string(model.RoleEnumAdmin))
		if err != nil {
			return err
		}

		// Initialize admin user profile
		err = user_ser.user_profile_repo.InitAdminProfile(adminUser.ID, defaultAdminEmail)
		if err != nil {
			return err
		}
	}

	return nil
}

func (user_ser *UserService) GetUserByID(id string) (*model.Users, error) {
	return user_ser.user_repo.GetUserByID(id)
}

func (user_ser *UserService) GetUserByusernameOrEmail(username, email string) (*model.Users, error) {
	return user_ser.user_repo.GetUserByusernameOrEmail(username, email)
}

func (user_ser *UserService) CreateUser(username, email, password string) (*model.Users, error) {
	return user_ser.user_repo.CreateUser(username, email, password)
}

func (user_ser *UserService) UpdateUser(user *model.Users) (*model.Users, error) {
	return user_ser.user_repo.UpdateUser(user)
}

func (user_ser *UserService) DeleteUser(id string) error {
	return user_ser.user_repo.DeleteUser(id)
}

func (user_ser *UserService) GetUserList(page, pageSize int) ([]dtos.UserListItem, int, int, int64, error) {
	users, total, err := user_ser.user_repo.GetUserList(page, pageSize)
	if err != nil {
		return nil, 0, 0, 0, err
	}

	var userList []dtos.UserListItem
	for _, user := range users {
		userList = append(userList, dtos.UserListItem{
			ID:       user.ID.String(),
			Username: user.Username,
			Email:    user.Email,
		})
	}

	return userList, page, pageSize, total, nil
}
