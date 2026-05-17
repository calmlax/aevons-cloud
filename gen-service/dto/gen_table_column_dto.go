/**
 * 代码生成表字段 DTO
 *
 * @author
 * @date 2026-04-08 03:51:33.21431843 +0000 UTC
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
type GenTableColumnQuery struct {

	// 所属表编号
	TableId int64 `form:"tableId" gorm:"column:table_id" q:"eq"`
	// DB字段名称
	ColumnName string `form:"columnName" gorm:"column:column_name" q:"like"`
	// DB字段解释
	ColumnComment string `form:"columnComment" gorm:"column:column_comment" q:"eq"`
	// 字段名称
	FieldName string `form:"fieldName" gorm:"column:field_name" q:"like"`

	base.BaseQuery
}

type CreateGenTableColumnDTO struct {

	//所属表编号
	TableId int64 `json:"tableId" binding:"required"`
	//DB字段名称
	ColumnName string `json:"columnName" binding:"required,max=50"`
	//DB字段解释
	ColumnComment string `json:"columnComment" binding:"max=50"`
	//DB字段数据类型
	ColumnType string `json:"columnType" binding:"required,max=20"`
	//数据类型
	DataType string `json:"dataType" binding:"required,max=20"`
	//字段名称
	FieldName string `json:"fieldName" binding:"required,max=50"`
	//JSON
	Json string `json:"json" binding:"required,max=100"`
	//是否主键（0否，1是）
	IsPrimaryKey bool `json:"isPrimaryKey"`
	//是否自增（0否，1是）
	IsAutoIncrement bool `json:"isAutoIncrement"`
	//是否必填（0否，1是）
	IsRequired bool `json:"isRequired"`
	//是否插入字段（0否，1是）
	IsInsert bool `json:"isInsert"`
	//是否编辑字段（0否，1是）
	IsEdit bool `json:"isEdit"`
	//是否列表字段（0否，1是）
	IsList bool `json:"isList"`
	//列表排序字段（0否，1是）
	Sortable bool `json:"sortable"`
	//列表筛选字段（0否，1是）
	Filterable bool `json:"filterable"`
	//查询条件
	Condition string `json:"condition" binding:"max=10"`
	//字典类型
	DictType string `json:"dictType" binding:"max=64"`
	//排序
	Sort int `json:"sort" binding:"required"`
	//UI组件
	Component string `json:"component" binding:"max=20"`
	//默认值
	DefaultValue string `json:"defaultValue" binding:"max=50"`
	//数据长款
	DataLength int `json:"dataLength"`
	//数据精度
	DataPrecision bool `json:"dataPrecision"`
}

type UpdateGenTableColumnDTO struct {
	ID            *int64  `json:"id,string" binding:"required" label:"编号"`       // 主键ID 👈 改成 string 就好了
	ColumnComment *string `json:"columnComment" binding:"max=50" label:"DB字段注释"` // DB字段注释
	DataType      *string `json:"dataType" binding:"max=20" label:"数据类型"`        // 数据类型
	FieldName     *string `json:"fieldName" binding:"max=50" label:"实体字段名"`      // 实体字段名
	Json          *string `json:"json" binding:"max=100" label:"JSON配置"`         // JSON
	IsRequired    *bool   `json:"isRequired" label:"是否必填"`
	IsInsert      *bool   `json:"isInsert" label:"是否插入字段"`                   // 是否插入字段
	IsEdit        *bool   `json:"isEdit" label:"是否编辑字段"`                     // 是否编辑字段
	IsList        *bool   `json:"isList" label:"是否列表字段"`                     // 是否列表字段（0否，1是）
	Sortable      *bool   `json:"sortable" label:"列表排序字段"`                   // 列表排序字段（0否，1是）
	Filterable    *bool   `json:"filterable" label:"列表筛选字段"`                 // 列表筛选字段（0否，1是）
	Condition     *string `json:"condition" binding:"max=10" label:"查询方式"`   // 查询条件
	DictType      *string `json:"dictType" binding:"max=64" label:"字典类型"`    // 字典类型
	Sort          *int    `json:"sort" label:"排序"`                           // 排序
	Component     *string `json:"component" binding:"max=20" label:"UI组件"`   // UI组件
	DefaultValue  *string `json:"defaultValue" binding:"max=50" label:"默认值"` // 默认值
}
