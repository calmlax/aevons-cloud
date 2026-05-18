// @title           Aevons Gen Service API
// @version         1.0.0
// @description     Aevons 代码生成服务接口文档，提供数据库表导入、字段设计、代码预览与代码下载等能力
// @host            localhost:10704
// @BasePath        /api/v1/gen
// @schemes         http https
// @securityDefinitions.apikey  BearerAuth
// @in                          header
// @name                        Authorization
// @description                 Bearer 访问令牌，格式：Bearer {token}

package main

import (
	"gen-service/cmd"
	"gen-service/internal/router"
	"os"

	"github.com/calmlax/aevons-framework/core"
	"github.com/calmlax/aevons-framework/db"
	"github.com/calmlax/aevons-framework/xlog"
)

func main() {
	app, err := core.Bootstrap()
	if err != nil {
		xlog.Fatal("failed to bootstrap app: %v", err)
	}

	cfg, err := app.RawConfig()
	if err != nil {
		xlog.Fatal("failed to read config from app: %v", err)
	}

	if len(os.Args) > 1 {
		cmd.Execute(cfg.Gen, db.DB)
		return
	}

	engine, err := router.Setup(app)
	if err != nil {
		xlog.Fatal("failed to setup router: %v", err)
	}

	if err := core.RunGin(app, engine); err != nil {
		xlog.Fatal("failed to run http server: %v", err)
	}
}
