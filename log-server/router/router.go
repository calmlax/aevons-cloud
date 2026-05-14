package router

import (
	"fmt"
	handler "log-server/hander"
	"log-server/repository"
	"log-server/service"

	"github.com/calmlax/aevons-framework/auth"
	"github.com/calmlax/aevons-framework/consts"
	"github.com/calmlax/aevons-framework/core"
	"github.com/calmlax/aevons-framework/core/server"
	"github.com/calmlax/aevons-framework/middleware"
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

	gin.SetMode(cfg.Server.Mode)

	r := gin.New()
	r.Use(middleware.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.RequestID())
	r.Use(middleware.CORS(cfg.CORS.Enabled, cfg.CORS.AllowedOrigins))
	r.Use(middleware.XSSMiddleware(cfg))
	server.RegisterHealthRoute(r, cfg.Server.Name)
	server.RegisterOpenApiRoute(r, cfg)
	r.Use(middleware.AuthMiddleware(auth.NewRedisTokenStore(redisClient), cfg.Auth.Excludes))

	v1 := r.Group("/api/v1")
	{
		registerLoginLogRoutes(v1, db)
		registerOperLogRoutes(v1, db)
	}

	return r, nil
}

func registerLoginLogRoutes(rg *gin.RouterGroup, db *gorm.DB) {
	h := handler.NewLoginLogHandler(service.NewLoginLogService(repository.NewLoginLogRepository(db)))
	g := rg.Group("/login/log")
	{
		g.GET("/list", middleware.HasPermission("monitor:login:log$query"), h.List)
		g.GET("/page", middleware.HasPermission("monitor:login:log$query"), h.Page)
		g.GET("/:id", middleware.HasPermission("monitor:login:log$query"), h.Get)
		g.DELETE("/:ids", middleware.HasPermission("monitor:log$delete"), middleware.OperLog(db, "LoginLog-[登录日志]", consts.DELETE), h.BatchDelete)
		g.DELETE("", middleware.HasPermission("monitor:login:log$clear"), middleware.OperLog(db, "LoginLog-[登录日志清空]", consts.CLEAN), h.Clear)
	}
}

func registerOperLogRoutes(rg *gin.RouterGroup, db *gorm.DB) {
	h := handler.NewOperLogHandler(service.NewOperLogService(repository.NewOperLogRepository(db)))
	g := rg.Group("/oper/log")
	{
		g.GET("/list", middleware.HasPermission("monitor:oper:log$query"), h.List)
		g.GET("/page", middleware.HasPermission("monitor:oper:log$query"), h.Page)
		g.GET("/:id", middleware.HasPermission("monitor:oper:log$query"), h.Get)
		g.DELETE("/:ids", middleware.HasPermission("monitor:oper:log$delete"), middleware.OperLog(db, "OperLog-[操作日志]", consts.DELETE), h.BatchDelete)
		g.DELETE("", middleware.HasPermission("monitor:oper:log$clear"), middleware.OperLog(db, "OperLog-[操作日志清空]", consts.CLEAN), h.Clear)
	}
}
