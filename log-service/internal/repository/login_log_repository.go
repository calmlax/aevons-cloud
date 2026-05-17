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
	"log-service/internal/model"

	"github.com/calmlax/aevons-framework/core/base"

	"gorm.io/gorm"
)

type LoginLogRepository interface {
	base.BaseRepository[model.LoginLog]
	GetLatestLoginLog(username string, limit int) ([]model.LoginLog, error)
	TruncateAll() error
}

type loginLogRepository struct {
	base.BaseRepository[model.LoginLog]
	db *gorm.DB
}

func NewLoginLogRepository(db *gorm.DB) LoginLogRepository {
	return &loginLogRepository{
		BaseRepository: base.NewBaseRepository[model.LoginLog](db),
		db:             db,
	}
}

func (r *loginLogRepository) GetLatestLoginLog(username string, limit int) ([]model.LoginLog, error) {
	var list []model.LoginLog
	err := r.db.Model(&model.LoginLog{}).
		Where("username = ?", username).
		Order("login_at DESC").
		Limit(limit).
		Find(&list).Error
	return list, err
}

func (r *loginLogRepository) TruncateAll() error {
	return r.db.Exec("TRUNCATE TABLE sys_login_log").Error
}
