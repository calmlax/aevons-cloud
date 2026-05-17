/**
 * 代码生成表 DTO
 *
 * @author
 * @date 2026-04-08 03:51:33.212851159 +0000 UTC
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
type GenTableQuery struct {

	// 表名
	TableName string `form:"tableName" gorm:"column:table_name" q:"like"`
	// 解释
	TableComment string `form:"tableComment" gorm:"column:table_comment" q:"like"`
	// 模块名称
	ModuleName string `form:"moduleName" gorm:"column:module_name" q:"like"`
	base.BaseQuery
}

type CreateGenTableDTO struct {

	//表名
	TableName string `json:"tableName" binding:"required,max=50"`
	//解释
	TableComment string `json:"tableComment" binding:"max=50"`
	//类名
	ClassName string `json:"className" binding:"max=50"`
	//模块名称
	ModuleName string `json:"moduleName" binding:"max=50"`
	//作者
	Author string `json:"author" binding:"max=255"`
	//路由
	Router string `json:"router" binding:"max=50"`
	//上级菜单编号
	MenuId int64 `json:"menuId"`
	//权限标识
	Permission string `json:"permission" binding:"max=64"`
	//备注
	Remark string `json:"remark" binding:"max=255"`
}

type UpdateGenTableDTO struct {
	//编号
	Id *int64 `json:"id,string"`
	//表名
	TableName *string `json:"tableName" binding:"required,max=50"`
	//解释
	TableComment *string `json:"tableComment" binding:"max=50"`
	//类名
	ClassName *string `json:"className" binding:"max=50"`
	//模块名称
	ModuleName *string `json:"moduleName" binding:"max=50"`
	//作者
	Author *string `json:"author" binding:"max=255"`
	//路由
	Router *string `json:"router" binding:"max=50"`
	//上级菜单编号
	MenuId *int64 `json:"menuId,string"`
	//权限标识
	Permission *string `json:"permission" binding:"max=64"`
	//备注
	Remark *string `json:"remark" binding:"max=255"`
}

type CodeDTO struct {
	Name string `json:"name"`

	FileName string `json:"-"`

	SaveFileName string `json:"-"`

	Code string `json:"code"`

	FileType string `json:"fileType"`
}
