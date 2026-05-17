/**
 * 定时任务执行日志表 Model
 *
 * @author
 * @date 2026-04-19
 * @website https://www.aevons.com
 * @email 1073602@qq.com
 * @copyright 2025-2026 Aevons
 */
package model

import "time"

type JobLog struct {

	//日志ID
	Id int64 `gorm:"column:id;primaryKey;autoIncrement;type:bigint;comment:日志ID" json:"id,string"`
	//任务ID
	JobId int64 `gorm:"column:job_id;type:bigint;comment:任务ID" json:"jobId,string"`
	//任务名称（冗余）
	JobName string `gorm:"column:job_name;type:varchar(64);comment:任务名称（冗余）" json:"jobName"`
	//任务分组（冗余）
	JobGroup string `gorm:"column:job_group;type:varchar(64);comment:任务分组（冗余）" json:"jobGroup"`
	//执行状态 0成功 1失败
	Status int16 `gorm:"column:status;type:tinyint;comment:执行状态 0成功 1失败" json:"status"`
	//执行日志/异常信息
	Message string `gorm:"column:message;type:text;comment:执行日志/异常信息" json:"message"`
	//执行耗时(毫秒)
	Duration int `gorm:"column:duration;type:int;comment:执行耗时(毫秒)" json:"duration"`
	//触发类型：auto/manual
	TriggerType string `gorm:"column:trigger_type;type:varchar(16);comment:触发类型：auto/manual" json:"triggerType"`
	//开始时间
	StartTime time.Time `gorm:"column:start_time;type:datetime;comment:开始时间" json:"startTime"`
	//结束时间
	EndTime time.Time `gorm:"column:end_time;type:datetime;comment:结束时间" json:"endTime"`
	//日志生成时间
	CreatedAt time.Time `gorm:"column:created_at;type:timestamp;comment:日志生成时间" json:"createdAt"`
}

// TableName 指定表名
func (JobLog) TableName() string {
	return "sys_job_log"
}
