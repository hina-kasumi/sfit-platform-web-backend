package services

import (
	"sfit-platform-web-backend/dtos"
	"sfit-platform-web-backend/entities"
	"sfit-platform-web-backend/repositories"
)

type UserService struct {
	user_repo *repositories.UserRepository
	role_repo *RoleService
}

func NewUserService(user_repo *repositories.UserRepository, role_repo *RoleService) *UserService {
	return &UserService{user_repo: user_repo, role_repo: role_repo}
}

func (user_ser *UserService) GetUserByID(id string) (*entities.Users, error) {
	return user_ser.user_repo.GetUserByID(id)
}

func (user_ser *UserService) GetUserByusernameOrEmail(username, email string) (*entities.Users, error) {
	return user_ser.user_repo.GetUserByusernameOrEmail(username, email)
}

func (user_ser *UserService) CreateUser(username, email, password string) (*entities.Users, error) {
	return user_ser.user_repo.CreateUser(username, email, password)
}

func (user_ser *UserService) UpdateUser(user *entities.Users) (*entities.Users, error) {
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
