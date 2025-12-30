package service

import (
	"aiops-cmdb/internal/model"
	"aiops-cmdb/internal/repository"
)

type HostGroupService struct {
	repo *repository.HostGroupRepository
}

func NewHostGroupService() *HostGroupService {
	return &HostGroupService{
		repo: repository.NewHostGroupRepository(),
	}
}

func (s *HostGroupService) Create(group *model.HostGroup) error {
	return s.repo.Create(group)
}

func (s *HostGroupService) Update(id uint, group *model.HostGroup) error {
	existing, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	if group.Name != "" {
		existing.Name = group.Name
	}
	if group.Description != "" {
		existing.Description = group.Description
	}
	if group.ParentID > 0 {
		existing.ParentID = group.ParentID
	}
	if group.Sort > 0 {
		existing.Sort = group.Sort
	}

	return s.repo.Update(existing)
}

func (s *HostGroupService) Delete(id uint) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	return s.repo.Delete(id)
}

func (s *HostGroupService) GetByID(id uint) (*model.HostGroup, error) {
	return s.repo.FindByID(id)
}

func (s *HostGroupService) List(params *model.QueryParams) ([]*model.HostGroup, int64, error) {
	return s.repo.List(params)
}

func (s *HostGroupService) GetAll() ([]*model.HostGroup, error) {
	return s.repo.GetAll()
}
