package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	consoleconfig "gateway-console/internal/config"
	"gateway-console/internal/router"

	"github.com/calmlax/aevons-framework/config"
	"github.com/calmlax/aevons-framework/core"
	"github.com/calmlax/aevons-framework/core/consul"
	"github.com/calmlax/aevons-framework/core/xjson"
	"github.com/calmlax/aevons-framework/core/xlog"
	"github.com/gin-gonic/gin/binding"
)

func init() {
	binding.EnableDecoderUseNumber = true
	xjson.InitGin()
}

func main() {
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
		xlog.Fatal("failed to load framework config: %v", err)
	}
	if err := xlog.Init(&cfg); err != nil {
		xlog.Fatal("failed to init logger: %v", err)
	}

	consoleCfg, err := consoleconfig.Load(configPath, env)
	if err != nil {
		xlog.Fatal("failed to load console config: %v", err)
	}

	app := core.NewApp(&cfg, nil, nil)

	engine, err := router.Setup(app, consoleCfg)
	if err != nil {
		xlog.Fatal("failed to setup router: %v", err)
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: engine,
	}

	listener, err := net.Listen("tcp", srv.Addr)
	if err != nil {
		xlog.Fatal("failed to listen on %s: %v", srv.Addr, err)
	}

	var consulManager *consul.Managed

	go func() {
		xlog.Info("gateway console starting on port %d", cfg.Server.Port)
		if err := srv.Serve(listener); err != nil && err != http.ErrServerClosed {
			xlog.Error("gateway console failed to start: %v", err)
			os.Exit(1)
		}
	}()

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
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	xlog.Info("shutting down gateway console...")

	if consulManager != nil {
		if err := consulManager.Deregister(); err != nil {
			xlog.Error("failed to deregister service from consul: %v", err)
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		xlog.Error("gateway console forced to shutdown: %v", err)
		os.Exit(1)
	}

	xlog.Info("gateway console exited gracefully")
}
