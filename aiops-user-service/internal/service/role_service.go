// internal/service/role_service.go
package service

import (
	"aiops-user-service/internal/domain/entity"
	"aiops-user-service/internal/repository"
)

type RoleService struct {
	roleRepo *repository.RoleRepository
	permRepo *repository.PermissionRepository
}

func NewRoleService() *RoleService {
	return &RoleService{
		roleRepo: repository.NewRoleRepository(),
		permRepo: repository.NewPermissionRepository(),
	}
}

// GetAll 获取所有角色
func (s *RoleService) GetAll() ([]*entity.Role, error) {
	return s.roleRepo.FindAll()
}

// GetByID 根据ID查询角色
func (s *RoleService) GetByID(id uint) (*entity.Role, error) {
	return s.roleRepo.FindByID(id)
}

// AssignPermissions 分配权限
func (s *RoleService) AssignPermissions(roleID uint, permissionIDs []uint) error {
	return s.roleRepo.AssignPermissions(roleID, permissionIDs)
}

// GetAllPermissions 获取所有权限（树形结构）
func (s *RoleService) GetAllPermissions() ([]*entity.Permission, error) {
	return s.permRepo.FindAll()
}
