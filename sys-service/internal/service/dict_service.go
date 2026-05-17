/**
 * 字典类型表 Service
 *
 * @author
 * @date 2026-04-09 01:08:50.442643548 +0000 UTC
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

// 继承BaseService
type DictService interface {
	base.BaseService[model.Dict, *dto.DictQuery]

	UpdateDict(ctx context.Context, id int64, updates map[string]any) error
}

type dictService struct {
	base.BaseService[model.Dict, *dto.DictQuery]
	repo repository.DictRepository
}

func NewDictService(repo repository.DictRepository) DictService {
	baseSrv := base.NewBaseService[model.Dict, *dto.DictQuery](repo)
	return &dictService{
		BaseService: baseSrv,
		repo:        repo,
	}
}

func (s *dictService) UpdateDict(ctx context.Context, id int64, updates map[string]any) error {

	return s.repo.UpdateDict(ctx, id, updates)
}
