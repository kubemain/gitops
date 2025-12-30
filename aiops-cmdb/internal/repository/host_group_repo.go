package repository

import (
	"aiops-cmdb/internal/model"
	"aiops-cmdb/pkg/database"

	"gorm.io/gorm"
)

type HostGroupRepository struct {
	db *gorm.DB
}

func NewHostGroupRepository() *HostGroupRepository {
	return &HostGroupRepository{
		db: database.GetDB(),
	}
}

func (r *HostGroupRepository) Create(group *model.HostGroup) error {
	return r.db.Create(group).Error
}

func (r *HostGroupRepository) Update(group *model.HostGroup) error {
	return r.db.Save(group).Error
}

func (r *HostGroupRepository) Delete(id uint) error {
	return r.db.Delete(&model.HostGroup{}, id).Error
}

func (r *HostGroupRepository) FindByID(id uint) (*model.HostGroup, error) {
	var group model.HostGroup
	err := r.db.First(&group, id).Error
	if err != nil {
		return nil, err
	}
	return &group, nil
}

func (r *HostGroupRepository) List(params *model.QueryParams) ([]*model.HostGroup, int64, error) {
	var groups []*model.HostGroup
	var total int64

	query := r.db.Model(&model.HostGroup{})

	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where("name LIKE ? OR description LIKE ?", keyword, keyword)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.
		Order(params.GetOrder()).
		Limit(params.GetPageSize()).
		Offset(params.GetOffset()).
		Find(&groups).Error

	// 统计每个分组的主机数
	for _, group := range groups {
		var count int64
		r.db.Model(&model.Host{}).Where("group_id = ?", group.ID).Count(&count)
		group.HostCount = int(count)
	}

	return groups, total, err
}

func (r *HostGroupRepository) GetAll() ([]*model.HostGroup, error) {
	var groups []*model.HostGroup
	err := r.db.Order("sort asc, id asc").Find(&groups).Error
	return groups, err
}
