package main

import (
	"os"

	consoleconfig "gateway-console/internal/config"
	"gateway-console/internal/router"

	"github.com/calmlax/aevons-framework/core"
	"github.com/calmlax/aevons-framework/xlog"
)

func main() {
	app, err := core.BootstrapWithOptions(core.BootstrapOptions{
		InitDB:      false,
		InitRedis:   false,
		InitGinJSON: true,
	})
	if err != nil {
		xlog.Fatal("failed to bootstrap app: %v", err)
	}

	configPath := "configs"
	for i, arg := range os.Args[1:] {
		if arg == "--config" && i+1 < len(os.Args) {
			configPath = os.Args[i+1]
			break
		}
	}
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}
	consoleCfg, err := consoleconfig.Load(configPath, env)
	if err != nil {
		xlog.Fatal("failed to load console config: %v", err)
	}

	engine, err := router.Setup(app, consoleCfg)
	if err != nil {
		xlog.Fatal("failed to setup router: %v", err)
	}

	if err := core.RunGin(app, engine); err != nil {
		xlog.Fatal("failed to run gateway console: %v", err)
	}
}
