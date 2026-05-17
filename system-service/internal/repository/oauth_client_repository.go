/**
 * 终端应用 Repository
 *
 * @author
 * @date 2026-04-09 01:26:40.390618977 +0000 UTC
 * @website https://www.aevons.com
 * @email 1073602@qq.com
 * @copyright 2025-2026 Aevons
 */
package repository

import (
	"system-service/internal/model"

	"github.com/calmlax/aevons-framework/core/base"

	"gorm.io/gorm"
)

// 继承BaseRepository
type OauthClientRepository interface {
	base.BaseRepository[model.OauthClient]
}

// oauthClientRepository 结构体
type oauthClientRepository struct {
	base.BaseRepository[model.OauthClient]
	db *gorm.DB
}

// 创建 OauthClientRepository 实例
func NewOauthClientRepository(db *gorm.DB) OauthClientRepository {
	return &oauthClientRepository{
		BaseRepository: base.NewBaseRepository[model.OauthClient](db),
		db:             db,
	}
}
