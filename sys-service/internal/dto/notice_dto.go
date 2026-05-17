package dto

/**
 * 通知公告 DTO
 *
 * @author
 * @date 2026-04-21
 * @website https://www.aevons.com
 * @email 1073602@qq.com
 * @copyright 2025-2026 Aevons
 */

import (
	"time"

	"github.com/calmlax/aevons-framework/core/base"
)

// q 支持查询操作: eq, ne, gt, gte, lt, lte, like, like_l, like_r, in, not_in, between, is_null, not_null
// gorm.column 必须对应真实数据库列名
type NoticeQuery struct {

	// 公告标题
	Title *string `form:"title" gorm:"column:title" q:"like"`
	// 公告类型（1通知 2公告）
	Type *int16 `form:"type" gorm:"column:type" q:"eq"`
	// 公告内容
	Content *string `form:"content" gorm:"column:content" q:"like"`
	// 状态（0正常 1关闭）
	Status *int16 `form:"status" gorm:"column:status" q:"eq"`
	base.BaseQuery
}

type CreateNoticeDTO struct {

	//公告标题
	Title string `json:"title" binding:"required,max=50"`
	//公告类型（1通知 2公告）
	Type int16 `json:"type" binding:"required"`
	//公告内容
	Content string `json:"content" binding:"required,max=3000"`
	//状态（0正常 1关闭）
	Status int16 `json:"status" binding:"required"`
	//备注
	Remark string `json:"remark" binding:"max=255"`
}

type UpdateNoticeDTO struct {
	//公告ID
	Id *int `json:"id" binding:"required"`
	//公告标题
	Title *string `json:"title" binding:"required,max=50"`
	//公告类型（1通知 2公告）
	Type *int16 `json:"type" binding:"required"`
	//公告内容
	Content *string `json:"content" binding:"required,max=3000"`
	//状态（0正常 1关闭）
	Status *int16 `json:"status" binding:"required"`
	//备注
	Remark *string `json:"remark" binding:"max=255"`
}

type NoticeDTO struct {
	//公告ID
	Id int `excel:"column:公告ID;index:0;dict:" json:"id" binding:"required"`
	//公告标题
	Title string `excel:"column:公告标题;index:1;dict:" json:"title" binding:"required,max=50"`
	//公告类型（1通知 2公告）
	Type int16 `excel:"column:公告类型（1通知 2公告）;index:2;dict:sys_notice_type" json:"type" binding:"required"`
	//公告内容
	Content string `excel:"column:公告内容;index:3;dict:" json:"content" binding:"required,max=3000"`
	//状态（0正常 1关闭）
	Status int16 `excel:"column:状态（0正常 1关闭）;index:4;dict:sys_common_status" json:"status" binding:"required"`
	//备注
	Remark string `excel:"column:备注;index:5;dict:" json:"remark" binding:"max=255"`
	//创建者
	CreatedBy int64 `excel:"column:创建者;index:6;dict:" json:"createdBy,string"`
	//创建时间
	CreatedAt time.Time `excel:"column:创建时间;index:7;dict:" json:"createdAt"`
	//更新者
	UpdatedBy int64 `excel:"column:更新者;index:8;dict:" json:"updatedBy,string"`
	//更新时间
	UpdatedAt time.Time `excel:"column:更新时间;index:9;dict:" json:"updatedAt"`
}
