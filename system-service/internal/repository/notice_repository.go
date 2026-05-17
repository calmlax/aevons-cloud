package repository

/**
 * 通知公告 Repository
 *
 * @author
 * @date 2026-04-21
 * @website https://www.aevons.com
 * @email 1073602@qq.com
 * @copyright 2025-2026 Aevons
 */

import (
	"system-service/internal/model"

	"github.com/calmlax/aevons-framework/core/base"

	"gorm.io/gorm"
)

// 继承BaseRepository
type NoticeRepository interface {
	base.BaseRepository[model.Notice]
}

// noticeRepository 结构体
type noticeRepository struct {
	base.BaseRepository[model.Notice]
	db *gorm.DB
}

// 创建 NoticeRepository 实例
func NewNoticeRepository(db *gorm.DB) NoticeRepository {
	return &noticeRepository{
		BaseRepository: base.NewBaseRepository[model.Notice](db),
		db:             db,
	}
}
