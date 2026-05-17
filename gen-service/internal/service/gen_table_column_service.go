/**
 * 代码生成表字段 Service
 *
 * @author
 * @date 2026-04-08 03:51:33.21431843 +0000 UTC
 * @website https://www.aevons.com
 * @email 1073602@qq.com
 * @copyright 2025-2026 Aevons
 */
package service

import (
	apperr "github.com/calmlax/aevons-framework/errors"

	"fmt"
	"gen-service/internal/dto"
	"gen-service/internal/model"
	"gen-service/internal/repository"

	"github.com/calmlax/aevons-framework/utils"

	"github.com/calmlax/aevons-framework/core/base"
)

// 继承BaseService
type GenTableColumnService interface {
	base.BaseService[model.GenTableColumn, *dto.GenTableColumnQuery]

	// BatchUpdate 批量更新数据表字段
	BatchUpdate(d []dto.UpdateGenTableColumnDTO) (bool, error)
}

type genTableColumnService struct {
	base.BaseService[model.GenTableColumn, *dto.GenTableColumnQuery]
	repo repository.GenTableColumnRepository
}

func NewGenTableColumnService(repo repository.GenTableColumnRepository) GenTableColumnService {
	baseSrv := base.NewBaseService[model.GenTableColumn, *dto.GenTableColumnQuery](repo)
	return &genTableColumnService{
		BaseService: baseSrv,
		repo:        repo,
	}
}

// BatchUpdate 批量更新数据表字段
func (s *genTableColumnService) BatchUpdate(columns []dto.UpdateGenTableColumnDTO) (bool, error) {
	// 空数据直接返回
	if len(columns) == 0 {
		return false, apperr.ErrorNoUpdateField
	}
	var updateMaps []map[string]any
	for _, d := range columns {
		fmt.Println(d)
		updateMaps = append(updateMaps, utils.StructToMapIgnoreNil(d))
	}
	return s.repo.BatchUpdate(updateMaps)
}
