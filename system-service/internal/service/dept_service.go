/**
 * 部门表 Service
 *
 * @author
 * @date 2026-04-09 01:02:34.478942247 +0000 UTC
 * @website https://www.aevons.com
 * @email 1073602@qq.com
 * @copyright 2025-2026 Aevons
 */
package service

import (
	"context"
	"system-service/internal/dto"
	"system-service/internal/model"
	"system-service/internal/repository"

	"github.com/calmlax/aevons-framework/core/base"
)

type DeptService interface {
	base.BaseService[model.Dept, *dto.DeptQuery]
	CreateDept(ctx context.Context, d dto.CreateDeptDTO) (*model.Dept, error)
	UpdateDept(ctx context.Context, id int64, d dto.UpdateDeptDTO) error
	ListTree(query dto.DeptQuery) ([]model.Dept, error)
	HasChildren(id int64) (bool, error)
}

type deptService struct {
	base.BaseService[model.Dept, *dto.DeptQuery]
	repo repository.DeptRepository
}

func NewDeptService(repo repository.DeptRepository) DeptService {
	return &deptService{
		BaseService: base.NewBaseService[model.Dept, *dto.DeptQuery](repo),
		repo:        repo,
	}
}

func (s *deptService) CreateDept(ctx context.Context, d dto.CreateDeptDTO) (*model.Dept, error) {
	return s.repo.CreateDept(ctx, d)
}

func (s *deptService) UpdateDept(ctx context.Context, id int64, d dto.UpdateDeptDTO) error {
	return s.repo.UpdateDept(ctx, id, d)
}

func (s *deptService) ListTree(query dto.DeptQuery) ([]model.Dept, error) {
	return s.repo.ListTree(query)
}

func (s *deptService) HasChildren(id int64) (bool, error) {
	return s.repo.HasChildren(id)
}
