/**
 * 语言资源 Model
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

type LangResource struct {

	//主键ID
	Id int64 `gorm:"column:id;primaryKey;autoIncrement;type:bigint;comment:主键ID" json:"id,string"`
	//资源标识
	ResourceKey string `gorm:"column:resource_key;type:varchar(128);comment:资源标识" json:"resourceKey"`
	//命名空间（default）
	Namespace string `gorm:"column:namespace;type:varchar(32);comment:命名空间（default）" json:"namespace"`
	//语言编码（如zh）
	LangCode string `gorm:"column:lang_code;type:varchar(10);comment:语言编码（如zh）" json:"langCode"`
	//内容
	Content string `gorm:"column:content;type:varchar(500);comment:内容" json:"content"`
	//状态（0正常,1停用）
	Status int16 `gorm:"column:status;type:tinyint(1);comment:状态（0正常,1停用）" json:"status"`
	base.BaseModel
}

// TableName 指定表名
func (LangResource) TableName() string {
	return "sys_lang_resource"
}
