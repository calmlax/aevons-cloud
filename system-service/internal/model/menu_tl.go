/**
 * 菜单翻译 Model
 *
 * @author
 * @date 2026-04-12
 * @website https://www.aevons.com
 * @email 1073602@qq.com
 * @copyright 2025-2026 Aevons
 */
package model

type MenuTl struct {

	//菜单Id
	MenuId int64 `gorm:"column:menu_id;primaryKey;type:bigint;comment:菜单Id" json:"menuId,string"`

	//语言标识
	LangCode string `gorm:"column:lang_code;primaryKey;type:varchar(10);comment:语言标识" json:"langCode"`

	//菜单名称
	Title string `gorm:"column:title;type:varchar(20);comment:菜单名称" json:"title"`
}

// TableName 指定表名
func (MenuTl) TableName() string {
	return "sys_menu_tl"
}
