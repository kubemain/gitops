// internal/domain/entity/role.go
package entity

import "time"

// Role 角色实体
type Role struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string    `gorm:"type:varchar(50);not null" json:"name"`
	Code        string    `gorm:"type:varchar(50);uniqueIndex;not null" json:"code"`
	Description string    `gorm:"type:varchar(255)" json:"description"`
	SortOrder   int       `gorm:"type:int;default:0" json:"sort_order"`
	Status      int8      `gorm:"type:tinyint;default:1" json:"status"`
	IsSystem    int8      `gorm:"type:tinyint;default:0" json:"is_system"` // 是否系统角色
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// 关联关系
	Permissions []Permission `gorm:"many2many:role_permissions" json:"permissions,omitempty"`
}

// TableName 指定表名
func (Role) TableName() string {
	return "roles"
}
