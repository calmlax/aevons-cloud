/**
 * 字典类型表 Repository
 *
 * @author
 * @date 2026-04-09 01:08:50.442643548 +0000 UTC
 * @website https://www.aevons.com
 * @email 1073602@qq.com
 * @copyright 2025-2026 Aevons
 */
package repository

import (
	"context"
	"fmt"
	"system-service/model"

	"github.com/calmlax/aevons-framework/core/base"

	"gorm.io/gorm"
)

// 继承BaseRepository
type DictRepository interface {
	base.BaseRepository[model.Dict]

	UpdateDict(ctx context.Context, id int64, updates map[string]any) error
}

// dictRepository 结构体
type dictRepository struct {
	base.BaseRepository[model.Dict]
	db *gorm.DB
}

// 创建 DictRepository 实例
func NewDictRepository(db *gorm.DB) DictRepository {
	return &dictRepository{
		BaseRepository: base.NewBaseRepository[model.Dict](db),
		db:             db,
	}
}

func (r *dictRepository) UpdateDict(ctx context.Context, id int64, updates map[string]any) error {

	return r.Transaction(ctx, func(tx *gorm.DB) error {
		var entity model.Dict
		if err := tx.First(&entity, id).Error; err != nil {
			return err
		}
		var oldDictType = entity.DictType
		if err := r.db.Model(&entity).Updates(updates).Error; err != nil {
			return err
		}
		var dictType = updates["dictType"]
		if dictType != oldDictType {
			if err := tx.Model(&model.DictData{}).
				Where("dict_type = ?", oldDictType).
				Updates(map[string]any{
					"dict_type": dictType,
				}).Error; err != nil {
				return fmt.Errorf("更新字典数据失败: %w", err)
			}
		}
		return nil
	})
}
