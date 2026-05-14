package router

import (
	"fmt"

	"github.com/calmlax/aevons-framework/auth"
	"github.com/calmlax/aevons-framework/core"
	"github.com/calmlax/aevons-framework/core/server"
	"github.com/calmlax/aevons-framework/middleware"

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

	gin.SetMode(cfg.Server.Mode)

	r := gin.New()
	r.Use(middleware.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.RequestID())
	r.Use(middleware.CORS(cfg.CORS.Enabled, cfg.CORS.AllowedOrigins))
	r.Use(middleware.XSSMiddleware(cfg))
	server.RegisterHealthRoute(r, cfg.Server.Name)
	server.RegisterOpenApiRoute(r, cfg)

	tokenStore := auth.NewRedisTokenStore(redisClient)
	r.Use(middleware.AuthMiddleware(tokenStore, cfg.Auth.Excludes))

	v1 := r.Group("/api/v1")
	{
		registerConfRoutes(v1, app)
	}

	return r, nil
}

func registerConfRoutes(rg *gin.RouterGroup, app *core.App) {
	// h := handler.NewConfHandler(service.NewConfService(repository.NewConfRepository(db)))
	// g := rg.Group("/conf")
	// {
	// 	g.GET("/list", middleware.HasPermission("sys:conf$list"), h.List)
	// 	g.GET("/page", middleware.HasPermission("sys:conf$list"), h.Page)
	// 	g.GET("/:id", middleware.HasPermission("sys:conf$query"), h.Get)
	// 	g.GET("/key/:key", h.GetConfByKey)
	// 	g.POST("", middleware.HasPermission("sys:conf$add"), middleware.OperLog(db, "参数配置", consts.INSERT), h.CreateConf)
	// 	g.POST("/refresh-cache", middleware.HasPermission("sys:conf$edit"), middleware.OperLog(db, "参数配置", consts.CLEAN), h.RefreshCache)
	// 	g.PUT("/:id", middleware.HasPermission("sys:conf$edit"), middleware.OperLog(db, "参数配置", consts.UPDATE), h.UpdateConf)
	// 	g.DELETE("/:id", middleware.HasPermission("sys:conf$delete"), middleware.OperLog(db, "参数配置", consts.DELETE), h.Delete)
	// }
	_ = rg
	_ = app
}
