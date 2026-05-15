/**
 * 参数配置表 Model
 *
 * @author
 * @date 2026-04-09 00:38:25.504785055 +0000 UTC
 * @website https://www.aevons.com
 * @email 1073602@qq.com
 * @copyright 2025-2026 Aevons
 */
package model

import (
	"github.com/calmlax/aevons-framework/core/base"
)

type Conf struct {

	//编号
	Id int64 `gorm:"column:id;primaryKey;autoIncrement;type:bigint;comment:编号" json:"id,string"`

	//名称
	Name string `gorm:"column:name;type:varchar(255);comment:名称" json:"name"`

	//配置KEY
	ConfKey string `gorm:"column:conf_key;type:varchar(64);comment:配置KEY" json:"confKey"`

	//配置值
	ConfValue string `gorm:"column:conf_value;type:text;comment:配置值" json:"confValue"`

	//是否系统内置（0否 1是）
	IsSys int16 `gorm:"column:is_sys;type:tinyint(1);comment:是否系统内置（0否 1是）" json:"isSys"`

	//范围（0公开配置，1登录配置，2后台服务配置）
	Scope int16 `gorm:"column:scope;type:tinyint(1);comment:范围（0公开配置，1登录配置，2后台服务配置）" json:"scope"`

	//是否加密/脱敏（0否 1是）
	IsSecret int16 `gorm:"column:is_secret;type:tinyint(1);comment:是否加密/脱敏（0否 1是）" json:"isSecret"`

	//备注
	Remark string `gorm:"column:remark;type:varchar(500);comment:备注" json:"remark"`

	base.DefaultModel
}

// TableName 指定表名
func (Conf) TableName() string {
	return "sys_conf"
}
