/**
 * 角色和菜单关联表 Model
 *
 * @author
 * @date 2026-04-09 01:49:40.929237188 +0000 UTC
 * @website https://www.aevons.com
 * @email 1073602@qq.com
 * @copyright 2025-2026 Aevons
 */
package model

type RoleMenu struct {

	//角色ID
	RoleId int64 `gorm:"column:role_id;primaryKey;type:bigint;comment:角色ID" json:"roleId,string"`

	//菜单ID
	MenuId int64 `gorm:"column:menu_id;primaryKey;type:bigint;comment:菜单ID" json:"menuId,string"`
}

// TableName 指定表名
func (RoleMenu) TableName() string {
	return "sys_role_menu"
}
