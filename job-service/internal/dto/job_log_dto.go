/**
 * 定时任务执行日志表 DTO
 *
 * @author
 * @date 2026-04-19
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
type JobLogQuery struct {

	// 任务ID
	JobId *int64 `form:"jobId" gorm:"column:job_id" q:"eq"`
	// 任务名称（冗余）
	JobName *string `form:"jobName" gorm:"column:job_name" q:"like"`
	// 任务分组（冗余）
	JobGroup *string `form:"jobGroup" gorm:"column:job_group" q:"eq"`
	// 执行状态 0成功 1失败 2进行中
	Status *int16 `form:"status" gorm:"column:status" q:"eq"`
	// 触发类型：自动/手动
	TriggerType *string `form:"triggerType" gorm:"column:trigger_type" q:"eq"`
	// 日志生成时间
	CreatedAt *[]time.Time `form:"createdAt" gorm:"column:created_at" q:"between"`
	base.BaseQuery
}

type CreateJobLogDTO struct {

	//任务ID
	JobId int64 `json:"jobId,string" binding:"required"`
	//任务名称（冗余）
	JobName string `json:"jobName" binding:"max=64"`
	//任务分组（冗余）
	JobGroup string `json:"jobGroup" binding:"max=64"`
	//执行状态 0成功 1失败 2进行中
	Status int16 `json:"status"`
	//执行日志/异常信息
	Message string `json:"message"`
	//执行耗时(毫秒)
	Duration int `json:"duration"`
	//触发类型：自动/手动
	TriggerType string `json:"triggerType" binding:"max=16"`
	//开始时间
	StartTime time.Time `json:"startTime"`
	//结束时间
	EndTime time.Time `json:"endTime"`
	//日志生成时间
	CreatedAt time.Time `json:"createdAt"`
}

type UpdateJobLogDTO struct {
}

type JobLogDTO struct {
	//日志ID
	Id int64 `excel:"column:日志ID;index:0;dict:" json:"id,string" binding:"required"`
	//任务ID
	JobId int64 `excel:"column:任务ID;index:1;dict:" json:"jobId,string" binding:"required"`
	//任务名称（冗余）
	JobName string `excel:"column:任务名称（冗余）;index:2;dict:" json:"jobName" binding:"max=64"`
	//任务分组（冗余）
	JobGroup string `excel:"column:任务分组（冗余）;index:3;dict:" json:"jobGroup" binding:"max=64"`
	//执行状态 0成功 1失败 2进行中
	Status int16 `excel:"column:执行状态 0成功 1失败 2进行中;index:4;dict:" json:"status"`
	//执行日志/异常信息
	Message string `excel:"column:执行日志/异常信息;index:5;dict:" json:"message"`
	//执行耗时(毫秒)
	Duration int `excel:"column:执行耗时(毫秒);index:6;dict:" json:"duration"`
	//触发类型：自动/手动
	TriggerType string `excel:"column:触发类型：自动/手动;index:7;dict:" json:"triggerType" binding:"max=16"`
	//开始时间
	StartTime time.Time `excel:"column:开始时间;index:8;dict:" json:"startTime"`
	//结束时间
	EndTime time.Time `excel:"column:结束时间;index:9;dict:" json:"endTime"`
	//日志生成时间
	CreatedAt time.Time `excel:"column:日志生成时间;index:10;dict:" json:"createdAt"`
}
