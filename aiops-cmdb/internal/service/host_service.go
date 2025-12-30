package service

import (
	"aiops-cmdb/internal/model"
	"aiops-cmdb/internal/repository"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type HostService struct {
	repo *repository.HostRepository
}

func NewHostService() *HostService {
	return &HostService{
		repo: repository.NewHostRepository(),
	}
}

func (s *HostService) Create(host *model.Host) error {
	// 检查主机名是否已存在
	existing, err := s.repo.FindByHostname(host.Hostname)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if existing != nil {
		return fmt.Errorf("主机名 %s 已存在", host.Hostname)
	}

	// 设置默认值
	if host.Status == "" {
		host.Status = "offline"
	}
	if host.Environment == "" {
		host.Environment = "development"
	}
	if host.Tags == nil {
		host.Tags = []string{}
	}
	if host.Labels == nil {
		host.Labels = make(map[string]interface{})
	}

	return s.repo.Create(host)
}

func (s *HostService) Update(id uint, host *model.Host) error {
	existing, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	// 如果修改了主机名，检查是否重复
	if host.Hostname != "" && host.Hostname != existing.Hostname {
		duplicate, err := s.repo.FindByHostname(host.Hostname)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		if duplicate != nil {
			return fmt.Errorf("主机名 %s 已存在", host.Hostname)
		}
		existing.Hostname = host.Hostname
	}

	// 更新字段
	if host.IP != "" {
		existing.IP = host.IP
	}
	if host.PublicIP != "" {
		existing.PublicIP = host.PublicIP
	}
	if host.OS != "" {
		existing.OS = host.OS
	}
	if host.OSVersion != "" {
		existing.OSVersion = host.OSVersion
	}
	if host.CPU > 0 {
		existing.CPU = host.CPU
	}
	if host.Memory > 0 {
		existing.Memory = host.Memory
	}
	if host.Disk > 0 {
		existing.Disk = host.Disk
	}
	if host.Status != "" {
		existing.Status = host.Status
	}
	if host.Environment != "" {
		existing.Environment = host.Environment
	}
	if host.Region != "" {
		existing.Region = host.Region
	}
	if host.IDC != "" {
		existing.IDC = host.IDC
	}
	if host.Cabinet != "" {
		existing.Cabinet = host.Cabinet
	}
	if host.GroupID > 0 {
		existing.GroupID = host.GroupID
	}
	if host.Tags != nil {
		existing.Tags = host.Tags
	}
	if host.Labels != nil {
		existing.Labels = host.Labels
	}
	if host.Remark != "" {
		existing.Remark = host.Remark
	}

	return s.repo.Update(existing)
}

func (s *HostService) Delete(id uint) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	return s.repo.Delete(id)
}

func (s *HostService) GetByID(id uint) (*model.Host, error) {
	return s.repo.FindByID(id)
}

func (s *HostService) List(params *model.HostQueryParams) ([]*model.Host, int64, error) {
	return s.repo.List(params)
}

func (s *HostService) GetStats() (*model.HostStats, error) {
	return s.repo.GetStats()
}

func (s *HostService) BatchDelete(ids []uint) error {
	if len(ids) == 0 {
		return fmt.Errorf("请选择要删除的主机")
	}
	return s.repo.BatchDelete(ids)
}

func (s *HostService) UpdateStatus(id uint, status string) error {
	validStatus := map[string]bool{
		"online":      true,
		"offline":     true,
		"maintenance": true,
	}
	if !validStatus[status] {
		return fmt.Errorf("无效的状态: %s", status)
	}
	return s.repo.UpdateStatus(id, status)
}
