/**
 * 菜单权限表 Model
 *
 * @author
 * @date 2026-04-18
 * @website https://www.aevons.com
 * @email 1073602@qq.com
 * @copyright 2025-2026 Aevons
 */
package model

import (
	"github.com/calmlax/aevons-framework/core/base"
)

type Menu struct {

	//菜单ID
	Id int64 `gorm:"column:id;primaryKey;autoIncrement;type:bigint;comment:菜单ID" json:"id,string"`
	//父级菜单Id
	ParentId int64 `gorm:"column:parent_id;type:bigint;comment:父级菜单Id" json:"parentId,string"`
	//类型（1目录 2菜单 3按钮）
	Type int16 `gorm:"column:type;type:tinyint(1);comment:类型（1目录 2菜单 3按钮）" json:"type"`
	//顺序
	Sort int `gorm:"column:sort;type:int;comment:顺序" json:"sort"`
	//路由地址
	Path string `gorm:"column:path;type:varchar(100);comment:路由地址" json:"path"`
	//组件路径
	Component string `gorm:"column:component;type:varchar(100);comment:组件路径" json:"component"`
	//路由参数
	Query string `gorm:"column:query;type:varchar(255);comment:路由参数" json:"query"`
	//是否可见（0隐藏 1显示）
	Visible int16 `gorm:"column:visible;type:tinyint(1);comment:是否可见（0隐藏 1显示）" json:"visible"`
	//状态（0正常 1停用）
	Status int16 `gorm:"column:status;type:tinyint(1);comment:状态（0正常 1停用）" json:"status"`
	//是否为外链（0否 1是）
	IsFrame int16 `gorm:"column:is_frame;type:tinyint(1);comment:是否为外链（0否 1是）" json:"isFrame"`
	//权限标识
	Permission string `gorm:"column:permission;type:varchar(32);comment:权限标识" json:"permission"`
	//图标
	Icon string `gorm:"column:icon;type:varchar(64);comment:图标" json:"icon"`
	//激活ID
	ActiveId int64 `gorm:"column:active_id;type:bigint;comment:激活ID" json:"activeId,string"`
	base.DefaultModel
}

// TableName 指定表名
func (Menu) TableName() string {
	return "sys_menu"
}
