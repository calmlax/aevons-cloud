// @title           Aevons Log Service API
// @version         1.0.0
// @description     Aevons 日志服务接口文档，提供登录日志、操作日志的查询、清理以及跨服务日志写入能力
// @host            localhost:10703
// @BasePath        /api/v1/log
// @schemes         http https
// @securityDefinitions.apikey  BearerAuth
// @in                          header
// @name                        Authorization
// @description                 Bearer 访问令牌，格式：Bearer {token}

package main

import (
	"context"
	"fmt"
	"log-service/router"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"aevons-grpc/log_grpc"
	"log-service/grpcs"

	"github.com/calmlax/aevons-framework/config"
	"github.com/calmlax/aevons-framework/core"
	"github.com/calmlax/aevons-framework/core/consul"
	"github.com/calmlax/aevons-framework/db"
	"github.com/calmlax/aevons-framework/grpcx"
	"github.com/calmlax/aevons-framework/redis"
	"github.com/calmlax/aevons-framework/xjson"
	"github.com/calmlax/aevons-framework/xlog"

	"github.com/gin-gonic/gin/binding"
	"google.golang.org/grpc"
)

func init() {
	// 启用 UseNumber，避免 JSON 数字在解码时默认丢失整数语义。
	binding.EnableDecoderUseNumber = true
	// 初始化自定义 JSON 绑定与编码行为。
	xjson.InitGin()
}

func main() {
	// 允许通过 --config 指定配置目录，默认使用当前项目下的 configs 目录。
	// APP_ENV 用于选择环境覆盖配置，例如 development、staging、production。
	configPath := "configs"
	for i, arg := range os.Args[1:] {
		if arg == "--config" && i+1 < len(os.Args[1:]) {
			configPath = os.Args[i+2]
			break
		}
	}
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}

	// 加载基础配置与环境覆盖配置。
	cfg, err := config.Load(configPath, env)
	if err != nil {
		xlog.Fatal("failed to load config: %v", err)
	}

	// 初始化数据库连接与连接池。
	gdb, err := db.Init(&cfg)
	if err != nil {
		xlog.Fatal("failed to init db: %v", err)
	}

	// 初始化 Redis 客户端。
	redisClient, err := redis.Init(&cfg)
	if err != nil {
		xlog.Fatal("failed to init redis: %v", err)
	}

	// App 仅作为装配层上下文使用，后续将显式拆出路由层需要的依赖。
	app := core.NewApp(&cfg, redisClient, gdb)

	// 构建 HTTP 路由与中间件链。
	engine, err := router.Setup(app)
	if err != nil {
		xlog.Fatal("failed to setup router: %v", err)
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: engine,
	}

	var grpcSrv *grpc.Server

	// 先显式监听端口，确保端口占用问题能在启动早期暴露出来。
	listener, err := net.Listen("tcp", srv.Addr)
	if err != nil {
		xlog.Fatal("failed to listen on %s: %v", srv.Addr, err)
	}

	var consulManager *consul.Managed

	// 启动 HTTP 服务主循环。
	go func() {
		xlog.Info("server starting on port %d", cfg.Server.Port)
		if err := srv.Serve(listener); err != nil && err != http.ErrServerClosed {
			xlog.Error("server failed to start: %v", err)
			os.Exit(1)
		}
	}()

	if cfg.Server.GRPCPort > 0 {
		grpcListener, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Server.GRPCPort))
		if err != nil {
			_ = listener.Close()
			xlog.Fatal("failed to listen on grpc port %d: %v", cfg.Server.GRPCPort, err)
		}

		grpcSrv = grpcx.NewServer()
		log_grpc.RegisterService(grpcSrv, grpcs.NewOperLogServiceServer(gdb))
		log_grpc.RegisterLoginService(grpcSrv, grpcs.NewLoginLogServiceServer(gdb))

		go func() {
			xlog.Info("grpc server starting on port %d", cfg.Server.GRPCPort)
			if err := grpcSrv.Serve(grpcListener); err != nil {
				xlog.Error("grpc server failed to start: %v", err)
				os.Exit(1)
			}
		}()
	}

	// 如果启用了 Consul，则在健康检查接口就绪后再执行服务注册，
	// 避免服务刚启动就因为探测过早而被标记为不健康。
	if cfg.Consul.Enabled {
		consulManager, err = consul.NewManaged(cfg.Consul, cfg.Server, consul.DefaultHealthPath)
		if err != nil {
			_ = listener.Close()
			xlog.Fatal("failed to init consul manager: %v", err)
		}

		if err := consulManager.Register(5 * time.Second); err != nil {
			_ = listener.Close()
			xlog.Fatal("failed to register service to consul: %v", err)
		}

		instances, discoverErr := consulManager.Discover()
		if discoverErr != nil {
			xlog.Warn("consul discover %s failed: %v", cfg.Server.Name, discoverErr)
		} else {
			xlog.Info("consul discovered %d healthy instance(s) for %s", len(instances), cfg.Server.Name)
		}
	}

	// 监听退出信号，进入优雅关闭流程。
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	xlog.Info("shutting down server...")

	// 先从 Consul 注销实例，避免下线中的节点继续接收流量。
	if consulManager != nil {
		if err := consulManager.Deregister(); err != nil {
			xlog.Error("failed to deregister service from consul: %v", err)
		}
	}

	// 为 HTTP 优雅关闭设置超时时间，给正在处理的请求留出收尾窗口。
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		xlog.Error("server forced to shutdown: %v", err)
		os.Exit(1)
	}
	if grpcSrv != nil {
		grpcSrv.GracefulStop()
	}

	// 回收基础资源连接。
	if err := redis.Close(); err != nil {
		xlog.Error("failed to close redis: %v", err)
	}
	if err := db.Close(); err != nil {
		xlog.Error("failed to close db: %v", err)
	}

	// rocketmq.Close()
	xlog.Info("server exited gracefully")
}
