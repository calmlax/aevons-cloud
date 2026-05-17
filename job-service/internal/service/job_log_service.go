/**
 * 定时任务执行日志表 Service
 *
 * @author
 * @date 2026-04-19
 * @website https://www.aevons.com
 * @email 1073602@qq.com
 * @copyright 2025-2026 Aevons
 */
package service

import (
	"job-service/internal/dto"
	"job-service/internal/model"
	"job-service/internal/repository"

	"github.com/calmlax/aevons-framework/core/base"
)

// 继承BaseService
type JobLogService interface {
	base.BaseService[model.JobLog, *dto.JobLogQuery]
}

type jobLogService struct {
	base.BaseService[model.JobLog, *dto.JobLogQuery]
	repo repository.JobLogRepository
}

func NewJobLogService(repo repository.JobLogRepository) JobLogService {
	baseSrv := base.NewBaseService[model.JobLog, *dto.JobLogQuery](repo)
	return &jobLogService{
		BaseService: baseSrv,
		repo:        repo,
	}
}
