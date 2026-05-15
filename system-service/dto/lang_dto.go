/**
 * 语言 DTO
 *
 * @author
 * @date 2026-04-19
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
type LangQuery struct {

	// 语言名称（如简体中文、English）
	LangName *string `form:"langName" gorm:"column:lang_name" q:"like"`
	// 状态（0正常,1停用）
	Status *int16 `form:"status" gorm:"column:status" q:"eq"`
	base.BaseQuery
}

type CreateLangDTO struct {

	//语言编码（如zh-CN、en-US）
	LangCode string `json:"langCode" binding:"required,max=10"`
	//语言名称（如简体中文、English）
	LangName string `json:"langName" binding:"required,max=50"`
	//是否默认语言（0否，1是）
	IsDefault int16 `json:"isDefault" binding:"required"`
	//排序值（升序）
	Sort int `json:"sort" binding:"required"`
	//状态（0正常,1停用）
	Status int16 `json:"status" binding:"required"`
	//备注
	Remark string `json:"remark" binding:"max=200"`
}

type UpdateLangDTO struct {
	//编号
	Id *int64 `json:"id,string" binding:"required"`
	//语言编码（如zh-CN、en-US）
	LangCode *string `json:"langCode" binding:"required,max=10"`
	//语言名称（如简体中文、English）
	LangName *string `json:"langName" binding:"required,max=50"`
	//是否默认语言（0否，1是）
	IsDefault *int16 `json:"isDefault" binding:"required"`
	//排序值（升序）
	Sort *int `json:"sort" binding:"required"`
	//状态（0正常,1停用）
	Status *int16 `json:"status" binding:"required"`
	//备注
	Remark *string `json:"remark" binding:"max=200"`
}

type LangDTO struct {
	//编号
	Id int64 `excel:"column:编号;index:0;dict:" json:"id,string" binding:"required"`
	//语言编码（如zh-CN、en-US）
	LangCode string `excel:"column:语言编码（如zh-CN、en-US）;index:1;dict:" json:"langCode" binding:"required,max=10"`
	//语言名称（如简体中文、English）
	LangName string `excel:"column:语言名称（如简体中文、English）;index:2;dict:" json:"langName" binding:"required,max=50"`
	//是否默认语言（0否，1是）
	IsDefault int16 `excel:"column:是否默认语言（0否，1是）;index:3;dict:sys_is" json:"isDefault" binding:"required"`
	//排序值（升序）
	Sort int `excel:"column:排序值（升序）;index:4;dict:" json:"sort" binding:"required"`
	//状态（0正常,1停用）
	Status int16 `excel:"column:状态（0正常,1停用）;index:5;dict:sys_common_status" json:"status" binding:"required"`
	//备注
	Remark string `excel:"column:备注;index:6;dict:" json:"remark" binding:"max=200"`
	//创建时间
	CreatedAt time.Time `excel:"column:创建时间;index:7;dict:" json:"createdAt"`
	//更新时间
	UpdatedAt time.Time `excel:"column:更新时间;index:8;dict:" json:"updatedAt"`
}
