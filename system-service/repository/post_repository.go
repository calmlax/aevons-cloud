/**
 * 岗位信息表 Repository
 *
 * @author
 * @date 2026-04-09 01:37:54.290392924 +0000 UTC
 * @website https://www.aevons.com
 * @email 1073602@qq.com
 * @copyright 2025-2026 Aevons
 */
package repository

import (
	"system-service/model"

	"github.com/calmlax/aevons-framework/core/base"

	"gorm.io/gorm"
)

type PostRepository interface {
	base.BaseRepository[model.Post]
}

type postRepository struct {
	base.BaseRepository[model.Post]
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) PostRepository {
	return &postRepository{
		BaseRepository: base.NewBaseRepository[model.Post](db),
		db:             db,
	}
}
