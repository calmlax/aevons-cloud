package router

import (
	"net/http"

	appconfig "github.com/calmlax/aevons-framework/config"
	"github.com/calmlax/aevons-framework/middleware"
	"github.com/calmlax/aevons-framework/xlog"

	"github.com/gin-gonic/gin"
	swagSpec "github.com/swaggo/swag"
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

	// tokenStore := pkgauth.NewRedisTokenStore(redis.Client)
	// r.Use(middleware.AuthMiddleware(tokenStore, cfg.Auth.Excludes))

	if cfg.Swagger.Enabled {
		r.GET("/apifox/openapi.json", func(c *gin.Context) {
			doc, err := swagSpec.ReadDoc()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.Data(http.StatusOK, "application/json; charset=utf-8", []byte(doc))
		})
		xlog.Info("apifox openapi:    http://localhost:%d/apifox/openapi.json", cfg.Server.Port)
	}

	// v1 := r.Group("/api/v1")

	// v1.GET("/ping", handler.Ping)
	// system.RegisterRoutes(v1, cfg)

	return r
}
