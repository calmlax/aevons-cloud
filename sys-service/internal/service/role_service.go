/**
 * 角色信息表 Service
 *
 * @author
 * @date 2026-04-09 01:49:40.926495422 +0000 UTC
 * @website https://www.aevons.com
 * @email 1073602@qq.com
 * @copyright 2025-2026 Aevons
 */
package service

import (
	"context"
	"sys-service/internal/dto"
	"sys-service/internal/model"
	"sys-service/internal/repository"

	"github.com/calmlax/aevons-framework/core/base"
)

type RoleService interface {
	base.BaseService[model.Role, *dto.RoleQuery]
	GetMenuIds(roleId int64) ([]int64, error)
	CreateWithMenus(ctx context.Context, d dto.CreateRoleDTO) (*model.Role, error)
	UpdateWithMenus(ctx context.Context, id int64, d dto.UpdateRoleDTO) error
}

type roleService struct {
	base.BaseService[model.Role, *dto.RoleQuery]
	repo repository.RoleRepository
}

func NewRoleService(repo repository.RoleRepository) RoleService {
	return &roleService{
		BaseService: base.NewBaseService[model.Role, *dto.RoleQuery](repo),
		repo:        repo,
	}
}

func (s *roleService) GetMenuIds(roleId int64) ([]int64, error) {
	return s.repo.GetMenuIds(roleId)
}

func (s *roleService) CreateWithMenus(ctx context.Context, d dto.CreateRoleDTO) (*model.Role, error) {
	return s.repo.CreateWithMenus(ctx, d)
}

func (s *roleService) UpdateWithMenus(ctx context.Context, id int64, d dto.UpdateRoleDTO) error {
	return s.repo.UpdateWithMenus(ctx, id, d)
}
