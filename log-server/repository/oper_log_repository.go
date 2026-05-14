/**
 * 操作日志记录 Repository
 *
 * @author
 * @date 2026-04-14
 * @website https://www.aevons.com
 * @email 1073602@qq.com
 * @copyright 2025-2026 Aevons
 */
package repository

import (
	"log-server/model"

	"github.com/calmlax/aevons-framework/core/base"

	"gorm.io/gorm"
)

type OperLogRepository interface {
	base.BaseRepository[model.OperLog]
	TruncateAll() error
}

type operLogRepository struct {
	base.BaseRepository[model.OperLog]
	db *gorm.DB
}

func NewOperLogRepository(db *gorm.DB) OperLogRepository {
	return &operLogRepository{
		BaseRepository: base.NewBaseRepository[model.OperLog](db),
		db:             db,
	}
}

func (r *operLogRepository) TruncateAll() error {
	return r.db.Exec("TRUNCATE TABLE sys_oper_log").Error
}
