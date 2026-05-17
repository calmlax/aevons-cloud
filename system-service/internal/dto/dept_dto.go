/**
 * 部门表 DTO
 *
 * @author
 * @date 2026-04-18
 * @website https://www.aevons.com
 * @email 1073602@qq.com
 * @copyright 2025-2026 Aevons
 */
package dto

import (
	"time"

	"github.com/calmlax/aevons-framework/core/base"
)

// q 支持查询操作: eq, ne, gt, gte, lt, lte, like, like_l, like_r, in, not_in, between, is_null, not_null
// gorm.column 必须对应真实数据库列名
type DeptQuery struct {

	// 父级编号
	ParentId *int64 `form:"parentId" gorm:"column:parent_id" q:"eq"`
	// 祖级
	Ancestors *string `form:"ancestors" gorm:"column:ancestors" q:"like"`
	// 类型（1机构，2部门）
	DeptType *int16 `form:"deptType" gorm:"column:dept_type" q:"eq"`
	// 部门名称
	DeptName *string `form:"deptName" gorm:"column:dept_name" q:"like"`
	// 状态（0正常 1停用）
	Status *int16 `form:"status" gorm:"column:status" q:"eq"`
	base.BaseQuery
}

type CreateDeptDTO struct {

	//父级编号
	ParentId int64 `json:"parentId,string" binding:"required"`

	//类型（1机构，2部门）
	DeptType int16 `json:"deptType" binding:"required"`
	//部门名称
	DeptName string `json:"deptName" binding:"required,max=30"`
	//顺序
	Sort int `json:"sort" binding:"required"`
	//状态（0正常 1停用）
	Status int16 `json:"status" binding:"required"`
	//备注
	Remark string `json:"remark" binding:"max=255"`
}

type UpdateDeptDTO struct {
	//部门编号
	Id *int64 `json:"id,string" binding:"required"`
	//父级编号
	ParentId *int64 `json:"parentId,string" binding:"required"`

	//类型（1机构，2部门）
	DeptType *int16 `json:"deptType" binding:"required"`
	//部门名称
	DeptName *string `json:"deptName" binding:"required,max=30"`
	//顺序
	Sort *int `json:"sort" binding:"required"`
	//状态（0正常 1停用）
	Status *int16 `json:"status" binding:"required"`
	//备注
	Remark *string `json:"remark" binding:"max=255"`
}

type DeptDTO struct {
	//部门编号
	Id int64 `excel:"column:部门编号;index:0;dict:" json:"id,string" binding:"required"`
	//父级编号
	ParentId int64 `excel:"column:父级编号;index:1;dict:" json:"parentId,string" binding:"required"`
	//祖级
	Ancestors string `excel:"column:祖级;index:2;dict:" json:"ancestors" binding:"max=255"`
	//类型（1机构，2部门）
	DeptType int16 `excel:"column:类型（1机构，2部门）;index:3;dict:sys_dept_type" json:"deptType" binding:"required"`
	//部门名称
	DeptName string `excel:"column:部门名称;index:4;dict:" json:"deptName" binding:"required,max=30"`
	//顺序
	Sort int `excel:"column:顺序;index:5;dict:" json:"sort" binding:"required"`
	//状态（0正常 1停用）
	Status int16 `excel:"column:状态（0正常 1停用）;index:6;dict:sys_common_status" json:"status" binding:"required"`
	//备注
	Remark string `excel:"column:备注;index:7;dict:" json:"remark" binding:"max=255"`
	//
	CreatedBy int64 `excel:"column:;index:8;dict:" json:"createdBy,string"`
	//
	CreatedAt time.Time `excel:"column:;index:9;dict:" json:"createdAt"`
	//
	UpdatedBy int64 `excel:"column:;index:10;dict:" json:"updatedBy,string"`
	//
	UpdatedAt time.Time `excel:"column:;index:11;dict:" json:"updatedAt"`
}
