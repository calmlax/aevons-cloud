/**
 * 菜单权限表 DTO
 *
 * @author
 * @date 2026-04-09 01:17:32.070579992 +0000 UTC
 * @website https://www.aevons.com
 * @email 1073602@qq.com
 * @copyright 2025-2026 Aevons
 */
package dto

import (
	"time"
)

type MenuDTO struct {
	//菜单ID
	Id int64 `excel:"column:菜单ID;index:0;dict:" json:"id,string" binding:"required"`
	//父级菜单Id
	ParentId int64 `excel:"column:父级菜单Id;index:1;dict:" json:"parentId,string" binding:"required"`
	//语言标识
	LangCode string `excel:"column:语言标识;index:2;dict:" json:"langCode" binding:"required,max=10"`
	//菜单名称
	Title string `excel:"column:菜单名称;index:2;dict:" json:"title" binding:"required,max=50"`
	//类型（1目录 2菜单 3按钮）
	Type int16 `excel:"column:类型（1目录 2菜单 3按钮）;index:3;dict:" json:"type" binding:"required"`
	//顺序
	Sort int `excel:"column:顺序;index:4;dict:" json:"sort" binding:"required"`
	//路由地址
	Path string `excel:"column:路由地址;index:5;dict:" json:"path" binding:"max=100"`
	//组件路径
	Component string `excel:"column:组件路径;index:6;dict:" json:"component" binding:"max=100"`
	//路由参数
	Query string `excel:"column:路由参数;index:7;dict:" json:"query" binding:"max=255"`
	//是否可见（0隐藏 1显示）
	Visible int16 `excel:"column:是否可见（0隐藏 1显示）;index:8;dict:sys_is" json:"visible" binding:"required"`
	//状态（0正常 1停用）
	Status int16 `excel:"column:状态（0正常 1停用）;index:9;dict:sys_common_status" json:"status" binding:"required"`
	//是否为外链（0否 1是）
	IsFrame int16 `excel:"column:是否为外链（0否 1是）;index:10;dict:sys_is" json:"isFrame" binding:"required"`
	//权限标识
	Permission string `excel:"column:权限标识;index:11;dict:" json:"permission" binding:"max=32"`
	//图标
	Icon string `excel:"column:图标;index:12;dict:" json:"icon" binding:"max=64"`
	//激活ID
	ActiveId int64 `excel:"column:激活ID;index:14;dict:" json:"activeId,string"`
	//
	CreatedBy int64 `excel:"column:;index:15;dict:" json:"createdBy,string"`
	//
	CreatedAt time.Time `excel:"column:;index:16;dict:" json:"createdAt"`
	//
	UpdatedBy int64 `excel:"column:;index:17;dict:" json:"updatedBy,string"`
	//
	UpdatedAt time.Time `excel:"column:;index:18;dict:" json:"updatedAt"`
	//语言信息
	Translations map[string]string `json:"translations" binding:"required,min=1"`
}
