/**
 * 字典类型表 Model
 *
 * @author
 * @date 2026-04-09 01:08:50.442643548 +0000 UTC
 * @website https://www.aevons.com
 * @email 1073602@qq.com
 * @copyright 2025-2026 Aevons
 */
package model

import (
	"github.com/calmlax/aevons-framework/core/base"
)

type Dict struct {

	//字典主键
	Id int64 `gorm:"column:id;primaryKey;autoIncrement;type:bigint;comment:字典主键" json:"id,string"`

	//字典类型
	DictType string `gorm:"column:dict_type;type:varchar(32);comment:字典类型" json:"dictType"`

	//字典名称
	DictName string `gorm:"column:dict_name;type:varchar(50);comment:字典名称" json:"dictName"`

	//状态（0正常 1停用）
	Status int16 `gorm:"column:status;type:tinyint(1);comment:状态（0正常 1停用）" json:"status"`

	//系统内置（0否 1是）
	IsSys int16 `gorm:"column:is_sys;type:tinyint(1);comment:系统内置（0否 1是）" json:"isSys"`

	//备注
	Remark string `gorm:"column:remark;type:varchar(500);comment:备注" json:"remark"`

	base.DefaultModel
}

// TableName 指定表名
func (Dict) TableName() string {
	return "sys_dict"
}
