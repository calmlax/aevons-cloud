package router

import (
	"fmt"
	"os"

	"aevons-grpc/log_grpc"

	"job-service/internal/handler"
	"job-service/internal/repository"
	"job-service/internal/scheduler"
	"job-service/internal/service"
	"job-service/internal/tasks"

	localmiddleware "job-service/internal/middleware"

	"github.com/calmlax/aevons-framework/auth/store"
	"github.com/calmlax/aevons-framework/consts"
	"github.com/calmlax/aevons-framework/core"
	"github.com/calmlax/aevons-framework/core/server"
	"github.com/calmlax/aevons-framework/middleware"
	"github.com/calmlax/aevons-framework/redis"
	"github.com/calmlax/aevons-framework/xlog"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

// Setup 配置 Gin 引擎，注册中间件和路由。
func Setup(app *core.App) (*gin.Engine, error) {
	cfg, err := app.RawConfig()
	if err != nil {
		return nil, fmt.Errorf("router: 读取应用配置失败: %w", err)
	}

	redisClient, err := app.RawRedis()
	if err != nil {
		return nil, fmt.Errorf("router: 读取 Redis 客户端失败: %w", err)
	}

	db, err := app.RawDatabase()
	if err != nil {
		return nil, fmt.Errorf("router: 读取数据库连接失败: %w", err)
	}

	logWriter := log_grpc.OperLogWriter(log_grpc.NopOperLogWriter{})
	client, err := log_grpc.NewOperLogClient(cfg.Consul)
	if err != nil {
		xlog.Warn("init oper log grpc client from consul failed: %v", err)
	} else {
		logWriter = client
	}

	gin.SetMode(cfg.Server.Mode)

	r := gin.New()
	r.Use(middleware.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.RequestID())
	server.RegisterHealthRoute(r, cfg.Server.Name)
	server.RegisterOpenApiRoute(r, cfg)
	r.Use(middleware.AuthMiddleware(store.NewRedisTokenStore(redisClient), cfg.Auth.Excludes))

	jobRepo := repository.NewJobRepository(db)
	logRepo := repository.NewJobLogRepository(db)
	svc := service.NewJobService(jobRepo, logRepo)
	h := handler.NewJobHandler(svc)

	v1 := r.Group("/api/v1/job")
	{
		v1.GET("ping", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "pong"})
		})

		v1.GET("/page", middleware.HasPermission("job$list"), h.Page)
		v1.GET("/:id", middleware.HasPermission("job$query"), h.Get)
		v1.POST("", middleware.HasPermission("job$add"), localmiddleware.OperLog(logWriter, "Job-[定时任务]", consts.INSERT), h.Create)
		v1.PUT("/:id", middleware.HasPermission("job$edit"), localmiddleware.OperLog(logWriter, "Job-[定时任务]", consts.UPDATE), h.Update)
		v1.DELETE("/:ids", middleware.HasPermission("job$delete"), localmiddleware.OperLog(logWriter, "Job-[定时任务]", consts.DELETE), h.BatchDelete)
		// 手动触发
		v1.POST("/:id/trigger", middleware.HasPermission("job$trigger"), h.Trigger)
		// 启动/暂停
		v1.PUT("/:id/status", middleware.HasPermission("job$edit"), h.ChangeStatus)

		registerJobRoutes(v1, db, logWriter)
	}
	nodeId := fmt.Sprintf("node-%d", os.Getpid())
	scheduler.Instance().SetRedis(redis.Client, nodeId)
	tasks.Init()
	if err := svc.InitScheduler(); err != nil {
		xlog.Error("scheduler init failed: %v", err)
	}
	return r, nil
}

func registerJobRoutes(rg *gin.RouterGroup, db *gorm.DB, logWriter log_grpc.OperLogWriter) {
	logHandler := handler.NewJobLogHandler(
		service.NewJobLogService(
			repository.NewJobLogRepository(db),
		),
	)
	g := rg.Group("/log")
	{
		g.GET("/page", middleware.HasPermission("job:log$list"), logHandler.Page)
		g.GET("/:id", middleware.HasPermission("job:log$query"), logHandler.Get)
		g.DELETE("/:ids", middleware.HasPermission("job:log$delete"), localmiddleware.OperLog(logWriter, "JobLog-[执行日志]", consts.DELETE), logHandler.BatchDelete)
	}

}
