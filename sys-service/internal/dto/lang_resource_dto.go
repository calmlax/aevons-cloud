/**
 * 语言资源 DTO
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
type LangResourceQuery struct {

	// 资源标识
	ResourceKey *string `form:"resourceKey" gorm:"column:resource_key" q:"eq"`
	// 命名空间（default）
	Namespace *string `form:"namespace" gorm:"column:namespace" q:"like"`
	// 语言编码（如zh）
	LangCode *string `form:"langCode" gorm:"column:lang_code" q:"eq"`
	// 内容
	Content *string `form:"content" gorm:"column:content" q:"like"`
	// 状态（0正常,1停用）
	Status *int16 `form:"status" gorm:"column:status" q:"eq"`
	base.BaseQuery
}

type CreateLangResourceDTO struct {

	//资源标识
	ResourceKey string `json:"resourceKey" binding:"required,max=128"`
	//命名空间（default）
	Namespace string `json:"namespace" binding:"required,max=32"`
	//语言编码（如zh）
	LangCode string `json:"langCode" binding:"required,max=10"`
	//内容
	Content string `json:"content" binding:"required,max=500"`
	//状态（0正常,1停用）
	Status int16 `json:"status" binding:"required"`
}

type UpdateLangResourceDTO struct {
	//主键ID
	Id *int64 `json:"id,string" binding:"required"`
	//资源标识
	ResourceKey *string `json:"resourceKey" binding:"required,max=128"`
	//命名空间（default）
	Namespace *string `json:"namespace" binding:"required,max=32"`
	//语言编码（如zh）
	LangCode *string `json:"langCode" binding:"required,max=10"`
	//内容
	Content *string `json:"content" binding:"required,max=500"`
	//状态（0正常,1停用）
	Status *int16 `json:"status" binding:"required"`
}

// TranslationItem 单条翻译
type TranslationItem struct {
	LangCode string `json:"langCode"`
	Content  string `json:"content"`
}

// SaveTranslationsDTO 批量保存翻译请求
type SaveTranslationsDTO struct {
	Namespace   string            `json:"namespace" binding:"required"`
	ResourceKey string            `json:"resourceKey" binding:"required"`
	Items       []TranslationItem `json:"items" binding:"required"`
}

type LangResourceDTO struct {
	//主键ID
	Id int64 `excel:"column:主键ID;index:0;dict:" json:"id,string" binding:"required"`
	//资源标识
	ResourceKey string `excel:"column:资源标识;index:1;dict:" json:"resourceKey" binding:"required,max=128"`
	//命名空间（default）
	Namespace string `excel:"column:命名空间（default）;index:2;dict:" json:"namespace" binding:"required,max=32"`
	//语言编码（如zh）
	LangCode string `excel:"column:语言编码（如zh）;index:3;dict:" json:"langCode" binding:"required,max=10"`
	//内容
	Content string `excel:"column:内容;index:4;dict:" json:"content" binding:"required,max=500"`
	//状态（0正常,1停用）
	Status int16 `excel:"column:状态（0正常,1停用）;index:5;dict:sys_common_status" json:"status" binding:"required"`
	//创建时间
	CreatedAt time.Time `excel:"column:创建时间;index:6;dict:" json:"createdAt"`
	//更新时间
	UpdatedAt time.Time `excel:"column:更新时间;index:7;dict:" json:"updatedAt"`
}
