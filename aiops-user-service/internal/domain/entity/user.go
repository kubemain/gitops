// internal/domain/entity/user.go
package entity

import (
	"time"
)

// User 用户实体
type User struct {
	ID           uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	Username     string     `gorm:"type:varchar(50);uniqueIndex;not null" json:"username"`
	Nickname     string     `gorm:"type:varchar(50)" json:"nickname"`
	Email        string     `gorm:"type:varchar(100)" json:"email"`
	Phone        string     `gorm:"type:varchar(20)" json:"phone"`
	Avatar       string     `gorm:"type:varchar(255)" json:"avatar"`
	PasswordHash string     `gorm:"type:varchar(255);not null" json:"-"`  // 不返回给前端
	Status       int8       `gorm:"type:tinyint;default:1" json:"status"` // 1=启用 0=禁用
	LastLoginAt  *time.Time `gorm:"type:timestamp" json:"last_login_at"`
	LastLoginIP  string     `gorm:"type:varchar(50)" json:"last_login_ip"`
	CreatedAt    time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time  `gorm:"autoUpdateTime" json:"updated_at"`

	// 关联关系
	Roles []Role `gorm:"many2many:user_roles" json:"roles,omitempty"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}
