package repositories

import (
	"errors"
	"sfit-platform-web-backend/internal/model"

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

func (ur *UserRepository) FindAdmin() (*model.Users, error) {
	var user model.Users
	err := ur.db.
		Preload("Roles").
		Joins("JOIN user_roles ON user_roles.user_id = users.id").
		Where("user_roles.role_id = ?", model.RoleEnumAdmin).
		First(&user).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, err
}

func (ur *UserRepository) GetUserByID(id string) (*model.Users, error) {
	userID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid user ID format")
	}

	user := model.Users{ID: userID}
	result := ur.db.Preload("Roles").First(&user, "id = ?", userID)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (ur *UserRepository) GetUserByusernameOrEmail(username, email string) (*model.Users, error) {
	var user *model.Users

	result := ur.db.Preload("Roles").Where("username = ? OR email = ?", username, email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}

func (ur *UserRepository) CreateUser(username, email, password string) (*model.Users, error) {
	user := model.NewUser(username, email, password)

	// Gán role mặc định
	roles := model.UserRole{
		RoleID: model.RoleEnumUser,
		UserID: user.ID,
	}

	result := ur.db.Create(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	if err := ur.db.Create(&roles).Error; err != nil {
		return nil, err
	}
	user.Roles = []model.UserRole{roles}
	return user, nil
}

func (ur *UserRepository) UpdateUser(user *model.Users) (*model.Users, error) {
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

	user := model.Users{
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

func (ur *UserRepository) GetUserList(page, pageSize int) ([]model.Users, int64, error) {
	var users []model.Users
	var total int64

	offset := (page - 1) * pageSize

	result := ur.db.Model(&model.Users{}).Count(&total)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	result = ur.db.Preload("Roles").Limit(pageSize).Offset(offset).Find(&users)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return users, total, nil
}
