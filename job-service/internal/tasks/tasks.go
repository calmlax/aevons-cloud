// Package tasks 注册所有业务定时任务处理函数。
// 新增任务：在此文件中调用 scheduler.Instance().Register("任务标识", 处理函数) 即可。
package tasks

import (
	"context"
	"fmt"
	"time"

	"job-service/internal/scheduler"

	"github.com/calmlax/aevons-framework/xlog"
)

// Init 注册所有任务处理函数，在应用启动时调用一次
func Init() {
	s := scheduler.Instance()

	// 示例任务：系统心跳检测
	s.Register("sys.heartbeat", func(ctx context.Context, param string) (string, error) {
		xlog.Info("[job] sys.heartbeat executed at %s", time.Now().Format(time.DateTime))
		return "ok", nil
	})

	// 示例任务：数据清理（可传 param 参数）
	s.Register("sys.cleanLog", func(ctx context.Context, param string) (string, error) {
		xlog.Info("[job] sys.cleanLog executed, param=%s", param)
		// TODO: 实现日志清理逻辑
		return fmt.Sprintf("cleaned with param=%s", param), nil
	})

	// 测试任务：打印当前时间 + param，用于验证调度器是否正常工作
	s.Register("sys.test", func(ctx context.Context, param string) (string, error) {
		msg := fmt.Sprintf("[job] sys.test fired at %s", time.Now().Format(time.DateTime))
		if param != "" {
			msg += fmt.Sprintf(", param=%s", param)
		}
		xlog.Info(msg)
		return msg, nil
	})

	xlog.Info("job tasks registered: %v", s.ListTasks())
}
