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

type OperLog struct {

	//编号
	Id int64 `gorm:"column:id;primaryKey;autoIncrement;type:bigint;comment:编号" json:"id,string"`

	//模块
	Module string `gorm:"column:module;type:varchar(50);comment:模块" json:"module"`

	//类型
	Type string `gorm:"column:type;type:char(10);comment:类型" json:"type"`

	//描述
	Description string `gorm:"column:description;type:varchar(128);comment:描述" json:"description"`

	//请求方法
	Method string `gorm:"column:method;type:varchar(10);comment:请求方法" json:"method"`

	//请求URL
	Url string `gorm:"column:url;type:varchar(255);comment:请求URL" json:"url"`

	//请求IP
	Ip string `gorm:"column:ip;type:varchar(128);comment:请求IP" json:"ip"`

	//请求地点
	Location string `gorm:"column:location;type:varchar(255);comment:请求地点" json:"location"`

	//请求参数
	Payload string `gorm:"column:payload;type:varchar(2000);comment:请求参数" json:"payload"`

	//响应数据
	Result string `gorm:"column:result;type:varchar(2000);comment:响应数据" json:"result"`

	//设备
	Device string `gorm:"column:device;type:varchar(50);comment:设备" json:"device"`

	//系统
	Os string `gorm:"column:os;type:varchar(50);comment:系统" json:"os"`

	//浏览器
	Browser string `gorm:"column:browser;type:varchar(50);comment:浏览器" json:"browser"`

	//状态（0失败 1成功）
	Status int16 `gorm:"column:status;type:tinyint(1);comment:状态（0失败 1成功）" json:"status"`

	//错误消息
	Error string `gorm:"column:error;type:varchar(2000);comment:错误消息" json:"error"`

	//耗时（毫秒）
	Time int64 `gorm:"column:time;type:bigint;comment:耗时（毫秒）" json:"time,string"`

	//用户编号
	UserId int64 `gorm:"column:user_id;type:bigint;comment:用户编号" json:"userId,string"`

	//用户名
	Username string `gorm:"column:username;type:varchar(50);comment:用户名" json:"username"`

	//操作时间
	OperAt time.Time `gorm:"column:oper_at;type:datetime;comment:操作时间" json:"operAt"`
}

// TableName 指定表名
func (OperLog) TableName() string {
	return "sys_oper_log"
}
