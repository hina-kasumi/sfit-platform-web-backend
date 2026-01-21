package repositories

import (
	"sfit-platform-web-backend/internal/config"
	"sfit-platform-web-backend/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type RoleRepository struct {
	db *gorm.DB
}

// NewRoleRepository creates a new RoleRepository
func NewRoleRepository(cfg config.AdminConfig, db *gorm.DB) *RoleRepository {
	return &RoleRepository{db: db}
}

func (rr *RoleRepository) CreateRoles(ids ...string) error {
	roles := make([]model.Role, 0, len(ids))
	for _, id := range ids {
		roles = append(roles, model.Role{ID: model.RoleEnum(id)})
	}
	return rr.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}}, // chỉ định cột key
		DoNothing: true,
	}).Create(&roles).Error
}

func (rr *RoleRepository) EnsureRolesExists() error {
	// Insert all RoleEnum values into the roles table if not exists
	allRoles := []model.Role{
		{ID: model.RoleEnumAdmin},
		{ID: model.RoleEnumUser},
		{ID: model.RoleEnumHead},
		{ID: model.RoleEnumMember},
		{ID: model.RoleEnumVice},
	}

	return rr.db.Clauses(clause.OnConflict{DoNothing: true}).Create(&allRoles).Error
}

// DeleteRole deletes a role from the database
func (rr *RoleRepository) DeleteRole(id string) error {
	var role model.Role
	if err := rr.db.First(&role, "id = ?", id).Error; err != nil {
		return err
	}
	if err := rr.db.Delete(&role).Error; err != nil {
		return err
	}
	var userRole model.UserRole
	if err := rr.db.Where("role_id = ?", id).First(&userRole).Error; err != nil {
		return err
	}
	if err := rr.db.Delete(&userRole).Error; err != nil {
		return err
	}
	return nil
}

// GetUserRoles retrieves the roles for a specific user
func (rr *RoleRepository) GetUserRoles(userID string) ([]model.RoleEnum, error) {
	var userRoles []model.UserRole
	err := rr.db.Where("user_id = ?", userID).Find(&userRoles).Error
	if err != nil {
		return nil, err
	}

	roleIDs := make([]model.RoleEnum, 0, len(userRoles))
	for _, ur := range userRoles {
		roleIDs = append(roleIDs, ur.RoleID)
	}
	return roleIDs, nil
}

// AddUserRole adds roles to a specific user
func (rr *RoleRepository) AddUserRole(userID string, roleIDs ...model.RoleEnum) error {
	uuidUserID, err := uuid.Parse(userID)
	if err != nil {
		return err
	}
	userRoles := make([]model.UserRole, 0, len(roleIDs))
	for _, roleID := range roleIDs {
		userRoles = append(userRoles, model.UserRole{UserID: uuidUserID, RoleID: roleID})
	}
	return rr.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "role_id"}, {Name: "user_id"}}, // chỉ định cột key
		DoNothing: true,
	}).Create(&userRoles).Error
}

// RemoveUserRole removes roles from a specific user
func (rr *RoleRepository) RemoveUserRole(userID string, roleIDs ...string) error {
	uuidUserID, err := uuid.Parse(userID)
	if err != nil {
		return err
	}
	return rr.db.Where("user_id = ? AND role_id IN ?", uuidUserID, roleIDs).Delete(&model.UserRole{}).Error
}

func (rr *RoleRepository) SyncRoles(userID string) error {
	var members []model.TeamMembers
	err := rr.db.Where("user_id = ?", userID).Find(&members).Error
	if err != nil {
		return err
	}
	roleSet := make(map[model.RoleEnum]bool)
	for _, m := range members {
		roleSet[m.Role] = true
	}

	uniqueRoles := make([]model.RoleEnum, 0, len(roleSet))
	for role := range roleSet {
		uniqueRoles = append(uniqueRoles, model.RoleEnum(role))
	}

	rr.AddUserRole(userID, uniqueRoles...)

	roleEnums, err := rr.GetUserRoles(userID)
	if err != nil {
		return err
	}

	listRoles := make([]string, 0, len(roleEnums))
	for _, v := range roleEnums {
		switch v {
		case model.RoleEnumHead:
			if _, exists := roleSet[model.RoleEnumHead]; !exists {
				listRoles = append(listRoles, string(model.RoleEnumHead))
			}
		case model.RoleEnumVice:
			if _, exists := roleSet[model.RoleEnumVice]; !exists {
				listRoles = append(listRoles, string(model.RoleEnumVice))
			}
		case model.RoleEnumMember:
			if _, exists := roleSet[model.RoleEnumMember]; !exists {
				listRoles = append(listRoles, string(model.RoleEnumMember))
			}
		}
	}
	return rr.RemoveUserRole(userID, listRoles...)
}
