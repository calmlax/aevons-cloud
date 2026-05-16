package router

import (
	"fmt"

	"aevons-cloud/gateway/console/handler"
	"aevons-cloud/gateway/console/internal/apisixadmin"
	consoleconfig "aevons-cloud/gateway/console/internal/config"
	"aevons-cloud/gateway/console/repository"
	"aevons-cloud/gateway/console/service"

	"github.com/calmlax/aevons-framework/core"
	"github.com/calmlax/aevons-framework/core/server"
	"github.com/calmlax/aevons-framework/middleware"
	"github.com/gin-gonic/gin"
)

func Setup(app *core.App, consoleCfg consoleconfig.Settings) (*gin.Engine, error) {
	cfg, err := app.RawConfig()
	if err != nil {
		return nil, fmt.Errorf("router: read app config failed: %w", err)
	}

	gin.SetMode(cfg.Server.Mode)

	repo := repository.NewStaticRepository()
	catalogService := service.NewCatalogService(repo)
	apisixClient := apisixadmin.New(consoleCfg.APISIXAdminURL, consoleCfg.APISIXAdminKey)
	publishService := service.NewPublishService(catalogService, apisixClient)
	catalogHandler := handler.NewCatalogHandler(catalogService, publishService, cfg.Server.Name)

	r := gin.New()
	r.Use(middleware.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.RequestID())
	r.Use(middleware.CORS(cfg.CORS.Enabled, cfg.CORS.AllowedOrigins))
	r.Use(middleware.XSSMiddleware(cfg))
	server.RegisterHealthRoute(r, cfg.Server.Name)

	v1 := r.Group("/api/v1/gateway")
	{
		v1.GET("/overview", catalogHandler.Overview)
		v1.GET("/routes", catalogHandler.Routes)
		v1.GET("/upstreams", catalogHandler.Upstreams)
		v1.GET("/consumers", catalogHandler.Consumers)
		v1.GET("/plugins", catalogHandler.Plugins)
		v1.GET("/policies", catalogHandler.Policies)
		v1.GET("/publish/plan", catalogHandler.PublishPlan)
		v1.GET("/publish/snapshot", catalogHandler.PublishSnapshot)
		v1.POST("/publish/run", catalogHandler.PublishToAPISIX)
		v1.GET("/healthz/control-plane", catalogHandler.ControlPlaneHealth)
	}

	return r, nil
}
