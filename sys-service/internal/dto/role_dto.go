/**
 * 角色信息表 DTO
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
type RoleQuery struct {

	// 角色名称
	RoleName *string `form:"roleName" gorm:"column:role_name" q:"like"`
	// 角色权限字符串
	RoleKey *string `form:"roleKey" gorm:"column:role_key" q:"like"`
	// 数据范围（0：全部 1：自定义 2：本部门 3：本部门及以下）
	DataScope *int16 `form:"dataScope" gorm:"column:data_scope" q:"eq"`
	// 菜单树选择项是否父子联动（0否 1是）
	MenuCheckStrictly *int16 `form:"menuCheckStrictly" gorm:"column:menu_check_strictly" q:"eq"`
	// 部门树选择项是否父子联动（0否 1是）
	DeptCheckStrictly *int16 `form:"deptCheckStrictly" gorm:"column:dept_check_strictly" q:"eq"`
	// 状态（0正常 1停用）
	Status *int16 `form:"status" gorm:"column:status" q:"eq"`
	base.BaseQuery
}

type CreateRoleDTO struct {

	//角色名称
	RoleName string `json:"roleName" binding:"required,max=50"`
	//角色权限字符串
	RoleKey string `json:"roleKey" binding:"max=32"`
	//显示顺序
	Sort int `json:"sort" binding:"required"`
	//数据范围
	DataScope int16 `json:"dataScope"`
	//菜单树选择项是否父子联动（0否 1是）
	MenuCheckStrictly int16 `json:"menuCheckStrictly"`
	//部门树选择项是否父子联动（0否 1是）
	DeptCheckStrictly int16 `json:"deptCheckStrictly"`
	//状态（0正常 1停用）
	Status int16 `json:"status"`
	//备注
	Remark string `json:"remark" binding:"max=500"`
	//部门ID列表
	DeptIds []int64 `json:"deptIds"`
	//菜单ID列表
	MenuIds []int64 `json:"menuIds"`
}

type UpdateRoleDTO struct {
	//角色ID
	Id *int64 `json:"id,string" binding:"required"`
	//角色名称
	RoleName *string `json:"roleName" binding:"required,max=50"`
	//角色权限字符串
	RoleKey *string `json:"roleKey" binding:"max=32"`
	//显示顺序
	Sort *int `json:"sort" binding:"required"`
	//数据范围
	DataScope *int16 `json:"dataScope"`
	//菜单树选择项是否父子联动（0否 1是）
	MenuCheckStrictly *int16 `json:"menuCheckStrictly"`
	//部门树选择项是否父子联动（0否 1是）
	DeptCheckStrictly *int16 `json:"deptCheckStrictly"`
	//状态（0正常 1停用）
	Status *int16 `json:"status"`
	//备注
	Remark *string `json:"remark" binding:"max=500"`
	//部门ID列表
	DeptIds []int64 `json:"deptIds"`
	//菜单ID列表
	MenuIds []int64 `json:"menuIds"`
}

type RoleDTO struct {
	//角色ID
	Id int64 `excel:"column:角色ID;index:0;dict:" json:"id,string" binding:"required"`
	//角色名称
	RoleName string `excel:"column:角色名称;index:1;dict:" json:"roleName" binding:"required,max=50"`
	//角色权限字符串
	RoleKey string `excel:"column:角色权限字符串;index:2;dict:" json:"roleKey" binding:"max=32"`
	//显示顺序
	Sort int `excel:"column:显示顺序;index:3;dict:" json:"sort,string" binding:"required"`
	//数据范围（0：全部 1：自定义 2：本部门 3：本部门及以下）
	DataScope int16 `excel:"column:数据范围（0：全部 1：自定义 2：本部门 3：本部门及以下）;index:4;dict:" json:"dataScope" binding:"required"`
	//菜单树选择项是否父子联动（0否 1是）
	MenuCheckStrictly int16 `json:"menuCheckStrictly"`
	//部门树选择项是否父子联动（0否 1是）
	DeptCheckStrictly int16 `json:"deptCheckStrictly"`
	//状态（0正常 1停用）
	Status int16 `excel:"column:状态（0正常 1停用）;index:5;dict:" json:"status" binding:"required"`
	//备注
	Remark string `excel:"column:备注;index:6;dict:" json:"remark" binding:"max=500"`
	//
	CreatedBy int64 `excel:"column:;index:7;dict:" json:"createdBy,string"`
	//
	CreatedAt time.Time `excel:"column:;index:8;dict:" json:"createdAt"`
	//
	UpdatedBy int64 `excel:"column:;index:9;dict:" json:"updatedBy,string"`
	//
	UpdatedAt time.Time `excel:"column:;index:10;dict:" json:"updatedAt"`
}
