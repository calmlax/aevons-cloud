/**
 * 岗位信息表 Model
 *
 * @author
 * @date 2026-04-18
 * @website https://www.aevons.com
 * @email 1073602@qq.com
 * @copyright 2025-2026 Aevons
 */
package model

import (
	"github.com/calmlax/aevons-framework/core/base"
)

type Post struct {

	//岗位编号
	Id int64 `gorm:"column:id;primaryKey;autoIncrement;type:bigint;comment:岗位编号" json:"id,string"`
	//岗位标识
	PostKey string `gorm:"column:post_key;type:varchar(32);comment:岗位标识" json:"postKey"`
	//岗位名称
	PostName string `gorm:"column:post_name;type:varchar(50);comment:岗位名称" json:"postName"`
	//顺序
	Sort int `gorm:"column:sort;type:int;comment:顺序" json:"sort"`
	//状态（0正常 1停用）
	Status int16 `gorm:"column:status;type:tinyint(1);comment:状态（0正常 1停用）" json:"status"`
	//备注
	Remark string `gorm:"column:remark;type:varchar(500);comment:备注" json:"remark"`
	base.DefaultModel
}

// TableName 指定表名
func (Post) TableName() string {
	return "sys_post"
}
