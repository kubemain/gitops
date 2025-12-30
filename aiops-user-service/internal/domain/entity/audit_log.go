// internal/domain/entity/audit_log.go
package entity

import "time"

// AuditLog 操作日志
type AuditLog struct {
	ID             uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID         uint      `gorm:"type:bigint unsigned;default:0" json:"user_id"`
	Username       string    `gorm:"type:varchar(50)" json:"username"`
	Action         string    `gorm:"type:varchar(50);not null" json:"action"`
	Resource       string    `gorm:"type:varchar(50)" json:"resource"`
	ResourceID     uint      `gorm:"type:bigint unsigned;default:0" json:"resource_id"`
	Method         string    `gorm:"type:varchar(10)" json:"method"`
	Path           string    `gorm:"type:varchar(255)" json:"path"`
	IPAddress      string    `gorm:"type:varchar(50)" json:"ip_address"`
	UserAgent      string    `gorm:"type:varchar(500)" json:"user_agent"`
	RequestData    string    `gorm:"type:text" json:"request_data"`
	ResponseStatus int       `gorm:"type:int;default:0" json:"response_status"`
	ErrorMsg       string    `gorm:"type:text" json:"error_msg"`
	Duration       int       `gorm:"type:int;default:0" json:"duration"` // 毫秒
	CreatedAt      time.Time `gorm:"autoCreateTime" json:"created_at"`
}

// TableName 指定表名
func (AuditLog) TableName() string {
	return "audit_logs"
}
