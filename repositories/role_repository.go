package repositories

import (
	"sfit-platform-web-backend/entities"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type RoleRepository struct {
	db *gorm.DB
}

// NewRoleRepository creates a new RoleRepository
func NewRoleRepository(db *gorm.DB) *RoleRepository {
	// Insert all RoleEnum values into the roles table if not exists
	allRoles := []entities.Role{
		{ID: entities.RoleEnumAdmin},
		{ID: entities.RoleEnumUser},
		{ID: entities.RoleEnumHead},
		{ID: entities.RoleEnumMember},
		{ID: entities.RoleEnumVice},
	}

	db.Clauses(clause.OnConflict{DoNothing: true}).Create(&allRoles)
	return &RoleRepository{db: db}
}

func (rr *RoleRepository) CreateRoles(ids ...string) error {
	roles := make([]entities.Role, 0, len(ids))
	for _, id := range ids {
		roles = append(roles, entities.Role{ID: entities.RoleEnum(id)})
	}
	return rr.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}}, // chỉ định cột key
		DoNothing: true,
	}).Create(&roles).Error
}

// DeleteRole deletes a role from the database
func (rr *RoleRepository) DeleteRole(id string) error {
	var role entities.Role
	if err := rr.db.First(&role, "id = ?", id).Error; err != nil {
		return err
	}
	if err := rr.db.Delete(&role).Error; err != nil {
		return err
	}
	var userRole entities.UserRole
	if err := rr.db.Where("role_id = ?", id).First(&userRole).Error; err != nil {
		return err
	}
	if err := rr.db.Delete(&userRole).Error; err != nil {
		return err
	}
	return nil
}

// GetUserRoles retrieves the roles for a specific user
func (rr *RoleRepository) GetUserRoles(userID string) ([]entities.RoleEnum, error) {
	var userRoles []entities.UserRole
	err := rr.db.Where("user_id = ?", userID).Find(&userRoles).Error
	if err != nil {
		return nil, err
	}

	roleIDs := make([]entities.RoleEnum, 0, len(userRoles))
	for _, ur := range userRoles {
		roleIDs = append(roleIDs, ur.RoleID)
	}
	return roleIDs, nil
}

// AddUserRole adds roles to a specific user
func (rr *RoleRepository) AddUserRole(userID string, roleIDs ...entities.RoleEnum) error {
	uuidUserID, err := uuid.Parse(userID)
	if err != nil {
		return err
	}
	userRoles := make([]entities.UserRole, 0, len(roleIDs))
	for _, roleID := range roleIDs {
		userRoles = append(userRoles, entities.UserRole{UserID: uuidUserID, RoleID: roleID})
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
	return rr.db.Where("user_id = ? AND role_id IN ?", uuidUserID, roleIDs).Delete(&entities.UserRole{}).Error
}

func (rr *RoleRepository) SyncRoles(userID string) error {
	var members []entities.TeamMembers
	err := rr.db.Where("user_id = ?", userID).Find(&members).Error
	if err != nil {
		return err
	}
	roleSet := make(map[entities.RoleEnum]bool)
	for _, m := range members {
		roleSet[m.Role] = true
	}

	uniqueRoles := make([]entities.RoleEnum, 0, len(roleSet))
	for role := range roleSet {
		uniqueRoles = append(uniqueRoles, entities.RoleEnum(role))
	}

	rr.AddUserRole(userID, uniqueRoles...)

	roleEnums, err := rr.GetUserRoles(userID)
	if err != nil {
		return err
	}

	listRoles := make([]string, 0, len(roleEnums))
	for _, v := range roleEnums {
		switch v {
		case entities.RoleEnumHead:
			if _, exists := roleSet[entities.RoleEnumHead]; !exists {
				listRoles = append(listRoles, string(entities.RoleEnumHead))
			}
		case entities.RoleEnumVice:
			if _, exists := roleSet[entities.RoleEnumVice]; !exists {
				listRoles = append(listRoles, string(entities.RoleEnumVice))
			}
		case entities.RoleEnumMember:
			if _, exists := roleSet[entities.RoleEnumMember]; !exists {
				listRoles = append(listRoles, string(entities.RoleEnumMember))
			}
		}
	}
	return rr.RemoveUserRole(userID, listRoles...)
}
