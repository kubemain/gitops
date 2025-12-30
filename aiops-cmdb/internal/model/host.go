package model

import (
	"database/sql/driver"
	"encoding/json"
)

type Host struct {
	BaseModel
	Hostname    string     `json:"hostname" gorm:"size:100;not null;uniqueIndex"`
	IP          string     `json:"ip" gorm:"size:50;not null"`
	PublicIP    string     `json:"public_ip" gorm:"size:50"`
	OS          string     `json:"os" gorm:"size:50"`
	OSVersion   string     `json:"os_version" gorm:"size:50"`
	CPU         int        `json:"cpu"`
	Memory      int        `json:"memory"` // MB
	Disk        int        `json:"disk"`   // GB
	Status      string     `json:"status" gorm:"size:20;default:offline;index"`
	Environment string     `json:"environment" gorm:"size:20;index"`
	Region      string     `json:"region" gorm:"size:50"`
	IDC         string     `json:"idc" gorm:"size:50"`
	Cabinet     string     `json:"cabinet" gorm:"size:50"`
	GroupID     uint       `json:"group_id" gorm:"index"`
	Group       *HostGroup `json:"group,omitempty" gorm:"foreignKey:GroupID"`
	Tags        TagList    `json:"tags" gorm:"type:json"`
	Labels      JSONMap    `json:"labels" gorm:"type:json"`
	Remark      string     `json:"remark" gorm:"type:text"`
	AgentStatus string     `json:"agent_status" gorm:"size:20;default:offline"`
	LastSeenAt  *int64     `json:"last_seen_at"`
}

type TagList []string

func (t TagList) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (t *TagList) Scan(value interface{}) error {
	if value == nil {
		*t = []string{}
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, t)
}

type JSONMap map[string]interface{}

func (j JSONMap) Value() (driver.Value, error) {
	return json.Marshal(j)
}

func (j *JSONMap) Scan(value interface{}) error {
	if value == nil {
		*j = make(map[string]interface{})
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, j)
}

type HostQueryParams struct {
	QueryParams
	Status      string `form:"status"`
	Environment string `form:"environment"`
	GroupID     uint   `form:"group_id"`
	Region      string `form:"region"`
	IDC         string `form:"idc"`
	Tags        string `form:"tags"`
}

type HostStats struct {
	Total         int64            `json:"total"`
	Online        int64            `json:"online"`
	Offline       int64            `json:"offline"`
	Maintenance   int64            `json:"maintenance"`
	ByEnvironment map[string]int64 `json:"by_environment"`
	ByRegion      map[string]int64 `json:"by_region"`
	ByIDC         map[string]int64 `json:"by_idc"`
}
