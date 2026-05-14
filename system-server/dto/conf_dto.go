/**
 * 参数配置表 DTO
 *
 * @author
 * @date 2026-04-09 00:38:25.504785055 +0000 UTC
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
type ConfQuery struct {

	// 编号
	Id *int64 `form:"id" gorm:"column:id" q:"eq"`
	// 名称
	Name *string `form:"name" gorm:"column:name" q:"like"`
	// 配置KEY
	ConfKey *string `form:"confKey" gorm:"column:conf_key" q:"eq"`
	// 配置值
	ConfValue *string `form:"confValue" gorm:"column:conf_value" q:"eq"`
	// 是否系统内置（0否 1是）
	IsSys *int16 `form:"isSys" gorm:"column:is_sys" q:"eq"`
	// 范围（0公开配置，1登录配置，2后台服务配置）
	Scope *int16 `form:"scope" gorm:"column:scope" q:"eq"`
	// 是否加密/脱敏（0否 1是）
	IsSecret *int16 `form:"isSecret" gorm:"column:is_secret" q:"eq"`
	base.BaseQuery
}

type CreateConfDTO struct {

	//名称
	Name string `json:"name" binding:"required,max=255"`
	//配置KEY
	ConfKey string `json:"confKey" binding:"required,max=64"`
	//配置值
	ConfValue string `json:"confValue" binding:"required"`
	//是否系统内置（0否 1是）
	IsSys int16 `json:"isSys"`
	//范围（0公开配置，1登录配置，2后台服务配置）
	Scope int16 `json:"scope"`
	//是否加密/脱敏（0否 1是）
	IsSecret int16 `json:"isSecret"`
	//备注
	Remark string `json:"remark" binding:"max=500"`
}

type UpdateConfDTO struct {
	//编号
	Id *int64 `json:"id" binding:"required"`
	//名称
	Name *string `json:"name" binding:"required,max=255"`
	//配置KEY
	ConfKey *string `json:"confKey" binding:"required,max=64"`
	//配置值
	ConfValue *string `json:"confValue" binding:"required"`
	//是否系统内置（0否 1是）
	IsSys *int16 `json:"isSys"`
	//范围（0公开配置，1登录配置，2后台服务配置）
	Scope *int16 `json:"scope"`
	//是否加密/脱敏（0否 1是）
	IsSecret *int16 `json:"isSecret"`
	//备注
	Remark *string `json:"remark" binding:"max=500"`
}
