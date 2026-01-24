package services

import (
	"context"
	"encoding/json"
	"sfit-platform-web-backend/internal/config"
	"sfit-platform-web-backend/internal/dtos"
	"sfit-platform-web-backend/internal/model"
	"sfit-platform-web-backend/internal/repositories"
	"time"

	"github.com/redis/go-redis/v9"
)

type UserService struct {
	redisClient       *redis.Client
	ctx               context.Context
	user_repo         *repositories.UserRepository
	user_profile_repo *repositories.UserProfileRepository
	role_repo         *RoleService
}

func NewUserService(ctx context.Context, redisClient *redis.Client, user_repo *repositories.UserRepository, user_profile_repo *repositories.UserProfileRepository, role_repo *RoleService) *UserService {
	return &UserService{
		redisClient:       redisClient,
		ctx:               ctx,
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
	cachedUser, err := user_ser.GetUserFromCache(id)
	if err != nil {
		return nil, err
	}
	if cachedUser != nil {
		return cachedUser, nil
	}

	user, err := user_ser.user_repo.GetUserByID(id)
	if err != nil {
		return nil, err
	}

	if user != nil {
		err = user_ser.SetUserToCache(id, user)
		if err != nil {
			return nil, err
		}
	}

	return user, nil
}

func (user_ser *UserService) GetUserByusernameOrEmail(username, email string) (*model.Users, error) {
	return user_ser.user_repo.GetUserByusernameOrEmail(username, email)
}

func (user_ser *UserService) CreateUser(username, email, password string) (*model.Users, error) {
	user, err := user_ser.user_repo.CreateUser(username, email, password)
	if err != nil {
		return nil, err
	}

	if user != nil {
		err = user_ser.SetUserToCache(user.ID.String(), user)
		if err != nil {
			return nil, err
		}
	}

	return user, nil
}

func (user_ser *UserService) UpdateUser(user *model.Users) (*model.Users, error) {
	updatedUser, err := user_ser.user_repo.UpdateUser(user)
	if err != nil {
		return nil, err
	}

	if updatedUser != nil {
		err = user_ser.SetUserToCache(updatedUser.ID.String(), updatedUser)
		if err != nil {
			return nil, err
		}
	}

	return updatedUser, nil
}

func (user_ser *UserService) DeleteUser(id string) error {
	err := user_ser.user_repo.DeleteUser(id)
	if err != nil {
		return err
	}

	err = user_ser.DeleteUserFromCache(id)
	if err != nil {
		return err
	}

	return nil
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

// Cache management methods
func (user_ser *UserService) generateCacheKey(identifier string) string {
	return "user:" + identifier
}

func (user_ser *UserService) GetUserFromCache(identifier string) (*model.Users, error) {
	cacheKey := user_ser.generateCacheKey(identifier)
	cachedUser, err := user_ser.redisClient.Get(user_ser.ctx, cacheKey).Result()
	if err == redis.Nil {
		return nil, nil // Cache miss
	} else if err != nil {
		return nil, err
	}

	// Deserialize cachedUser (assuming JSON serialization)
	var user model.Users
	err = json.Unmarshal([]byte(cachedUser), &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (user_ser *UserService) SetUserToCache(identifier string, user *model.Users) error {
	cacheKey := user_ser.generateCacheKey(identifier)
	// Serialize user (assuming JSON serialization)
	userData, err := json.Marshal(user)
	if err != nil {
		return err
	}

	err = user_ser.redisClient.Set(user_ser.ctx, cacheKey, userData, 10*time.Minute).Err()
	if err != nil {
		return err
	}

	return nil
}

func (user_ser *UserService) DeleteUserFromCache(identifier string) error {
	cacheKey := user_ser.generateCacheKey(identifier)
	err := user_ser.redisClient.Del(user_ser.ctx, cacheKey).Err()
	if err != nil {
		return err
	}

	return nil
}
