/**
 * 定时任务配置表 Repository
 *
 * @author
 * @date 2026-04-19
 * @website https://www.aevons.com
 * @email 1073602@qq.com
 * @copyright 2025-2026 Aevons
 */
package repository

import (
	"job-service/internal/model"

	"github.com/calmlax/aevons-framework/core/base"

	"gorm.io/gorm"
)

// 继承BaseRepository
type JobRepository interface {
	base.BaseRepository[model.Job]
}

// jobRepository 结构体
type jobRepository struct {
	base.BaseRepository[model.Job]
	db *gorm.DB
}

// 创建 JobRepository 实例
func NewJobRepository(db *gorm.DB) JobRepository {
	return &jobRepository{
		BaseRepository: base.NewBaseRepository[model.Job](db),
		db:             db,
	}
}
