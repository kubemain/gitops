// internal/repository/permission_repo.go
package repository

import (
	"aiops-user-service/internal/domain/entity"
	"aiops-user-service/pkg/database"

	"gorm.io/gorm"
)

type PermissionRepository struct {
	db *gorm.DB
}

func NewPermissionRepository() *PermissionRepository {
	return &PermissionRepository{db: database.DB}
}

// FindAll 查询所有权限
func (r *PermissionRepository) FindAll() ([]*entity.Permission, error) {
	var permissions []*entity.Permission
	err := r.db.Order("sort_order ASC").Find(&permissions).Error
	return permissions, err
}

// FindByUserID 查询用户的所有权限
func (r *PermissionRepository) FindByUserID(userID uint) ([]*entity.Permission, error) {
	var permissions []*entity.Permission
	err := r.db.Raw(`
		SELECT DISTINCT p.*
		FROM permissions p
		INNER JOIN role_permissions rp ON p.id = rp.permission_id
		INNER JOIN user_roles ur ON rp.role_id = ur.role_id
		WHERE ur.user_id = ?
		ORDER BY p.sort_order ASC
	`, userID).Scan(&permissions).Error
	return permissions, err
}

// FindByRoleID 查询角色的所有权限
func (r *PermissionRepository) FindByRoleID(roleID uint) ([]*entity.Permission, error) {
	var permissions []*entity.Permission
	err := r.db.Raw(`
		SELECT p.*
		FROM permissions p
		INNER JOIN role_permissions rp ON p.id = rp.permission_id
		WHERE rp.role_id = ?
		ORDER BY p.sort_order ASC
	`, roleID).Scan(&permissions).Error
	return permissions, err
}
