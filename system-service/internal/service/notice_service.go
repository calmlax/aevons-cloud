package service

/**
 * 通知公告 Service
 *
 * @author
 * @date 2026-04-21
 * @website https://www.aevons.com
 * @email 1073602@qq.com
 * @copyright 2025-2026 Aevons
 */

import (
	"system-service/internal/dto"
	"system-service/internal/model"
	"system-service/internal/repository"

	"github.com/calmlax/aevons-framework/core/base"
)

// 继承BaseService
type NoticeService interface {
	base.BaseService[model.Notice, *dto.NoticeQuery]
}

type noticeService struct {
	base.BaseService[model.Notice, *dto.NoticeQuery]
	repo repository.NoticeRepository
}

func NewNoticeService(repo repository.NoticeRepository) NoticeService {
	baseSrv := base.NewBaseService[model.Notice, *dto.NoticeQuery](repo)
	return &noticeService{
		BaseService: baseSrv,
		repo:        repo,
	}
}
