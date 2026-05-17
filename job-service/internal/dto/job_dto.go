/**
 * 定时任务配置表 DTO
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
type JobQuery struct {

	// 任务名称
	JobName *string `form:"jobName" gorm:"column:job_name" q:"like"`
	// 任务唯一标识(不可重复)
	JobKey *string `form:"jobKey" gorm:"column:job_key" q:"eq"`
	// 状态 0正常 1暂停
	Status *int16 `form:"status" gorm:"column:status" q:"eq"`
	base.BaseQuery
}

type CreateJobDTO struct {

	//任务名称
	JobName string `json:"jobName" binding:"required,max=64"`
	//任务分组
	JobGroup string `json:"jobGroup" binding:"max=64"`
	//任务唯一标识(不可重复)
	JobKey string `json:"jobKey" binding:"required,max=64"`
	//cron执行表达式
	CronExpr string `json:"cronExpr" binding:"required,max=32"`
	//执行目标：服务.方法名
	InvokeTarget string `json:"invokeTarget" binding:"max=128"`
	//状态 0正常 1暂停
	Status int16 `json:"status" binding:"required"`
	//是否并发 0禁止 1允许
	Concurrent int16 `json:"concurrent" binding:"required"`
	//失败重试次数
	RetryCount int `json:"retryCount" binding:"required"`
	//执行超时时间(秒)
	Timeout int `json:"timeout" binding:"required"`
	//备注说明
	Remark string `json:"remark" binding:"max=255"`
}

type UpdateJobDTO struct {
	//主键ID
	Id *int64 `json:"id,string" binding:"required"`
	//任务名称
	JobName *string `json:"jobName" binding:"required,max=64"`
	//任务分组
	JobGroup *string `json:"jobGroup" binding:"max=64"`
	//任务唯一标识(不可重复)
	JobKey *string `json:"jobKey" binding:"required,max=64"`
	//cron执行表达式
	CronExpr *string `json:"cronExpr" binding:"required,max=32"`
	//执行目标：服务.方法名
	InvokeTarget *string `json:"invokeTarget" binding:"max=128"`
	//状态 0正常 1暂停
	Status *int16 `json:"status" binding:"required"`
	//是否并发 0禁止 1允许
	Concurrent *int16 `json:"concurrent" binding:"required"`
	//失败重试次数
	RetryCount *int `json:"retryCount" binding:"required"`
	//执行超时时间(秒)
	Timeout *int `json:"timeout" binding:"required"`
	//备注说明
	Remark *string `json:"remark" binding:"max=255"`
}

type JobDTO struct {
	//主键ID
	Id int64 `excel:"column:主键ID;index:0;dict:" json:"id,string" binding:"required"`
	//任务名称
	JobName string `excel:"column:任务名称;index:1;dict:" json:"jobName" binding:"required,max=64"`
	//任务分组
	JobGroup string `excel:"column:任务分组;index:2;dict:" json:"jobGroup" binding:"max=64"`
	//任务唯一标识(不可重复)
	JobKey string `excel:"column:任务唯一标识(不可重复);index:3;dict:" json:"jobKey" binding:"required,max=64"`
	//cron执行表达式
	CronExpr string `excel:"column:cron执行表达式;index:4;dict:" json:"cronExpr" binding:"required,max=32"`
	//执行目标：服务.方法名
	InvokeTarget string `excel:"column:执行目标：服务.方法名;index:5;dict:" json:"invokeTarget" binding:"max=128"`
	//状态 0正常 1暂停
	Status int16 `excel:"column:状态 0正常 1暂停;index:6;dict:sys_active" json:"status" binding:"required"`
	//是否并发 0禁止 1允许
	Concurrent int16 `excel:"column:是否并发 0禁止 1允许;index:7;dict:sys_is" json:"concurrent" binding:"required"`
	//失败重试次数
	RetryCount int `excel:"column:失败重试次数;index:8;dict:" json:"retryCount" binding:"required"`
	//执行超时时间(秒)
	Timeout int `excel:"column:执行超时时间(秒);index:9;dict:" json:"timeout" binding:"required"`
	//备注说明
	Remark string `excel:"column:备注说明;index:10;dict:" json:"remark" binding:"max=255"`
	//创建人ID
	CreatedBy int64 `excel:"column:创建人ID;index:11;dict:" json:"createdBy,string"`
	//创建时间
	CreatedAt time.Time `excel:"column:创建时间;index:12;dict:" json:"createdAt"`
	//更新时间
	UpdatedAt time.Time `excel:"column:更新时间;index:13;dict:" json:"updatedAt"`
	//更新人ID
	UpdatedBy int64 `excel:"column:更新人ID;index:14;dict:" json:"updatedBy,string"`
}
