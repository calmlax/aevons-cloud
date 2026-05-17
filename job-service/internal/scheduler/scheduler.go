// Package scheduler 提供基于 robfig/cron 的全局任务调度器，支持 Redis 分布式锁防重复执行。
package scheduler

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/calmlax/aevons-framework/xlog"

	"github.com/redis/go-redis/v9"
	"github.com/robfig/cron/v3"
)

// TaskFunc 任务执行函数签名
type TaskFunc func(ctx context.Context, param string) (string, error)

// LogCallback 执行完成后的日志回调
type LogCallback func(jobId int64, triggerType string, status int16, msg string, duration int, start, end time.Time)

// Scheduler 全局调度器
type Scheduler struct {
	c           *cron.Cron
	mu          sync.RWMutex
	entries     map[int64]cron.EntryID // jobId -> cronEntryId
	tasks       map[string]TaskFunc    // invokeTarget -> handler
	redisClient *redis.Client
	nodeId      string // 当前节点标识
}

var (
	instance *Scheduler
	once     sync.Once
)

// Instance 获取全局调度器单例
func Instance() *Scheduler {
	once.Do(func() {
		instance = &Scheduler{
			c:       cron.New(cron.WithSeconds()),
			entries: make(map[int64]cron.EntryID),
			tasks:   make(map[string]TaskFunc),
		}
		instance.c.Start()
	})
	return instance
}

// SetRedis 注入 Redis 客户端（用于分布式锁）
func (s *Scheduler) SetRedis(client *redis.Client, nodeId string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.redisClient = client
	s.nodeId = nodeId
}

// Register 注册任务处理函数（应用启动时调用）
func (s *Scheduler) Register(invokeTarget string, fn TaskFunc) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.tasks[invokeTarget] = fn
}

// AddJob 添加/替换一个 cron 任务
func (s *Scheduler) AddJob(jobId int64, cronExpr, invokeTarget, param string,
	concurrent bool, timeout int, retryCount int, onLog LogCallback,
) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 先移除旧的
	if eid, ok := s.entries[jobId]; ok {
		s.c.Remove(eid)
		delete(s.entries, jobId)
	}

	fn, ok := s.tasks[invokeTarget]
	if !ok {
		return fmt.Errorf("task handler not found: %s", invokeTarget)
	}

	// 并发控制：禁止并发时用互斥锁
	var running sync.Mutex

	eid, err := s.c.AddFunc(cronExpr, func() {
		if !concurrent {
			if !running.TryLock() {
				xlog.Info("[scheduler] job %d (%s) skipped: previous run still in progress", jobId, invokeTarget)
				return
			}
			defer running.Unlock()
		}

		// Redis 分布式锁：多实例部署时防重复执行
		if !s.acquireLock(jobId, timeout) {
			xlog.Info("[scheduler] job %d (%s) skipped: failed to acquire distributed lock", jobId, invokeTarget)
			return
		}
		defer s.releaseLock(jobId)

		xlog.Info("[scheduler] job %d (%s) starting, triggerType=auto", jobId, invokeTarget)
		s.execute(jobId, invokeTarget, param, "auto", timeout, retryCount, fn, onLog)
	})
	if err != nil {
		return err
	}
	s.entries[jobId] = eid
	xlog.Info("[scheduler] job %d (%s) registered with cron=%q", jobId, invokeTarget, cronExpr)
	return nil
}

// RemoveJob 移除 cron 任务
func (s *Scheduler) RemoveJob(jobId int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if eid, ok := s.entries[jobId]; ok {
		s.c.Remove(eid)
		delete(s.entries, jobId)
	}
}

// TriggerJob 手动触发任务（异步）
func (s *Scheduler) TriggerJob(jobId int64, invokeTarget, param string,
	timeout, retryCount int, onLog LogCallback,
) error {
	s.mu.RLock()
	fn, ok := s.tasks[invokeTarget]
	s.mu.RUnlock()
	if !ok {
		return fmt.Errorf("task handler not found: %s", invokeTarget)
	}
	go s.execute(jobId, invokeTarget, param, "manual", timeout, retryCount, fn, onLog)
	return nil
}

