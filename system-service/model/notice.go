package model

/**
 * 通知公告 Model
 *
 * @author
 * @date 2026-04-21
 * @website https://www.aevons.com
 * @email 1073602@qq.com
 * @copyright 2025-2026 Aevons
 */

import (
	"github.com/calmlax/aevons-framework/core/base"
)

type Notice struct {

	//公告ID
	Id int `gorm:"column:id;primaryKey;autoIncrement;type:int;comment:公告ID" json:"id"`
	//公告标题
	Title string `gorm:"column:title;type:varchar(50);comment:公告标题" json:"title"`
	//公告类型（1通知 2公告）
	Type int16 `gorm:"column:type;type:tinyint(1);comment:公告类型（1通知 2公告）" json:"type"`
	//公告内容
	Content string `gorm:"column:content;type:varchar(3000);comment:公告内容" json:"content"`
	//状态（0正常 1关闭）
	Status int16 `gorm:"column:status;type:tinyint(1);comment:状态（0正常 1关闭）" json:"status"`
	//备注
	Remark string `gorm:"column:remark;type:varchar(255);comment:备注" json:"remark"`
	base.DefaultModel
}

// TableName 指定表名
func (Notice) TableName() string {
	return "sys_notice"
}
