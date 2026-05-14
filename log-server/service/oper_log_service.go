/**
 * 操作日志记录 Service
 *
 * @author
 * @date 2026-04-14
 * @website https://www.aevons.com
 * @email 1073602@qq.com
 * @copyright 2025-2026 Aevons
 */
package service

import (
	"log-server/dto"
	"log-server/model"
	"log-server/repository"

	"github.com/calmlax/aevons-framework/core/base"
)

type OperLogService interface {
	base.BaseService[model.OperLog, *dto.OperLogQuery]
	TruncateAll() error
}

type operLogService struct {
	base.BaseService[model.OperLog, *dto.OperLogQuery]
	repo repository.OperLogRepository
}

func NewOperLogService(repo repository.OperLogRepository) OperLogService {
	return &operLogService{
		BaseService: base.NewBaseService[model.OperLog, *dto.OperLogQuery](repo),
		repo:        repo,
	}
}

func (s *operLogService) TruncateAll() error {
	return s.repo.TruncateAll()
}
