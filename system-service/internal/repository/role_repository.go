/**
 * 角色信息表 Repository
 *
 * @author
 * @date 2026-04-09 01:49:40.926495422 +0000 UTC
 * @website https://www.aevons.com
 * @email 1073602@qq.com
 * @copyright 2025-2026 Aevons
 */
package repository

import (
	"context"
	"system-service/internal/dto"
	"system-service/internal/model"

	"github.com/calmlax/aevons-framework/utils"

	"github.com/calmlax/aevons-framework/core/base"

	"gorm.io/gorm"
)

type RoleRepository interface {
	base.BaseRepository[model.Role]
	GetMenuIds(roleId int64) ([]int64, error)
	CreateWithMenus(ctx context.Context, d dto.CreateRoleDTO) (*model.Role, error)
	UpdateWithMenus(ctx context.Context, id int64, d dto.UpdateRoleDTO) error
}

type roleRepository struct {
	base.BaseRepository[model.Role]
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &roleRepository{
		BaseRepository: base.NewBaseRepository[model.Role](db),
		db:             db,
	}
}

func (r *roleRepository) GetMenuIds(roleId int64) ([]int64, error) {
	var ids []int64
	err := r.db.Model(&model.RoleMenu{}).
		Where("role_id = ?", roleId).
		Pluck("menu_id", &ids).Error
	return ids, err
}

// expandWithParents 父子联动模式下，自动把所有祖先节点 ID 补进来
func expandWithParents(db *gorm.DB, menuIds []int64) []int64 {
	all := make(map[int64]struct{})
	for _, id := range menuIds {
		all[id] = struct{}{}
	}
	queue := append([]int64{}, menuIds...)
	for len(queue) > 0 {
		var parents []struct {
			Id       int64 `gorm:"column:id"`
			ParentId int64 `gorm:"column:parent_id"`
		}
		db.Table("sys_menu").Select("id, parent_id").
			Where("id IN ? AND parent_id != 0", queue).Find(&parents)
		queue = queue[:0]
		for _, p := range parents {
			if p.ParentId != 0 {
				if _, exists := all[p.ParentId]; !exists {
					all[p.ParentId] = struct{}{}
					queue = append(queue, p.ParentId)
				}
			}
		}
	}
	result := make([]int64, 0, len(all))
	for id := range all {
		result = append(result, id)
	}
	return result
}

func saveRoleMenus(tx *gorm.DB, roleId int64, menuIds []int64) error {
	if err := tx.Where("role_id = ?", roleId).Delete(&model.RoleMenu{}).Error; err != nil {
		return err
	}
	if len(menuIds) == 0 {
		return nil
	}
	menuIds = expandWithParents(tx, menuIds)
	records := make([]model.RoleMenu, 0, len(menuIds))
	for _, mid := range menuIds {
		records = append(records, model.RoleMenu{RoleId: roleId, MenuId: mid})
	}
	return tx.Create(&records).Error
}

// CreateWithMenus 事务：创建角色 + 关联菜单
func (r *roleRepository) CreateWithMenus(ctx context.Context, d dto.CreateRoleDTO) (*model.Role, error) {
	var role model.Role
	utils.Copy(&role, d)
	err := r.Transaction(ctx, func(tx *gorm.DB) error {
		if err := tx.Create(&role).Error; err != nil {
			return err
		}
		return saveRoleMenus(tx, role.Id, d.MenuIds)
	})
	return &role, err
}

// UpdateWithMenus 事务：更新角色 + 重置关联菜单
func (r *roleRepository) UpdateWithMenus(ctx context.Context, id int64, d dto.UpdateRoleDTO) error {
	mp := utils.StructToMapIgnoreNil(d)
	delete(mp, "id")
	delete(mp, "menuIds")
	return r.Transaction(ctx, func(tx *gorm.DB) error {
		if len(mp) > 0 {
			if err := tx.Model(&model.Role{}).Where("id = ?", id).Updates(mp).Error; err != nil {
				return err
			}
		}
		return saveRoleMenus(tx, id, d.MenuIds)
	})
}
