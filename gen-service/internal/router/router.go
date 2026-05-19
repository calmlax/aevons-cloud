package router

import (
	"fmt"

	"internal-grpc/log_grpc"

	"gen-service/internal/handler"
	localmiddleware "gen-service/internal/middleware"
	"gen-service/internal/repository"
	"gen-service/internal/service"

	"github.com/calmlax/aevons-framework/auth/store"
	"github.com/calmlax/aevons-framework/config"
	"github.com/calmlax/aevons-framework/consts"
	"github.com/calmlax/aevons-framework/core"
	"github.com/calmlax/aevons-framework/core/server"
	"github.com/calmlax/aevons-framework/middleware"
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

	v1 := r.Group("/api/gen/v1")
	{
		v1.GET("ping", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "pong"})
		})
		registerGenTableRoutes(v1, db, logWriter, cfg.Gen)
	}

	return r, nil
}

func registerGenTableRoutes(rg *gin.RouterGroup, db *gorm.DB, logWriter log_grpc.OperLogWriter, genCfg config.GenConfig) {
	gtrepo := repository.NewGenTableRepository(db)
	gtcrepo := repository.NewGenTableColumnRepository(db)

	tableHandler := handler.NewGenTableHandler(service.NewGenTableService(gtrepo, gtcrepo, genCfg))
	columnHandler := handler.NewGenTableColumnHandler(service.NewGenTableColumnService(gtcrepo))

	g := rg.Group("/table")
	{
		g.GET("/list", middleware.HasPermission("gen:table$list"), tableHandler.List)
		g.GET("/page", middleware.HasPermission("gen:table$list"), tableHandler.Page)
		g.GET("/db", middleware.HasPermission("gen:table$import"), tableHandler.DBTables)
		g.POST("", middleware.HasPermission("gen:table$add"), localmiddleware.OperLog(logWriter, "代码生成表", consts.INSERT), tableHandler.Create)
		g.POST("/import", middleware.HasPermission("gen:table$import"), localmiddleware.OperLog(logWriter, "导入数据库表", consts.IMPORT), tableHandler.ImportTables)
		g.GET("/:id", middleware.HasPermission("gen:table$query"), tableHandler.Get)
		g.PUT("/:id", middleware.HasPermission("gen:table$edit"), localmiddleware.OperLog(logWriter, "代码生成表", consts.UPDATE), tableHandler.Update)
		g.DELETE("/:ids", middleware.HasPermission("gen:table$delete"), localmiddleware.OperLog(logWriter, "代码生成表", consts.DELETE), tableHandler.BatchDelete)
		g.GET("/download", middleware.HasPermission("gen:table$download"), localmiddleware.OperLog(logWriter, "代码生成表", consts.EXPORT), tableHandler.Download)
		g.GET("/preview", middleware.HasPermission("gen:table$preview"), tableHandler.CodePreview)
		g.GET("/synch", middleware.HasPermission("gen:table$import"), localmiddleware.OperLog(logWriter, "代码生成表", consts.SYNCH), tableHandler.SynchDbTable)

		column := g.Group("/column", middleware.HasPermission("gen:table$design"))
		{
			column.GET("/list", columnHandler.List)
			column.PUT("/batch-update", localmiddleware.OperLog(logWriter, "代码生成字段", consts.UPDATE), columnHandler.BatchUpdate)
		}
	}
}
