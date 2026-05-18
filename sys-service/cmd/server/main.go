// @title           Aevons System Service API
// @version         1.0.0
// @description     Aevons 系统服务接口文档，提供用户、角色、部门、岗位、菜单、字典、语言、配置、通知公告和终端应用等系统管理能力
// @host            localhost:10702
// @BasePath        /api/v1/sys
// @schemes         http https
// @securityDefinitions.apikey  BearerAuth
// @in                          header
// @name                        Authorization
// @description                 Bearer 访问令牌，格式：Bearer {token}

package main

import (
	"sys-service/internal/router"

	"github.com/calmlax/aevons-framework/core"
	"github.com/calmlax/aevons-framework/xlog"
)

func main() {
	// 统一完成配置、日志、JSON、数据库和 Redis 的基础装配。
	app, err := core.Bootstrap()
	if err != nil {
		xlog.Fatal("failed to bootstrap app: %v", err)
	}

	// sys-service 只保留服务特有的路由与依赖装配。
	engine, err := router.Setup(app)
	if err != nil {
		xlog.Fatal("failed to setup router: %v", err)
	}

	// HTTP 启动、Consul 注册和优雅关闭统一由框架接管。
	if err := core.RunGin(app, engine); err != nil {
		xlog.Fatal("failed to run http server: %v", err)
	}
}