// HasTask 判断 invokeTarget 是否已注册
func (s *Scheduler) HasTask(invokeTarget string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, ok := s.tasks[invokeTarget]
	return ok
}

// ListTasks 返回所有已注册的 invokeTarget 列表
func (s *Scheduler) ListTasks() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	keys := make([]string, 0, len(s.tasks))
	for k := range s.tasks {
		keys = append(keys, k)
	}
	return keys
}

// acquireLock 尝试获取 Redis 分布式锁，无 Redis 时直接返回 true
func (s *Scheduler) acquireLock(jobId int64, timeout int) bool {
	s.mu.RLock()
	client := s.redisClient
	nodeId := s.nodeId
	s.mu.RUnlock()

	if client == nil {
		return true
	}
	ttl := time.Duration(timeout+5) * time.Second
	if ttl < 15*time.Second {
		ttl = 15 * time.Second
	}
	key := fmt.Sprintf("job:lock:%d", jobId)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	ok, err := client.SetNX(ctx, key, nodeId, ttl).Result()
	if err != nil {
		xlog.Error("[scheduler] acquireLock job %d error: %v", jobId, err)
		return false
	}
	return ok
}

// releaseLock 释放 Redis 分布式锁（仅释放自己持有的）
func (s *Scheduler) releaseLock(jobId int64) {
	s.mu.RLock()
	client := s.redisClient
	nodeId := s.nodeId
	s.mu.RUnlock()

	if client == nil {
		return
	}
	key := fmt.Sprintf("job:lock:%d", jobId)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	script := `if redis.call("get",KEYS[1]) == ARGV[1] then return redis.call("del",KEYS[1]) else return 0 end`
	client.Eval(ctx, script, []string{key}, nodeId)
}

// execute 执行任务（含超时、重试、panic 捕获、日志回调）
func (s *Scheduler) execute(
	jobId int64, invokeTarget, param, triggerType string,
	timeout, retryCount int,
	fn TaskFunc,
	onLog LogCallback,
) {
	start := time.Now()
	var (
		msg    string
		status int16 = 0
		err    error
	)

	// 兜底：捕获 panic，确保日志一定被写入
	defer func() {
		if r := recover(); r != nil {
			status = 1
			msg = fmt.Sprintf("panic: %v", r)
			xlog.Error("[scheduler] job %d (%s) panic: %v", jobId, invokeTarget, r)
		}
		end := time.Now()
		duration := int(end.Sub(start).Milliseconds())
		if onLog != nil {
			onLog(jobId, triggerType, status, msg, duration, start, end)
		}
	}()

	ctx := context.Background()
	if timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
		defer cancel()
	}

	// 重试（失败后退避重试）
	maxTry := retryCount + 1
	for i := 0; i < maxTry; i++ {
		msg, err = fn(ctx, param)
		if err == nil {
			break
		}
		xlog.Warn("[scheduler] job %d (%s) attempt %d/%d failed: %v", jobId, invokeTarget, i+1, maxTry, err)
		if i < maxTry-1 {
			time.Sleep(time.Second * time.Duration(i+1))
		}
	}

	if err != nil {
		status = 1
		errStr := err.Error()
		if ctx.Err() == context.DeadlineExceeded {
			errStr = fmt.Sprintf("timeout(%ds): %s", timeout, errStr)
		}
		// 合并任务输出和错误信息到 message
		if msg == "" {
			msg = errStr
		} else {
			msg = fmt.Sprintf("%s\n[error] %s", msg, errStr)
		}
		xlog.Error("[scheduler] job %d (%s) failed: %s", jobId, invokeTarget, msg)
	} else {
		xlog.Info("[scheduler] job %d (%s) succeeded: %s", jobId, invokeTarget, msg)
	}
}

// Stop 停止调度器
func (s *Scheduler) Stop() {
	s.c.Stop()
}
