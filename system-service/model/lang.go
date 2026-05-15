/**
 * 语言 Model
 *
 * @author
 * @date 2026-04-19
 * @website https://www.aevons.com
 * @email 1073602@qq.com
 * @copyright 2025-2026 Aevons
 */
package model

import (
	"github.com/calmlax/aevons-framework/core/base"
)

type Lang struct {

	//编号
	Id int64 `gorm:"column:id;primaryKey;autoIncrement;type:bigint;comment:编号" json:"id,string"`
	//语言编码（如zh-CN、en-US）
	LangCode string `gorm:"column:lang_code;type:varchar(10);comment:语言编码（如zh-CN、en-US）" json:"langCode"`
	//语言名称（如简体中文、English）
	LangName string `gorm:"column:lang_name;type:varchar(50);comment:语言名称（如简体中文、English）" json:"langName"`
	//是否默认语言（0否，1是）
	IsDefault int16 `gorm:"column:is_default;type:tinyint(1);comment:是否默认语言（0否，1是）" json:"isDefault"`
	//排序值（升序）
	Sort int `gorm:"column:sort;type:int;comment:排序值（升序）" json:"sort"`
	//状态（0正常,1停用）
	Status int16 `gorm:"column:status;type:tinyint(1);comment:状态（0正常,1停用）" json:"status"`
	//备注
	Remark string `gorm:"column:remark;type:varchar(200);comment:备注" json:"remark"`
	base.BaseModel
}

// TableName 指定表名
func (Lang) TableName() string {
	return "sys_lang"
}
