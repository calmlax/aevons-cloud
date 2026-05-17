/**
 * 代码生成表 Model
 *
 * @author
 * @date 2026-04-08 03:51:33.212851159 +0000 UTC
 * @website https://www.aevons.com
 * @email 1073602@qq.com
 * @copyright 2025-2026 Aevons
 */
package model

import (
	"github.com/calmlax/aevons-framework/core/base"
)

type GenTable struct {

	//编号
	Id int64 `gorm:"column:id;primaryKey;autoIncrement;type:bigint;comment:编号" json:"id,string"`

	//表名
	DbTableName string `gorm:"column:table_name;type:varchar(50);comment:表名" json:"tableName"`

	//解释
	DbTableComment string `gorm:"column:table_comment;type:varchar(50);comment:解释" json:"tableComment"`

	//类名
	ClassName string `gorm:"column:class_name;type:varchar(50);comment:类名" json:"className"`

	//模块名称
	ModuleName string `gorm:"column:module_name;type:varchar(50);comment:模块名称" json:"moduleName"`

	//作者
	Author string `gorm:"column:author;type:varchar(255);comment:作者" json:"author"`

	//路由
	Router string `gorm:"column:router;type:varchar(50);comment:路由" json:"router"`

	//上级菜单编号
	MenuId int64 `gorm:"column:menu_id;type:bigint;comment:上级菜单编号" json:"menuId,string"`

	//权限标识
	Permission string `gorm:"column:permission;type:varchar(64);comment:权限标识" json:"permission"`

	//备注
	Remark string `gorm:"column:remark;type:varchar(255);comment:备注" json:"remark"`

	base.BaseModel
}

// TableName 指定表名
func (GenTable) TableName() string {
	return "gen_table"
}
