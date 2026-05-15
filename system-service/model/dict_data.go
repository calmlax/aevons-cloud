/**
 * 字典数据表 Model
 *
 * @author
 * @date 2026-04-09 01:08:50.443674979 +0000 UTC
 * @website https://www.aevons.com
 * @email 1073602@qq.com
 * @copyright 2025-2026 Aevons
 */
package model

import (
	"github.com/calmlax/aevons-framework/core/base"
)

type DictData struct {

	//字典编号
	Id int64 `gorm:"column:id;primaryKey;autoIncrement;type:bigint;comment:字典编号" json:"id,string"`

	//字典类型
	DictType string `gorm:"column:dict_type;type:varchar(32);comment:字典类型" json:"dictType"`

	//字典键值
	DictValue string `gorm:"column:dict_value;type:varchar(32);comment:字典键值" json:"dictValue"`

	//状态（0正常 1停用）
	Status int16 `gorm:"column:status;type:tinyint(1);comment:状态（0正常 1停用）" json:"status"`

	//顺序
	Sort int `gorm:"column:sort;type:int;comment:顺序" json:"sort"`

	//标签风格
	TagType string `gorm:"column:tag_type;type:varchar(10);comment:标签风格" json:"tagType"`

	//样式类名
	TagClass string `gorm:"column:tag_class;type:varchar(10);comment:样式类名" json:"tagClass"`

	base.DefaultModel
}

// TableName 指定表名
func (DictData) TableName() string {
	return "sys_dict_data"
}
