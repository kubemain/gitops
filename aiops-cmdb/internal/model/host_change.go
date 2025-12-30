package model

type HostChange struct {
	BaseModel
	HostID     uint   `json:"host_id" gorm:"index;not null"`
	Host       *Host  `json:"host,omitempty" gorm:"foreignKey:HostID"`
	ChangeType string `json:"change_type" gorm:"size:50;not null"`
	OldValue   string `json:"old_value" gorm:"type:text"`
	NewValue   string `json:"new_value" gorm:"type:text"`
	Operator   string `json:"operator" gorm:"size:100"`
	Remark     string `json:"remark" gorm:"type:text"`
}
