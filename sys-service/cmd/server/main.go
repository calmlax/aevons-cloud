// @title           Aevons System Service API
// @version         1.0.0
// @description     Aevons 系统服务接口文档，提供用户、角色、部门、岗位、菜单、字典、语言、配置、通知公告和终端应用等系统管理能力
// @host            localhost:10702
// @BasePath        /api/sys/v1
// @schemes         http https
// @securityDefinitions.apikey  BearerAuth
// @in                          header
// @name                        Authorization
// @description                 Bearer 访问令牌，格式：Bearer {token}

package main

import (
	"internal-grpc/sys_grpc"
	"sys-service/internal/grpcs"
	"sys-service/internal/router"

	"github.com/calmlax/aevons-framework/core"
	"github.com/calmlax/aevons-framework/grpcx"
	"github.com/calmlax/aevons-framework/xlog"
	"google.golang.org/grpc"
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

	cfg, err := app.RawConfig()
	if err != nil {
		xlog.Fatal("failed to read config from app: %v", err)
	}

	var grpcSrv *grpc.Server

	if cfg.Server.GRPCPort > 0 {
		grpcSrv = grpcx.NewServer()
		dbConn, dbErr := app.RawDatabase()
		if dbErr != nil {
			xlog.Fatal("failed to read database from app: %v", dbErr)
		}
		sys_grpc.RegisterService(grpcSrv, grpcs.NewSysServiceServer(dbConn))
	}

	// HTTP、gRPC、Consul 注册和优雅关闭统一由框架接管。
	if err := core.RunGinAndGRPC(app, engine, grpcSrv, cfg.Server.GRPCPort); err != nil {
		xlog.Fatal("failed to run http/grpc server: %v", err)
	}
}
