// internal/repository/user_repo.go
package repository

import (
	"aiops-user-service/internal/domain/entity"
	"aiops-user-service/pkg/database"
	"errors"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository() *UserRepository {
	return &UserRepository{db: database.DB}
}

// Create 创建用户
func (r *UserRepository) Create(user *entity.User) error {
	return r.db.Create(user).Error
}

// FindByID 根据ID查询用户
func (r *UserRepository) FindByID(id uint) (*entity.User, error) {
	var user entity.User
	err := r.db.Preload("Roles").First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}
	return &user, nil
}

// FindByUsername 根据用户名查询
func (r *UserRepository) FindByUsername(username string) (*entity.User, error) {
	var user entity.User
	err := r.db.Preload("Roles").Where("username = ?", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}
	return &user, nil
}

// Update 更新用户
func (r *UserRepository) Update(user *entity.User) error {
	return r.db.Save(user).Error
}

// Delete 删除用户
func (r *UserRepository) Delete(id uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 删除用户角色关联
		if err := tx.Exec("DELETE FROM user_roles WHERE user_id = ?", id).Error; err != nil {
			return err
		}
		// 删除用户
		return tx.Delete(&entity.User{}, id).Error
	})
}

// List 用户列表（分页）
func (r *UserRepository) List(page, pageSize int, filters map[string]interface{}) ([]*entity.User, int64, error) {
	var users []*entity.User
	var total int64

	query := r.db.Model(&entity.User{})

	// 应用过滤条件
	if username, ok := filters["username"].(string); ok && username != "" {
		query = query.Where("username LIKE ?", "%"+username+"%")
	}
	if status, ok := filters["status"].(int); ok {
		query = query.Where("status = ?", status)
	}

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err := query.Preload("Roles").
		Offset(offset).
		Limit(pageSize).
		Order("id DESC").
		Find(&users).Error

	return users, total, err
}

// AssignRoles 分配角色
func (r *UserRepository) AssignRoles(userID uint, roleIDs []uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 先删除旧的角色关联
		if err := tx.Exec("DELETE FROM user_roles WHERE user_id = ?", userID).Error; err != nil {
			return err
		}

		// 插入新的角色关联
		for _, roleID := range roleIDs {
			if err := tx.Exec("INSERT INTO user_roles (user_id, role_id) VALUES (?, ?)", userID, roleID).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// UpdateLastLogin 更新最后登录信息
func (r *UserRepository) UpdateLastLogin(userID uint, ip string) error {
	return r.db.Model(&entity.User{}).
		Where("id = ?", userID).
		Updates(map[string]interface{}{
			"last_login_at": gorm.Expr("NOW()"),
			"last_login_ip": ip,
		}).Error
}
