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
type LoginLogQuery struct {

	// 登录用户名
	Username *string `form:"username" gorm:"column:username" q:"like"`
	// 客户端ID
	ClientId *string `form:"clientId" gorm:"column:client_id" q:"eq"`
	// 授权类型
	GrantType *string `form:"grantType" gorm:"column:grant_type" q:"eq"`
	// 系统
	Os *string `form:"os" gorm:"column:os" q:"eq"`
	// 状态（0失败 1成功）
	Status *int16 `form:"status" gorm:"column:status" q:"eq"`
	base.BaseQuery
}

type CreateLoginLogDTO struct {
}

type UpdateLoginLogDTO struct {
}

type LoginLogDTO struct {
	//编号
	Id int64 `json:"id,string" binding:"required"`
	//登录用户名
	Username string `json:"username" binding:"required,max=50"`
	//客户端ID
	ClientId string `json:"clientId" binding:"required,max=32"`
	//授权类型
	GrantType string `json:"grantType" binding:"max=20"`
	//系统
	Os string `json:"os" binding:"max=50"`
	//浏览器
	Browser string `json:"browser" binding:"max=50"`
	//登录IP
	Ip string `json:"ip" binding:"max=128"`
	//登录地点
	Location string `json:"location" binding:"max=255"`
	//状态（0失败 1成功）
	Status int16 `json:"status"`
	//模块标题
	Msg string `json:"msg" binding:"max=50"`
	//登录时间
	LoginAt time.Time `json:"loginAt"`
}
