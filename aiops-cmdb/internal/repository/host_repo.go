package repository

import (
	"aiops-cmdb/internal/model"
	"aiops-cmdb/pkg/database"
	"fmt"

	"gorm.io/gorm"
)

type HostRepository struct {
	db *gorm.DB
}

func NewHostRepository() *HostRepository {
	return &HostRepository{
		db: database.GetDB(),
	}
}

func (r *HostRepository) Create(host *model.Host) error {
	return r.db.Create(host).Error
}

func (r *HostRepository) Update(host *model.Host) error {
	return r.db.Save(host).Error
}

func (r *HostRepository) Delete(id uint) error {
	return r.db.Delete(&model.Host{}, id).Error
}

func (r *HostRepository) FindByID(id uint) (*model.Host, error) {
	var host model.Host
	err := r.db.Preload("Group").First(&host, id).Error
	if err != nil {
		return nil, err
	}
	return &host, nil
}

func (r *HostRepository) FindByHostname(hostname string) (*model.Host, error) {
	var host model.Host
	err := r.db.Where("hostname = ?", hostname).First(&host).Error
	if err != nil {
		return nil, err
	}
	return &host, nil
}

func (r *HostRepository) List(params *model.HostQueryParams) ([]*model.Host, int64, error) {
	var hosts []*model.Host
	var total int64

	query := r.db.Model(&model.Host{}).Preload("Group")

	// 关键字搜索
	if params.Keyword != "" {
		keyword := "%" + params.Keyword + "%"
		query = query.Where("hostname LIKE ? OR ip LIKE ? OR public_ip LIKE ?", keyword, keyword, keyword)
	}

	// 状态过滤
	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}

	// 环境过滤
	if params.Environment != "" {
		query = query.Where("environment = ?", params.Environment)
	}

	// 分组过滤
	if params.GroupID > 0 {
		query = query.Where("group_id = ?", params.GroupID)
	}

	// 区域过滤
	if params.Region != "" {
		query = query.Where("region = ?", params.Region)
	}

	// IDC 过滤
	if params.IDC != "" {
		query = query.Where("idc = ?", params.IDC)
	}

	// 标签过滤
	if params.Tags != "" {
		query = query.Where("JSON_CONTAINS(tags, ?)", fmt.Sprintf(`"%s"`, params.Tags))
	}

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	err := query.
		Order(params.GetOrder()).
		Limit(params.GetPageSize()).
		Offset(params.GetOffset()).
		Find(&hosts).Error

	return hosts, total, err
}

func (r *HostRepository) GetStats() (*model.HostStats, error) {
	stats := &model.HostStats{
		ByEnvironment: make(map[string]int64),
		ByRegion:      make(map[string]int64),
		ByIDC:         make(map[string]int64),
	}

	// 总数
	r.db.Model(&model.Host{}).Count(&stats.Total)

	// 按状态统计
	r.db.Model(&model.Host{}).Where("status = ?", "online").Count(&stats.Online)
	r.db.Model(&model.Host{}).Where("status = ?", "offline").Count(&stats.Offline)
	r.db.Model(&model.Host{}).Where("status = ?", "maintenance").Count(&stats.Maintenance)

	// 按环境统计
	var envStats []struct {
		Environment string
		Count       int64
	}
	r.db.Model(&model.Host{}).
		Select("environment, COUNT(*) as count").
		Group("environment").
		Scan(&envStats)
	for _, s := range envStats {
		stats.ByEnvironment[s.Environment] = s.Count
	}

	// 按区域统计
	var regionStats []struct {
		Region string
		Count  int64
	}
	r.db.Model(&model.Host{}).
		Select("region, COUNT(*) as count").
		Where("region != ''").
		Group("region").
		Scan(&regionStats)
	for _, s := range regionStats {
		stats.ByRegion[s.Region] = s.Count
	}

	// 按 IDC 统计
	var idcStats []struct {
		IDC   string
		Count int64
	}
	r.db.Model(&model.Host{}).
		Select("idc, COUNT(*) as count").
		Where("idc != ''").
		Group("idc").
		Scan(&idcStats)
	for _, s := range idcStats {
		stats.ByIDC[s.IDC] = s.Count
	}

	return stats, nil
}

func (r *HostRepository) BatchDelete(ids []uint) error {
	return r.db.Delete(&model.Host{}, ids).Error
}

func (r *HostRepository) UpdateStatus(id uint, status string) error {
	return r.db.Model(&model.Host{}).Where("id = ?", id).Update("status", status).Error
}
