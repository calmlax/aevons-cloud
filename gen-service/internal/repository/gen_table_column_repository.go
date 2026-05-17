/**
 * 代码生成表字段 Repository
 *
 * @author
 * @date 2026-04-08 03:51:33.21431843 +0000 UTC
 * @website https://www.aevons.com
 * @email 1073602@qq.com
 * @copyright 2025-2026 Aevons
 */
package repository

import (
	"errors"
	"gen-service/internal/dto"
	"gen-service/internal/model"

	apperr "github.com/calmlax/aevons-framework/errors"

	"github.com/calmlax/aevons-framework/core/base"

	"gorm.io/gorm"
)

// 继承BaseRepository
type GenTableColumnRepository interface {
	base.BaseRepository[model.GenTableColumn]

	// BatchUpdate 根据ID更新数据表字段
	BatchUpdate(updates []map[string]any) (bool, error)

	// DeleteByTableIds 根据表ID字段
	DeleteByTableIds(tableIds []int64) error

	// ListByTableId 根据tableId条件查询数据表字段列表
	ListByTableId(tableId int64) ([]model.GenTableColumn, error)
}

// genTableColumnRepository 结构体
type genTableColumnRepository struct {
	base.BaseRepository[model.GenTableColumn]
	db *gorm.DB
}

// 创建 GenTableColumnRepository 实例
func NewGenTableColumnRepository(db *gorm.DB) GenTableColumnRepository {
	return &genTableColumnRepository{
		BaseRepository: base.NewBaseRepository[model.GenTableColumn](db),
		db:             db,
	}
}

// BatchUpdate 根据ID更新数据表字段
func (s *genTableColumnRepository) BatchUpdate(updates []map[string]any) (bool, error) {
	if len(updates) == 0 {
		return false, apperr.ErrorNoUpdateField
	}

	tx := s.db.Begin()
	if tx.Error != nil {
		return false, tx.Error
	}
	for _, column := range updates {
		columnId, ok := column["id"].(int64)
		if !ok {
			tx.Rollback()
			return false, errors.New("字段ID不存在或格式错误")
		}
		delete(column, "id")
		result := tx.Model(&model.GenTableColumn{}).
			Where("id = ?", columnId).
			Updates(column)

		if result.Error != nil {
			tx.Rollback()
			return false, result.Error
		}
	}
	// 4. 提交事务
	if err := tx.Commit().Error; err != nil {
		return false, err
	}

	return true, nil
}

// DeleteByTableIds 根据表ID批量删除字段（级联删除）
func (r *genTableColumnRepository) DeleteByTableIds(tableIds []int64) error {
	return r.db.Where("table_id IN ?", tableIds).Delete(&model.GenTableColumn{}).Error
}

// ListByTableId 根据tableId条件查询数据表字段列表
func (r *genTableColumnRepository) ListByTableId(tableId int64) ([]model.GenTableColumn, error) {
	q := dto.GenTableColumnQuery{
		TableId: tableId,
	}
	return r.BaseRepository.List(&q)
}
