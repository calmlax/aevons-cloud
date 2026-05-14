package router

import (
	"net/http"

	appconfig "github.com/calmlax/aevons-framework/config"
	"github.com/calmlax/aevons-framework/core/server"
	"github.com/calmlax/aevons-framework/middleware"

	"github.com/gin-gonic/gin"
)

// Setup 配置 Gin 引擎，注册中间件和路由。
func Setup(cfg *appconfig.Config) *gin.Engine {
	gin.SetMode(cfg.Server.Mode)

	r := gin.New()
	r.Use(middleware.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.RequestID())
	r.Use(middleware.CORS(cfg.CORS.Enabled, cfg.CORS.AllowedOrigins))
	r.Use(middleware.XSSMiddleware(cfg))
	server.RegisterOpenApiRoute(r, cfg)
	server.RegisterHealthRoute(r, cfg.Server.Name)

	// tokenStore := pkgauth.NewRedisTokenStore(redis.Client)
	// r.Use(middleware.AuthMiddleware(tokenStore, cfg.Auth.Excludes))

	v1 := r.Group("/api/v1")
	v1.GET("/test/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code":    "OK",
			"message": "pong",
			"service": cfg.Server.Name,
		})
	})

	return r
}
