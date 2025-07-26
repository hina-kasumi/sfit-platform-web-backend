package userservice

import (
	"sfit-platform-web-backend/entities"
	userrepository "sfit-platform-web-backend/repositories/userRepository"
)

func GetUserByID(id string) (*entities.Users, error) {
	return userrepository.GetUserByID(id)
}

func GetUserByusernameOrEmail(username, email string) (*entities.Users, error) {
	return userrepository.GetUserByusernameOrEmail(username, email)
}

func CreateUser(username, email, password string) (*entities.Users, error) {
	return userrepository.CreateUser(username, email, password)
}

func UpdateUser(user *entities.Users) (*entities.Users, error) {
	return userrepository.UpdateUser(user)
}

func DeleteUser(id string) error {
	return userrepository.DeleteUser(id)
}
