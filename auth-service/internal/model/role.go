/**
 * 角色信息表 Model
 *
 * @author
 * @date 2026-04-09 01:49:40.926495422 +0000 UTC
 * @website https://www.aevons.com
 * @email 1073602@qq.com
 * @copyright 2025-2026 Aevons
 */
package model

import (
	"github.com/calmlax/aevons-framework/core/base"
)

type Role struct {

	//角色ID
	Id int64 `gorm:"column:id;primaryKey;autoIncrement;type:bigint;comment:角色ID" json:"id,string"`

	//角色名称
	RoleName string `gorm:"column:role_name;type:varchar(50);comment:角色名称" json:"roleName"`

	//角色权限字符串
	RoleKey string `gorm:"column:role_key;type:varchar(32);comment:角色权限字符串" json:"roleKey"`

	//显示顺序
	Sort int `gorm:"column:sort;type:int;comment:显示顺序" json:"sort"`

	//数据范围（0：全部 1：自定义 2：本部门 3：本部门及以下）
	DataScope int16 `gorm:"column:data_scope;type:tinyint;comment:数据范围" json:"dataScope"`

	//状态（0正常 1停用）
	Status int16 `gorm:"column:status;type:tinyint(1);comment:状态（0正常 1停用）" json:"status"`

	//备注
	Remark string `gorm:"column:remark;type:varchar(500);comment:备注" json:"remark"`

	base.DefaultModel
}

// TableName 指定表名
func (Role) TableName() string {
	return "sys_role"
}
