// @title           Aevons Admin API
// @version         1.0
// @description     Aevons 后台管理系统接口文档，供 Apifox 同步使用
// @host            localhost:8021
// @BasePath        /api/v1
// @schemes         http https

// @securityDefinitions.apikey  BearerAuth
// @in                          header
// @name                        Authorization
// @description                 Bearer token，格式：Bearer {token}

package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"system-server/router"
	"time"

	"github.com/calmlax/aevons-framework/config"
	"github.com/calmlax/aevons-framework/xjson"
	"github.com/calmlax/aevons-framework/xlog"

	"github.com/gin-gonic/gin/binding"
)

func init() {
	binding.EnableDecoderUseNumber = true
	xjson.InitGin()
}

func main() {
	// --config flag allows overriding the config file path at runtime.
	// Defaults to configs; APP_ENV env var selects the overlay file.
	// Examples:
	//   ./main                                    -> configs/config.yaml
	//   APP_ENV=production ./main                 -> configs/config.yaml + configs/config.production.yaml
	//   ./main --config configs       -> explicit path
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
	cfg, err := config.Load(configPath, env)
	if err != nil {
		xlog.Fatal("failed to load config: %v", err)
	}

	// if err := db.Init(cfg); err != nil {
	// 	xlog.Fatal("failed to init db: %v", err)
	// }

	// if err := redis.Init(cfg); err != nil {
	// 	xlog.Fatal("failed to init redis: %v", err)
	// }

	// if err := rocketmq.Init(cfg); err != nil {
	// 	xlog.Fatal("failed to init rocketmq: %v", err)
	// }

	// if len(os.Args) > 1 {
	// 	cmd.Execute(cfg.Generator, db.DB)
	// 	return
	// }

	// if err := db.DB.AutoMigrate(
	// 	&model.User{},
	// 	&model.Role{},
	// 	&model.Dept{},
	// 	&model.Post{},
	// 	&model.Conf{},
	// 	&model.Dict{},
	// 	&model.DictData{},
	// 	&model.Menu{},
	// 	&model.LoginLog{},
	// 	&model.OperLog{},
	// ); err != nil {
	// 	xlog.Fatal("failed to migrate db: %v", err)
	// }

	engine := router.Setup(&cfg)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: engine,
	}

	go func() {
		xlog.Info("server starting on port %d", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			xlog.Error("server failed to start: %v", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	xlog.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		xlog.Error("server forced to shutdown: %v", err)
		os.Exit(1)
	}

	// rocketmq.Close()
	xlog.Info("server exited gracefully")
}
