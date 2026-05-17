/**
 * 定时任务配置表 Model
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

type Job struct {

	//主键ID
	Id int64 `gorm:"column:id;primaryKey;autoIncrement;type:bigint;comment:主键ID" json:"id,string"`
	//任务名称
	JobName string `gorm:"column:job_name;type:varchar(64);comment:任务名称" json:"jobName"`
	//任务分组
	JobGroup string `gorm:"column:job_group;type:varchar(64);comment:任务分组" json:"jobGroup"`
	//任务唯一标识(不可重复)
	JobKey string `gorm:"column:job_key;type:varchar(64);comment:任务唯一标识(不可重复)" json:"jobKey"`
	//cron执行表达式
	CronExpr string `gorm:"column:cron_expr;type:varchar(32);comment:cron执行表达式" json:"cronExpr"`
	//执行目标：服务.方法名
	InvokeTarget string `gorm:"column:invoke_target;type:varchar(128);comment:执行目标：服务.方法名" json:"invokeTarget"`
	//状态 0正常 1暂停
	Status int16 `gorm:"column:status;type:tinyint;comment:状态 0正常 1暂停" json:"status"`
	//是否并发 0禁止 1允许
	Concurrent int16 `gorm:"column:concurrent;type:tinyint;comment:是否并发 0禁止 1允许" json:"concurrent"`
	//失败重试次数
	RetryCount int `gorm:"column:retry_count;type:int;comment:失败重试次数" json:"retryCount"`
	//执行超时时间(秒)
	Timeout int `gorm:"column:timeout;type:int;comment:执行超时时间(秒)" json:"timeout"`
	//备注说明
	Remark string `gorm:"column:remark;type:varchar(255);comment:备注说明" json:"remark"`
	base.DefaultModel
}

// TableName 指定表名
func (Job) TableName() string {
	return "sys_job"
}
