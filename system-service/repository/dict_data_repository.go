/**
 * 字典数据表 Repository
 *
 * @author
 * @date 2026-04-09 01:08:50.443674979 +0000 UTC
 * @website https://www.aevons.com
 * @email 1073602@qq.com
 * @copyright 2025-2026 Aevons
 */
package repository

import (
	"context"
	"fmt"
	"system-service/dto"
	"system-service/model"

	apperr "github.com/calmlax/aevons-framework/errors"

	"github.com/calmlax/aevons-framework/utils"

	"github.com/calmlax/aevons-framework/core/base"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// 继承BaseRepository
type DictDataRepository interface {
	base.BaseRepository[model.DictData]
	CreateDictData(ctx context.Context, dto dto.CreateDictDataDTO) (int64, error)
	UpdateDictData(ctx context.Context, dto dto.UpdateDictDataDTO) error
	ListByDictType(dictType string) ([]dto.DictDataDTO, error)
	ListByDictTypeAndLangCode(dictType string, langCode string) ([]dto.DictDataDTO, error)
	GetTranslationsById(dictDataId int64) ([]model.DictDataTl, error)
	DeleteByIds(ctx context.Context, ids []int64) error
	UpdateSort(ctx context.Context, items []dto.SortItemDTO) error
}

// dictDataRepository 结构体
type dictDataRepository struct {
	base.BaseRepository[model.DictData]
	db *gorm.DB
}

// 创建 DictDataRepository 实例
func NewDictDataRepository(db *gorm.DB) DictDataRepository {
	return &dictDataRepository{
		BaseRepository: base.NewBaseRepository[model.DictData](db),
		db:             db,
	}
}
func (r *dictDataRepository) CreateDictData(ctx context.Context, dto dto.CreateDictDataDTO) (int64, error) {

	dictData := model.DictData{
		DictType:  dto.DictType,
		DictValue: dto.DictValue,
		Status:    dto.Status,
		Sort:      dto.Sort,
		TagType:   dto.TagType,
		TagClass:  dto.TagClass,
	}

	err := r.Transaction(ctx, func(tx *gorm.DB) error {
		var count int64
		tx.Model(model.DictData{}).Where("dict_type = ? and dict_value = ?", dictData.DictType, dictData.DictValue).Count(&count)
		if count > 0 {
			return apperr.ErrorExisting
		}
		// 直接插入（依赖唯一索引）
		if err := tx.Create(&dictData).Error; err != nil {
			return fmt.Errorf("保存字典[%s:%s]失败: %w", dictData.DictType, dictData.DictValue, err)
		}

		tls := make([]model.DictDataTl, 0, len(dto.Translations))
		for langCode, trans := range dto.Translations {
			if utils.IsEmpty(trans.Label) {
				continue
			}
			tls = append(tls, model.DictDataTl{
				DictDataId: dictData.Id,
				LangCode:   langCode,
				Label:      trans.Label,
				Tip:        trans.Tip,
			})
		}
		if len(tls) == 0 {
			return fmt.Errorf("至少填写一种语言的翻译")
		}
		if err := tx.Create(&tls).Error; err != nil {
			return fmt.Errorf("保存翻译失败: %w", err)
		}

		return nil
	})

	return dictData.Id, err
}

func (r *dictDataRepository) UpdateDictData(ctx context.Context, dto dto.UpdateDictDataDTO) error {

	return r.Transaction(ctx, func(tx *gorm.DB) error {
		id := dto.Id
		if err := tx.Model(&model.DictData{}).
			Where("id = ?", id).
			Updates(map[string]any{
				"dict_type":  dto.DictType,
				"dict_value": dto.DictValue,
				"status":     dto.Status,
				"sort":       dto.Sort,
				"tag_type":   dto.TagType,
				"tag_class":  dto.TagClass,
			}).Error; err != nil {
			return fmt.Errorf("更新字典失败: %w", err)
		}

		//删除旧翻译
		if err := tx.
			Where("dict_data_id = ?", id).
			Delete(&model.DictDataTl{}).Error; err != nil {
			return fmt.Errorf("删除翻译失败: %w", err)
		}
		//插入新翻译
		if len(dto.Translations) == 0 {
			return fmt.Errorf("必须填写标签信息")
		}

		tls := make([]model.DictDataTl, 0, len(dto.Translations))
		for langCode, trans := range dto.Translations {
			if utils.IsEmpty(trans.Label) {
				continue
			}
			tls = append(tls, model.DictDataTl{
				DictDataId: *id,
				LangCode:   langCode,
				Label:      trans.Label,
				Tip:        trans.Tip,
			})
		}
		if len(tls) == 0 {
			return fmt.Errorf("至少填写一种语言的翻译")
		}

		if err := tx.Create(&tls).Error; err != nil {
			return fmt.Errorf("保存翻译失败: %w", err)
		}

		return nil
	})
}
func (r *dictDataRepository) ListByDictTypeAndLangCode(dictType string, langCode string) ([]dto.DictDataDTO, error) {

	var ddlist []model.DictData

	err := r.db.Model(&model.DictData{}).Where("dict_type = ?", dictType).
		Order(clause.OrderBy{Columns: []clause.OrderByColumn{
			{Column: clause.Column{Name: "dict_type"}, Desc: false},
			{Column: clause.Column{Name: "sort"}, Desc: false},
		}}).
		Find(&ddlist).Error
	if err != nil {
		return []dto.DictDataDTO{}, err
	}
	var list []dto.DictDataDTO
	for _, item := range ddlist {
		var dictDataTl model.DictDataTl
		r.db.Table("sys_dict_data_tl").
			Select("lang_code, label, tip").
			Where("dict_data_id = ? and lang_code = ?", item.Id, langCode).First(&dictDataTl)

		list = append(list, dto.DictDataDTO{
			Id:        item.Id,
			DictType:  item.DictType,
			DictValue: item.DictValue,
			Status:    item.Status,
			Sort:      item.Sort,
			TagType:   item.TagType,
			TagClass:  item.TagClass,
			LangCode:  dictDataTl.LangCode,
			Label:     dictDataTl.Label,
			Tip:       dictDataTl.Tip,
		})
	}
	return list, nil
}

func (r *dictDataRepository) ListByDictType(dictType string) ([]dto.DictDataDTO, error) {
	var dictData []dto.DictDataDTO
	db := r.db.Table("sys_dict_data").
		Select("sys_dict_data.id, sys_dict_data.dict_type,sys_dict_data.dict_value, sys_dict_data.status, sys_dict_data.sort, sys_dict_data.tag_type, sys_dict_data.tag_class, sys_dict_data_tl.lang_code, sys_dict_data_tl.label, sys_dict_data_tl.tip").
		Joins("INNER JOIN sys_dict_data_tl ON sys_dict_data.id = sys_dict_data_tl.dict_data_id").
		Where("sys_dict_data.status = ?", 0)

	if utils.IsNotEmpty(dictType) {
		db = db.Where("sys_dict_data.dict_type = ?", dictType)
	}
	db = db.Order(clause.OrderBy{Columns: []clause.OrderByColumn{
		{Column: clause.Column{Name: "dict_type"}, Desc: false},
		{Column: clause.Column{Name: "sort"}, Desc: false},
	}})
	err := db.Find(&dictData).Error
	return dictData, err
}

func (r *dictDataRepository) GetTranslationsById(dictDataId int64) ([]model.DictDataTl, error) {
	var dictDataTl []model.DictDataTl
	err := r.db.Table("sys_dict_data_tl").
		Select("lang_code, label, tip").
		Where("dict_data_id = ?", dictDataId).Find(&dictDataTl).Error
	return dictDataTl, err
}

// Delete 根据 ID 删除
func (r *dictDataRepository) DeleteByIds(ctx context.Context, ids []int64) error {

	return r.Transaction(ctx, func(tx *gorm.DB) error {
		if err := tx.
			Where("id in ?", ids).
			Delete(&model.DictData{}).Error; err != nil {
			return fmt.Errorf("删除失败: %w", err)
		}
		if err := tx.
			Where("dict_data_id in ?", ids).
			Delete(&model.DictDataTl{}).Error; err != nil {
			return fmt.Errorf("删除失败: %w", err)
		}
		return nil
	})
}

// UpdateSort 批量更新排序
func (r *dictDataRepository) UpdateSort(ctx context.Context, items []dto.SortItemDTO) error {
	return r.Transaction(ctx, func(tx *gorm.DB) error {
		for _, item := range items {
			if err := tx.Model(&model.DictData{}).
				Where("id = ?", item.Id).
				Update("sort", item.Sort).Error; err != nil {
				return fmt.Errorf("更新排序失败: %w", err)
			}
		}
		return nil
	})
}
