/**
 * 用户信息表 Model
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

type User struct {

	//用户编号
	Id int64 `gorm:"column:id;primaryKey;autoIncrement;type:bigint;comment:用户编号" json:"id,string"`
	//用户名
	Username string `gorm:"column:username;type:varchar(32);comment:用户名" json:"username"`
	//呢称
	Nickname string `gorm:"column:nickname;type:varchar(50);comment:呢称" json:"nickname"`
	//用户类型
	Type int16 `gorm:"column:type;type:tinyint;comment:用户类型" json:"type"`
	//邮箱
	Email string `gorm:"column:email;type:varchar(64);comment:邮箱" json:"email"`
	//手机号
	Mobile string `gorm:"column:mobile;type:varchar(11);comment:手机号" json:"mobile"`
	//性别（0未知 1男 2女）
	Sex int16 `gorm:"column:sex;type:tinyint;comment:性别（0未知 1男 2女）" json:"sex"`
	//头像
	Avatar string `gorm:"column:avatar;type:varchar(255);comment:头像" json:"avatar"`
	//电子签名
	Autograph string `gorm:"column:autograph;type:varchar(255);comment:电子签名" json:"autograph"`
	//密码
	Password string `gorm:"column:password;type:varchar(128);comment:密码" json:"-"`
	//状态（0正常 1停用）
	Status int16 `gorm:"column:status;type:tinyint;comment:状态（0正常 1停用）" json:"status"`
	base.DefaultModel
}

// TableName 指定表名
func (User) TableName() string {
	return "sys_user"
}
