/**
 * 操作日志记录 Model
 *
 * @author
 * @date 2026-04-14
 * @website https://www.aevons.com
 * @email 1073602@qq.com
 * @copyright 2025-2026 Aevons
 */
package model

import "time"

type LoginLog struct {

	//编号
	Id int64 `gorm:"column:id;primaryKey;autoIncrement;type:bigint;comment:编号" json:"id,string"`

	//登录用户名
	Username string `gorm:"column:username;type:varchar(50);comment:登录用户名" json:"username"`

	//客户端ID
	ClientId string `gorm:"column:client_id;type:varchar(20);comment:客户端ID" json:"clientId"`

	//授权类型
	GrantType string `gorm:"column:grant_type;type:varchar(20);comment:授权类型" json:"grantType"`

	//系统
	Os string `gorm:"column:os;type:varchar(50);comment:系统" json:"os"`

	//浏览器
	Browser string `gorm:"column:browser;type:varchar(50);comment:浏览器" json:"browser"`

	//登录IP
	Ip string `gorm:"column:ip;type:varchar(128);comment:登录IP" json:"ip"`

	//登录地点
	Location string `gorm:"column:location;type:varchar(255);comment:登录地点" json:"location"`

	//状态（0失败 1成功）
	Status int16 `gorm:"column:status;type:tinyint(1);comment:状态（0失败 1成功）" json:"status"`

	//模块标题
	Msg string `gorm:"column:msg;type:varchar(50);comment:模块标题" json:"msg"`

	//登录时间
	LoginAt time.Time `gorm:"column:login_at;type:datetime;comment:登录时间" json:"loginAt"`
}

// TableName 指定表名
func (LoginLog) TableName() string {
	return "sys_login_log"
}
