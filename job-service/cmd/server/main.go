// @title           Aevons Job Service API
// @version         1.0.0
// @description     Aevons 定时任务服务接口文档，提供任务配置、状态切换、手动触发以及任务执行日志查询等能力
// @host            localhost:10705
// @BasePath        /api/job/v1
// @schemes         http https
// @securityDefinitions.apikey  BearerAuth
// @in                          header
// @name                        Authorization
// @description                 Bearer 访问令牌，格式：Bearer {token}

package main

import (
	"job-service/internal/router"

	"github.com/calmlax/aevons-framework/core"
	"github.com/calmlax/aevons-framework/xlog"
)

func main() {
	app, err := core.Bootstrap()
	if err != nil {
		xlog.Fatal("failed to bootstrap app: %v", err)
	}

	engine, err := router.Setup(app)
	if err != nil {
		xlog.Fatal("failed to setup router: %v", err)
	}

	if err := core.RunGin(app, engine); err != nil {
		xlog.Fatal("failed to run http server: %v", err)
	}
}
