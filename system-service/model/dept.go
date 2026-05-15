/**
 * 部门表 Model
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

type Dept struct {

	//部门编号
	Id int64 `gorm:"column:id;primaryKey;autoIncrement;type:bigint;comment:部门编号" json:"id,string"`
	//父级编号
	ParentId int64 `gorm:"column:parent_id;type:bigint;comment:父级编号" json:"parentId,string"`
	//祖级
	Ancestors string `gorm:"column:ancestors;type:varchar(255);comment:祖级" json:"ancestors"`
	//类型（1机构，2部门）
	DeptType int16 `gorm:"column:dept_type;type:tinyint;comment:类型（1机构，2部门）" json:"deptType"`
	//部门名称
	DeptName string `gorm:"column:dept_name;type:varchar(30);comment:部门名称" json:"deptName"`
	//顺序
	Sort int `gorm:"column:sort;type:int;comment:顺序" json:"sort"`
	//状态（0正常 1停用）
	Status int16 `gorm:"column:status;type:tinyint(1);comment:状态（0正常 1停用）" json:"status"`
	//备注
	Remark string `gorm:"column:remark;type:varchar(255);comment:备注" json:"remark"`
	base.DefaultModel
}

// TableName 指定表名
func (Dept) TableName() string {
	return "sys_dept"
}
