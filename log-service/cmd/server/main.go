// @title           Aevons Log Service API
// @version         1.0.0
// @description     Aevons 日志服务接口文档，提供登录日志、操作日志的查询、清理以及跨服务日志写入能力
// @host            localhost:10703
// @BasePath        /api/log/v1
// @schemes         http https
// @securityDefinitions.apikey  BearerAuth
// @in                          header
// @name                        Authorization
// @description                 Bearer 访问令牌，格式：Bearer {token}

package main

import (
	"log-service/internal/router"

	"internal-grpc/log_grpc"
	"log-service/internal/grpcs"

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

	// gRPC 端口属于服务级配置，框架启动后再从 App 中读取。
	cfg, err := app.RawConfig()
	if err != nil {
		xlog.Fatal("failed to read config from app: %v", err)
	}

	// log-service 只保留路由装配和 gRPC 服务注册这类业务特有逻辑。
	engine, err := router.Setup(app)
	if err != nil {
		xlog.Fatal("failed to setup router: %v", err)
	}

	var grpcSrv *grpc.Server

	if cfg.Server.GRPCPort > 0 {
		grpcSrv = grpcx.NewServer()
		dbConn, dbErr := app.RawDatabase()
		if dbErr != nil {
			xlog.Fatal("failed to read database from app: %v", dbErr)
		}
		log_grpc.RegisterService(grpcSrv, grpcs.NewOperLogServiceServer(dbConn))
		log_grpc.RegisterLoginService(grpcSrv, grpcs.NewLoginLogServiceServer(dbConn))
	}

	// HTTP、gRPC、Consul 注册和优雅关闭统一由框架接管。
	if err := core.RunGinAndGRPC(app, engine, grpcSrv, cfg.Server.GRPCPort); err != nil {
		xlog.Fatal("failed to run http/grpc server: %v", err)
	}
}
