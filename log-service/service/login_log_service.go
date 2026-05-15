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
	"log-service/dto"
	"log-service/model"
	"log-service/repository"

	"github.com/calmlax/aevons-framework/core/base"
)

type LoginLogService interface {
	base.BaseService[model.LoginLog, *dto.LoginLogQuery]
	GetProfileLoginLogs(username string) ([]model.LoginLog, error)
	TruncateAll() error
}

type loginLogService struct {
	base.BaseService[model.LoginLog, *dto.LoginLogQuery]
	repo repository.LoginLogRepository
}

func NewLoginLogService(repo repository.LoginLogRepository) LoginLogService {
	return &loginLogService{
		BaseService: base.NewBaseService[model.LoginLog, *dto.LoginLogQuery](repo),
		repo:        repo,
	}
}

func (s *loginLogService) GetProfileLoginLogs(username string) ([]model.LoginLog, error) {
	return s.repo.GetLatestLoginLogs(username, 10)
}

func (s *loginLogService) TruncateAll() error {
	return s.repo.TruncateAll()
}
