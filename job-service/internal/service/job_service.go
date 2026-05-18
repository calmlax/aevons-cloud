package service

import (
	"job-service/internal/dto"
	"job-service/internal/model"
	"job-service/internal/repository"
	"job-service/internal/scheduler"
	"time"

	"github.com/calmlax/aevons-framework/core/xlog"

	"github.com/calmlax/aevons-framework/core/base"
)

// JobService 继承 BaseService 并扩展调度管理能力
type JobService interface {
	base.BaseService[model.Job, *dto.JobQuery]
	// 启动任务（注册到 cron 引擎）
	StartJob(job *model.Job) error
	// 停止任务（从 cron 引擎移除）
	StopJob(jobId int64) error
	// 手动触发任务
	TriggerJob(jobId int64) error
	// 更新任务状态（0正常/1暂停）并同步调度引擎
	ChangeStatus(jobId int64, status int16) error
	// 初始化：加载所有启用任务到调度引擎
	InitScheduler() error
}

type jobService struct {
	base.BaseService[model.Job, *dto.JobQuery]
	repo    repository.JobRepository
	logRepo repository.JobLogRepository
}

func NewJobService(repo repository.JobRepository, logRepo repository.JobLogRepository) JobService {
	baseSrv := base.NewBaseService[model.Job, *dto.JobQuery](repo)
	return &jobService{
		BaseService: baseSrv,
		repo:        repo,
		logRepo:     logRepo,
	}
}

// InitScheduler 服务启动时加载所有状态=0（正常）的任务
func (s *jobService) InitScheduler() error {
	jobs, err := s.repo.ListByField("status", int16(0))
	if err != nil {
		return err
	}
	for _, job := range jobs {
		j := job
		if err := s.StartJob(&j); err != nil {
			xlog.Error("init scheduler: failed to start job %d (%s): %v", j.Id, j.JobKey, err)
		}
	}
	xlog.Info("scheduler initialized, loaded %d jobs", len(jobs))
	return nil
}

// StartJob 将任务注册到 cron 引擎
func (s *jobService) StartJob(job *model.Job) error {
	sch := scheduler.Instance()
	return sch.AddJob(
		job.Id,
		job.CronExpr,
		job.InvokeTarget,
		"",
		job.Concurrent == 1,
		job.Timeout,
		job.RetryCount,
		s.buildLogCallback(job),
	)
}

// StopJob 从 cron 引擎移除任务
func (s *jobService) StopJob(jobId int64) error {
	scheduler.Instance().RemoveJob(jobId)
	return nil
}

// TriggerJob 手动触发任务
func (s *jobService) TriggerJob(jobId int64) error {
	job, err := s.repo.GetById(jobId)
	if err != nil {
		return err
	}
	return scheduler.Instance().TriggerJob(
		job.Id,
		job.InvokeTarget,
		"",
		job.Timeout,
		job.RetryCount,
		s.buildLogCallback(job),
	)
}

// ChangeStatus 更新任务状态并同步调度引擎
func (s *jobService) ChangeStatus(jobId int64, status int16) error {
	job, err := s.repo.GetById(jobId)
	if err != nil {
		return err
	}
	if _, err := s.repo.Update(jobId, map[string]any{"status": status}); err != nil {
		return err
	}
	if status == 0 {
		job.Status = 0
		return s.StartJob(job)
	}
	return s.StopJob(jobId)
}

// buildLogCallback 构建日志写入回调
func (s *jobService) buildLogCallback(job *model.Job) scheduler.LogCallback {
	return func(jobId int64, triggerType string, status int16, msg string, duration int, start, end time.Time) {
		entry := &model.JobLog{
			JobId:       jobId,
			JobName:     job.JobName,
			JobGroup:    job.JobGroup,
			Status:      status,
			Message:     msg,
			Duration:    duration,
			TriggerType: triggerType,
			StartTime:   start,
			EndTime:     end,
			CreatedAt:   time.Now(),
		}
		if err := s.logRepo.Create(entry); err != nil {
			xlog.Error("job log write failed: jobId=%d status=%d err=%v", jobId, status, err)
		}
	}
}
