package services

import (
	"golang.org/x/crypto/bcrypt"
	"log"
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

func (user_ser *UserService) ChangePassword(userID, oldPass, newPass string) error {
	user, err := user_ser.user_repo.GetUserByID(userID)
	if err != nil {
		return err
	}
	if err := user.IsValidPasswrod(oldPass); err != nil {
		log.Println("Old password is invalid")
		return err
	}

	hashedNew, _ := bcrypt.GenerateFromPassword([]byte(newPass), bcrypt.DefaultCost)
	user.Password = string(hashedNew)

	_, err = user_ser.user_repo.UpdateUser(user)
	return err
}
