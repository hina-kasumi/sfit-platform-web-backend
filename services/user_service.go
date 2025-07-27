package services

import (
	"sfit-platform-web-backend/entities"
	"sfit-platform-web-backend/repositories"
)

type UserService struct {
	user_repo *repositories.UserRepository
}

func NewUserService(user_repo *repositories.UserRepository) *UserService {
	return &UserService{user_repo: user_repo}
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
