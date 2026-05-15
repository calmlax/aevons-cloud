/**
 * 岗位信息表 DTO
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
type PostQuery struct {

	// 岗位名称
	PostName *string `form:"postName" gorm:"column:post_name" q:"like"`
	// 状态（0正常 1停用）
	Status *int16 `form:"status" gorm:"column:status" q:"eq"`
	base.BaseQuery
}

type CreatePostDTO struct {

	//岗位标识
	PostKey string `json:"postKey" binding:"required,max=32"`
	//岗位名称
	PostName string `json:"postName" binding:"required,max=50"`
	//顺序
	Sort int `json:"sort" binding:"required"`
	//状态（0正常 1停用）
	Status int16 `json:"status" binding:"required"`
	//备注
	Remark string `json:"remark" binding:"max=500"`
}

type UpdatePostDTO struct {
	//岗位编号
	Id *int64 `json:"id,string" binding:"required"`
	//岗位标识
	PostKey *string `json:"postKey" binding:"required,max=32"`
	//岗位名称
	PostName *string `json:"postName" binding:"required,max=50"`
	//顺序
	Sort *int `json:"sort" binding:"required"`
	//状态（0正常 1停用）
	Status *int16 `json:"status" binding:"required"`
	//备注
	Remark *string `json:"remark" binding:"max=500"`
}

type PostDTO struct {
	//岗位编号
	Id int64 `excel:"column:岗位编号;index:0;dict:" json:"id,string" binding:"required"`
	//岗位标识
	PostKey string `excel:"column:岗位标识;index:1;dict:" json:"postKey" binding:"required,max=32"`
	//岗位名称
	PostName string `excel:"column:岗位名称;index:2;dict:" json:"postName" binding:"required,max=50"`
	//顺序
	Sort int `excel:"column:顺序;index:3;dict:" json:"sort" binding:"required"`
	//状态（0正常 1停用）
	Status int16 `excel:"column:状态（0正常 1停用）;index:5;dict:sys_common_status" json:"status" binding:"required"`
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
