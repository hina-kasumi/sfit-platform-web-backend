package repositories

import (
	"errors"
	"sfit-platform-web-backend/entities"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (ur *UserRepository) GetUserByID(id string) (*entities.Users, error) {
	userID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid user ID format")
	}

	user := entities.Users{ID: userID}
	result := ur.db.Preload("Roles").First(&user, "id = ?", userID)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (ur *UserRepository) GetUserByusernameOrEmail(username, email string) (*entities.Users, error) {
	var user *entities.Users

	result := ur.db.Preload("Roles").Where("username = ? OR email = ?", username, email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}

func (ur *UserRepository) CreateUser(username, email, password string) (*entities.Users, error) {
	user := entities.NewUser(username, email, password)

	// Gán role mặc định
	roles := entities.UserRole{
		RoleID: entities.RoleEnumUser,
		UserID: user.ID,
	}

	result := ur.db.Create(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	if err := ur.db.Create(&roles).Error; err != nil {
		return nil, err
	}
	user.Roles = []entities.UserRole{roles}
	return user, nil
}

func (ur *UserRepository) UpdateUser(user *entities.Users) (*entities.Users, error) {
	result := ur.db.Save(user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}

func (ur *UserRepository) DeleteUser(id string) error {
	userID, err := uuid.Parse(id)
	if err != nil {
		return errors.New("invalid user ID format")
	}

	user := entities.Users{
		ID: userID,
	}
	result := ur.db.Delete(&user)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("user not found")
	}
	return nil
}

func (ur *UserRepository) GetUserList(page, pageSize int) ([]entities.Users, int64, error) {
	var users []entities.Users
	var total int64

	offset := (page - 1) * pageSize

	result := ur.db.Model(&entities.Users{}).Count(&total)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	result = ur.db.Preload("Roles").Limit(pageSize).Offset(offset).Find(&users)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return users, total, nil
}
