// internal/repository/role_repo.go
package repository

import (
	"aiops-user-service/internal/domain/entity"
	"aiops-user-service/pkg/database"
	"errors"

	"gorm.io/gorm"
)

type RoleRepository struct {
	db *gorm.DB
}

func NewRoleRepository() *RoleRepository {
	return &RoleRepository{db: database.DB}
}

// FindAll 查询所有角色
func (r *RoleRepository) FindAll() ([]*entity.Role, error) {
	var roles []*entity.Role
	err := r.db.Order("sort_order ASC").Find(&roles).Error
	return roles, err
}

// FindByID 根据ID查询角色
func (r *RoleRepository) FindByID(id uint) (*entity.Role, error) {
	var role entity.Role
	err := r.db.Preload("Permissions").First(&role, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("角色不存在")
		}
		return nil, err
	}
	return &role, nil
}

// Create 创建角色
func (r *RoleRepository) Create(role *entity.Role) error {
	return r.db.Create(role).Error
}

// Update 更新角色
func (r *RoleRepository) Update(role *entity.Role) error {
	return r.db.Save(role).Error
}

// Delete 删除角色
func (r *RoleRepository) Delete(id uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 检查是否为系统角色
		var role entity.Role
		if err := tx.First(&role, id).Error; err != nil {
			return err
		}
		if role.IsSystem == 1 {
			return errors.New("系统角色不可删除")
		}

		// 删除角色权限关联
		if err := tx.Exec("DELETE FROM role_permissions WHERE role_id = ?", id).Error; err != nil {
			return err
		}

		// 删除角色
		return tx.Delete(&entity.Role{}, id).Error
	})
}

// AssignPermissions 分配权限
func (r *RoleRepository) AssignPermissions(roleID uint, permissionIDs []uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 删除旧的权限关联
		if err := tx.Exec("DELETE FROM role_permissions WHERE role_id = ?", roleID).Error; err != nil {
			return err
		}

		// 插入新的权限关联
		for _, permID := range permissionIDs {
			if err := tx.Exec("INSERT INTO role_permissions (role_id, permission_id) VALUES (?, ?)", roleID, permID).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
