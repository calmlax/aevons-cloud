/**
 * 菜单权限表 Repository
 *
 * @author
 * @date 2026-04-09 01:17:32.070579992 +0000 UTC
 * @website https://www.aevons.com
 * @email 1073602@qq.com
 * @copyright 2025-2026 Aevons
 */
package repository

import (
	"context"
	"fmt"
	"sys-service/internal/dto"
	"sys-service/internal/model"

	apperr "github.com/calmlax/aevons-framework/errors"

	"github.com/calmlax/aevons-framework/utils"

	"github.com/calmlax/aevons-framework/core/base"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// 继承BaseRepository
type MenuRepository interface {
	base.BaseRepository[model.Menu]
	CreateMenu(ctx context.Context, dto dto.CreateMenuDTO) (int64, error)
	UpdateMenu(ctx context.Context, dto dto.UpdateMenuDTO) error
	ListByLangCode(langCode string) ([]dto.MenuDTO, error)
	GetTranslationsById(menuId int64) ([]model.MenuTl, error)
	DeleteByIds(ctx context.Context, ids []int64) error
}

// menuRepository 结构体
type menuRepository struct {
	base.BaseRepository[model.Menu]
	db *gorm.DB
}

// 创建 MenuRepository 实例
func NewMenuRepository(db *gorm.DB) MenuRepository {
	return &menuRepository{
		BaseRepository: base.NewBaseRepository[model.Menu](db),
		db:             db,
	}
}

func (r *menuRepository) CreateMenu(ctx context.Context, dto dto.CreateMenuDTO) (int64, error) {

	menu := model.Menu{
		ParentId:   dto.ParentId,
		Type:       dto.Type,
		Status:     dto.Status,
		Sort:       dto.Sort,
		Path:       dto.Path,
		Component:  dto.Component,
		Query:      dto.Query,
		Visible:    dto.Visible,
		IsFrame:    dto.IsFrame,
		Permission: dto.Permission,
		Icon:       dto.Icon,
		ActiveId:   dto.ActiveId,
	}

	err := r.Transaction(ctx, func(tx *gorm.DB) error {

		// 直接插入（依赖唯一索引）
		if err := tx.Create(&menu).Error; err != nil {
			return fmt.Errorf("保存失败")
		}

		tls := make([]model.MenuTl, 0, len(dto.Translations))

		for langCode, trans := range dto.Translations {
			if utils.IsEmpty(trans) {
				continue
			}
			tls = append(tls, model.MenuTl{
				MenuId:   menu.Id,
				LangCode: langCode,
				Title:    trans,
			})
		}
		if len(tls) == 0 {
			return fmt.Errorf("至少填写一种语言的翻译")
		}

		if err := tx.Create(&tls).Error; err != nil {
			return fmt.Errorf("保存失败: %w", err)
		}

		return nil
	})

	return menu.Id, err
}

func (r *menuRepository) UpdateMenu(ctx context.Context, dto dto.UpdateMenuDTO) error {

	return r.Transaction(ctx, func(tx *gorm.DB) error {
		id := dto.Id
		if err := tx.Model(&model.Menu{}).
			Where("id = ?", id).
			Updates(map[string]any{
				"parent_id":  dto.ParentId,
				"type":       dto.Type,
				"status":     dto.Status,
				"sort":       dto.Sort,
				"path":       dto.Path,
				"component":  dto.Component,
				"query":      dto.Query,
				"visible":    dto.Visible,
				"is_frame":   dto.IsFrame,
				"permission": dto.Permission,
				"icon":       dto.Icon,
				"active_id":  dto.ActiveId,
			}).Error; err != nil {
			return fmt.Errorf("更新失败: %w", err)
		}

		//删除旧翻译
		if err := tx.
			Where("menu_id = ?", id).
			Delete(&model.MenuTl{}).Error; err != nil {
			return fmt.Errorf("更新失败: %w", err)
		}
		tls := make([]model.MenuTl, 0, len(dto.Translations))

		for langCode, trans := range dto.Translations {
			if utils.IsEmpty(trans) {
				continue
			}
			tls = append(tls, model.MenuTl{
				MenuId:   *id,
				LangCode: langCode,
				Title:    trans,
			})
		}
		if len(tls) == 0 {
			return fmt.Errorf("至少填写一种语言的翻译")
		}

		if err := tx.Create(&tls).Error; err != nil {
			return fmt.Errorf("保存失败: %w", err)
		}

		return nil
	})
}

func (r *menuRepository) ListByLangCode(langCode string) ([]dto.MenuDTO, error) {
	if utils.IsEmpty(langCode) {
		return []dto.MenuDTO{}, nil
	}
	var list []dto.MenuDTO
	err := r.db.Table("sys_menu").
		Select("id,parent_id,tl.lang_code,tl.title,type,sort,path,component,`query`,visible,`status`,is_frame,permission,icon,active_id,created_by,created_at,updated_by,updated_at").
		Joins("left join sys_menu_tl tl on sys_menu.id = tl.menu_id and lang_code = ?", langCode).
		Order(clause.OrderBy{Columns: []clause.OrderByColumn{
			{Column: clause.Column{Name: "parent_id"}, Desc: false},
			{Column: clause.Column{Name: "sort"}, Desc: false},
		}}).Find(&list).Error
	if err != nil {
		return []dto.MenuDTO{}, nil
	}
	return list, nil
}

func (r *menuRepository) GetTranslationsById(menuId int64) ([]model.MenuTl, error) {
	var menuTl []model.MenuTl
	err := r.db.Table("sys_menu_tl").
		Select("lang_code, title").
		Where("menu_id = ?", menuId).Find(&menuTl).Error
	return menuTl, err
}

// Delete 根据 ID 删除
func (r *menuRepository) DeleteByIds(ctx context.Context, ids []int64) error {

	return r.Transaction(ctx, func(tx *gorm.DB) error {
		var count int64
		e := tx.Model(&model.Menu{}).Where("parent_id in ?", ids).Count(&count).Error
		if e != nil || count > 0 {
			return apperr.ErrorExisting
		}
		if err := tx.
			Where("id in ?", ids).
			Delete(&model.Menu{}).Error; err != nil {
			return fmt.Errorf("删除失败: %w", err)
		}
		if err := tx.
			Where("menu_id in ?", ids).
			Delete(&model.MenuTl{}).Error; err != nil {
			return fmt.Errorf("删除失败: %w", err)
		}
		return nil
	})
}
