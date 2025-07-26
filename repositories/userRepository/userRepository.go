package userrepository

import (
	"errors"
	"sfit-platform-web-backend/entities"
	"sfit-platform-web-backend/infrastructures"

	"github.com/google/uuid"
)

func GetUserByID(id string) (*entities.Users, error) {
	userID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid user ID format")
	}

	user := entities.Users{ID: userID}
	result := infrastructures.GetDB().First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func GetUserByusernameOrEmail(username, email string) (*entities.Users, error) {
	var user *entities.Users

	result := infrastructures.GetDB().Where("username = ? OR email = ?", username, email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}

func CreateUser(username, email, password string) (*entities.Users, error) {
	user := entities.NewUser(username, email, password)
	result := infrastructures.GetDB().Create(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}

func UpdateUser(user *entities.Users) (*entities.Users, error) {
	result := infrastructures.GetDB().Save(user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}

func DeleteUser(id string) error {
	userID, err := uuid.Parse(id)
	if err != nil {
		return errors.New("invalid user ID format")
	}

	user := entities.Users{
		ID: userID,
	}
	result := infrastructures.GetDB().Delete(&user)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("user not found")
	}
	return nil
}
