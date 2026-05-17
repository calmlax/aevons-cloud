/**
 * 字典数据翻译 Model
 *
 * @author
 * @date 2026-04-12
 * @website https://www.aevons.com
 * @email 1073602@qq.com
 * @copyright 2025-2026 Aevons
 */
package model

type DictDataTl struct {

	//字典数据ID
	DictDataId int64 `gorm:"column:dict_data_id;primaryKey;type:bigint;comment:字典数据ID" json:"dictDataId,string"`

	//语言标识
	LangCode string `gorm:"column:lang_code;primaryKey;type:varchar(10);comment:语言标识" json:"langCode"`

	//标签翻译
	Label string `gorm:"column:label;type:varchar(100);comment:标签翻译" json:"label"`

	//提示翻译
	Tip string `gorm:"column:tip;type:varchar(200);comment:提示翻译" json:"tip"`
}

// TableName 指定表名
func (DictDataTl) TableName() string {
	return "sys_dict_data_tl"
}
