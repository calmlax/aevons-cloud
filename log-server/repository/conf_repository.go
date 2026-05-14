/**
 * 参数配置表 Repository
 *
 * @author
 * @date 2026-04-09 00:38:25.504785055 +0000 UTC
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

// 继承BaseRepository
type ConfRepository interface {
	base.BaseRepository[model.Conf]
}

// confRepository 结构体
type confRepository struct {
	base.BaseRepository[model.Conf]
	db *gorm.DB
}

// 创建 ConfRepository 实例
func NewConfRepository(db *gorm.DB) ConfRepository {
	return &confRepository{
		BaseRepository: base.NewBaseRepository[model.Conf](db),
		db:             db,
	}
}
