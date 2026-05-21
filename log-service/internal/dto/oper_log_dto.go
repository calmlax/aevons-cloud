/**
 * 操作日志记录 DTO
 *
 * @author
 * @date 2026-04-14
 * @website https://www.aevons.com
 * @email 1073602@qq.com
 * @copyright 2025-2026 Aevons
 */
package dto

import (
	"time"

	"github.com/calmlax/aevons-framework/core/base"
)

// q 支持查询操作: eq, ne, gt, gte, lt, lte, like, like_l, like_r, in, not_in, between, is_null, not_null
// gorm.column 必须对应真实数据库列名
type OperLogQuery struct {

	// 模块
	Module *string `form:"module" gorm:"column:module" q:"eq"`
	// 类型
	Type *string `form:"type" gorm:"column:type" q:"eq"`
	// 描述
	Description *string `form:"description" gorm:"column:description" q:"like"`
	// 请求方法
	Method *string `form:"method" gorm:"column:method" q:"eq"`
	// 设备
	Device *string `form:"device" gorm:"column:device" q:"like"`
	// 系统
	Os *string `form:"os" gorm:"column:os" q:"like"`
	// 状态（0失败 1成功）
	Status *int16 `form:"status" gorm:"column:status" q:"eq"`
	// 用户编号
	UserId *int64 `form:"userId" gorm:"column:user_id" q:"eq"`
	// 用户名
	Username *string `form:"username" gorm:"column:username" q:"like"`
	// 操作时间
	OperAt *time.Time `form:"operAt" gorm:"column:oper_at" q:"eq"`
	base.BaseQuery
}

type CreateOperLogDTO struct {

	//模块
	Module string `json:"module" binding:"max=50"`
	//类型
	Type string `json:"type" binding:"max=10"`
	//描述
	Description string `json:"description" binding:"max=128"`
	//请求方法
	Method string `json:"method" binding:"max=10"`
	//请求URL
	Url string `json:"url" binding:"max=255"`
	//请求IP
	Ip string `json:"ip" binding:"max=128"`
	//请求地点
	Location string `json:"location" binding:"max=255"`
	//请求参数
	Payload string `json:"payload" binding:"max=2000"`
	//响应数据
	Result string `json:"result" binding:"max=2000"`
	//设备
	Device string `json:"device" binding:"max=50"`
	//系统
	Os string `json:"os" binding:"max=50"`
	//浏览器
	Browser string `json:"browser" binding:"max=50"`
	//状态（0失败 1成功）
	Status int16 `json:"status"`
	//错误消息
	Error string `json:"error" binding:"max=2000"`
	//耗时（毫秒）
	Time int64 `json:"time,string"`
	//用户编号
	UserId int64 `json:"userId,string"`
	//用户名
	Username string `json:"username" binding:"max=50"`
	//操作时间
	OperAt time.Time `json:"operAt"`
}

type UpdateOperLogDTO struct {
	//编号
	Id *int64 `json:"id,string" binding:"required"`
	//模块
	Module *string `json:"module" binding:"max=50"`
	//类型
	Type *string `json:"type" binding:"max=10"`
	//描述
	Description *string `json:"description" binding:"max=128"`
	//请求方法
	Method *string `json:"method" binding:"max=10"`
	//请求URL
	Url *string `json:"url" binding:"max=255"`
	//请求IP
	Ip *string `json:"ip" binding:"max=128"`
	//请求地点
	Location *string `json:"location" binding:"max=255"`
	//请求参数
	Payload *string `json:"payload" binding:"max=2000"`
	//响应数据
	Result *string `json:"result" binding:"max=2000"`
	//设备
	Device *string `json:"device" binding:"max=50"`
	//系统
	Os *string `json:"os" binding:"max=50"`
	//浏览器
	Browser *string `json:"browser" binding:"max=50"`
	//状态（0失败 1成功）
	Status *int16 `json:"status"`
	//错误消息
	Error *string `json:"error" binding:"max=2000"`
	//耗时（毫秒）
	Time *int64 `json:"time,string"`
	//用户编号
	UserId *int64 `json:"userId,string"`
	//用户名
	Username *string `json:"username" binding:"max=50"`
	//操作时间
	OperAt *time.Time `json:"operAt"`
}

type OperLogDTO struct {
	//编号
	Id int64 `json:"id,string" binding:"required"`
	//模块
	Module string `json:"module" binding:"max=50"`
	//类型
	Type string `json:"type" binding:"max=10"`
	//描述
	Description string `json:"description" binding:"max=128"`
	//请求方法
	Method string `json:"method" binding:"max=10"`
	//请求URL
	Url string `json:"url" binding:"max=255"`
	//请求IP
	Ip string `json:"ip" binding:"max=128"`
	//请求地点
	Location string `json:"location" binding:"max=255"`
	//请求参数
	Payload string `json:"payload" binding:"max=2000"`
	//响应数据
	Result string `json:"result" binding:"max=2000"`
	//设备
	Device string `json:"device" binding:"max=50"`
	//系统
	Os string `json:"os" binding:"max=50"`
	//浏览器
	Browser string `json:"browser" binding:"max=50"`
	//状态（0失败 1成功）
	Status int16 `json:"status"`
	//错误消息
	Error string `json:"error" binding:"max=2000"`
	//耗时（毫秒）
	Time int64 `json:"time,string"`
	//用户编号
	UserId int64 `json:"userId,string"`
	//用户名
	Username string `json:"username" binding:"max=50"`
	//操作时间
	OperAt time.Time `json:"operAt"`
}
