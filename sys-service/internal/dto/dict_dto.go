/**
 * 字典类型表 DTO
 *
 * @author
 * @date 2026-04-09 01:08:50.442643548 +0000 UTC
 * @website https://www.aevons.com
 * @email 1073602@qq.com
 * @copyright 2025-2026 Aevons
 */
package dto

import (
	"github.com/calmlax/aevons-framework/core/base"
)

// q 支持查询操作: eq, ne, gt, gte, lt, lte, like, like_l, like_r, in, not_in, between, is_null, not_null
// gorm.column 必须对应真实数据库列名
type DictQuery struct {

	// 字典类型
	DictType *string `form:"dictType" gorm:"column:dict_type" q:"like"`
	// 字典名称
	DictName *string `form:"dictName" gorm:"column:dict_name" q:"like"`
	// 状态（0正常 1停用）
	Status *int16 `form:"status" gorm:"column:status" q:"eq"`
	// 系统内置（0否 1是）
	IsSys *int16 `form:"isSys" gorm:"column:is_sys" q:"eq"`
	base.BaseQuery
}

type CreateDictDTO struct {

	//字典类型
	DictType string `json:"dictType" binding:"max=32"`
	//字典名称
	DictName string `json:"dictName" binding:"max=50"`
	//状态（0正常 1停用）
	Status int16 `json:"status"`
	//系统内置（0否 1是）
	IsSys int16 `json:"isSys"`
	//备注
	Remark string `json:"remark" binding:"max=500"`
}

type UpdateDictDTO struct {
	//字典主键
	Id *int64 `json:"id,string" binding:"required"`
	//字典类型
	DictType *string `json:"dictType" binding:"max=32"`
	//字典名称
	DictName *string `json:"dictName" binding:"max=50"`
	//状态（0正常 1停用）
	Status *int16 `json:"status"`
	//系统内置（0否 1是）
	IsSys *int16 `json:"isSys"`
	//备注
	Remark *string `json:"remark" binding:"max=500"`
}
