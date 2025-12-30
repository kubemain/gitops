// internal/domain/entity/permission.go
package entity

import "time"

// Permission 权限实体
type Permission struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	ParentID  uint      `gorm:"type:bigint unsigned;default:0" json:"parent_id"`
	Name      string    `gorm:"type:varchar(50);not null" json:"name"`
	Code      string    `gorm:"type:varchar(100);uniqueIndex;not null" json:"code"`
	Type      string    `gorm:"type:varchar(20);default:'menu'" json:"type"` // menu/button/api
	Resource  string    `gorm:"type:varchar(50)" json:"resource"`
	Action    string    `gorm:"type:varchar(50)" json:"action"`
	Path      string    `gorm:"type:varchar(200)" json:"path"`
	Method    string    `gorm:"type:varchar(10)" json:"method"`
	SortOrder int       `gorm:"type:int;default:0" json:"sort_order"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName 指定表名
func (Permission) TableName() string {
	return "permissions"
}
