/**
 * 字典数据表 DTO
 *
 * @author
 * @date 2026-04-09 01:08:50.443674979 +0000 UTC
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
type DictDataQuery struct {

	// 字典编号
	Id *int64 `form:"id" gorm:"column:id" q:"eq"`
	// 字典类型
	DictType *string `form:"dictType" gorm:"column:dict_type" q:"eq"`
	base.BaseQuery
}

type CreateDictDataDTO struct {

	//字典类型
	DictType string `json:"dictType" binding:"required,max=32"`
	//字典键值
	DictValue string `json:"dictValue" binding:"required,max=32"`
	//状态（0正常 1停用）
	Status int16 `json:"status"`
	//顺序
	Sort int `json:"sort" binding:"required"`
	//标签风格
	TagType string `json:"tagType" binding:"max=10"`
	//样式类名
	TagClass string `json:"tagClass" binding:"max=50"`
	//语言字典信息
	Translations map[string]DictDataTlDTO `json:"translations" binding:"required,min=1"`
}

type UpdateDictDataDTO struct {
	//字典编号
	Id *int64 `json:"id,string" binding:"required"`
	//字典类型
	DictType *string `json:"dictType" binding:"required,max=32"`
	//字典键值
	DictValue *string `json:"dictValue" binding:"required,max=32"`
	//状态（0正常 1停用）
	Status *int16 `json:"status"`
	//顺序
	Sort *int `json:"sort" binding:"required"`
	//标签风格
	TagType *string `json:"tagType" binding:"max=10"`
	//样式类名
	TagClass string `json:"tagClass" binding:"max=50"`
	//语言字典信息
	Translations map[string]DictDataTlDTO `json:"translations" binding:"required,min=1"`
}

type DictDataDTO struct {
	//字典编号
	Id int64 `excel:"column:字典编号;index:0;dict:" json:"id,string" binding:"required"`
	//字典类型
	DictType string `excel:"column:字典类型;index:1;dict:" json:"dictType" binding:"max=32"`
	//字典键值
	DictValue string `excel:"column:字典键值;index:2;dict:" json:"dictValue" binding:"max=32"`
	//状态（0正常 1停用）
	Status int16 `excel:"column:状态;index:3;dict:" json:"status"`
	//顺序
	Sort int `excel:"column:顺序;index:5;dict:" json:"sort,string"`
	//标签风格
	TagType string `excel:"column:标签风格;index:6;dict:" json:"tagType" binding:"max=10"`
	//样式类名
	TagClass string `excel:"column:标签风格;index:6;dict:" json:"tagClass" binding:"max=50"`
	//语言标识
	LangCode string `excel:"column:语言标识;index:9;dict:" json:"langCode" binding:"required,max=10"`
	//标签翻译
	Label string `excel:"column:标签翻译;index:10;dict:" json:"label" binding:"required,max=100"`
	//提示翻译
	Tip string `excel:"column:提示翻译;index:11;dict:" json:"tip" binding:"max=200"`
	//语言字典信息
	Translations map[string]DictDataTlDTO `excel:"column:翻译;index:12;dict:" json:"translations,omitempty"`
}

type DictDataTlDTO struct {
	//标签翻译
	Label string `json:"label" binding:"required,max=100"`
	//提示翻译
	Tip string `json:"tip" binding:"max=200"`
}

type SortItemDTO struct {
	Id   int64 `json:"id,string" binding:"required"`
	Sort int   `json:"sort"`
}
