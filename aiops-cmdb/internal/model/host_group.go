package model

type HostGroup struct {
	BaseModel
	Name        string `json:"name" gorm:"size:100;not null;uniqueIndex"`
	Description string `json:"description" gorm:"type:text"`
	ParentID    uint   `json:"parent_id" gorm:"index"`
	Sort        int    `json:"sort" gorm:"default:0"`
	HostCount   int    `json:"host_count" gorm:"-"`
}
