/**
 * 定时任务执行日志表 Repository
 *
 * @author
 * @date 2026-04-19
 * @website https://www.aevons.com
 * @email 1073602@qq.com
 * @copyright 2025-2026 Aevonsns
 */
package repository

import (
	"job-service/internal/model"

	"github.com/calmlax/aevons-framework/core/base"

	"gorm.io/gorm"
)

// 继承BaseRepository
type JobLogRepository interface {
	base.BaseRepository[model.JobLog]
}

// jobLogRepository 结构体
type jobLogRepository struct {
	base.BaseRepository[model.JobLog]
	db *gorm.DB
}

// 创建 JobLogRepository 实例
func NewJobLogRepository(db *gorm.DB) JobLogRepository {
	return &jobLogRepository{
		BaseRepository: base.NewBaseRepository[model.JobLog](db),
		db:             db,
	}
}
