// @title           Aevons Auth Service API
// @version         1.0.0
// @description     Aevons 认证服务接口文档，提供登录、令牌刷新、退出登录、邮箱验证码、注册、密码重置、SSO 授权码和通行密钥认证等能力
// @host            localhost:10701
// @BasePath        /api/v1/auth
// @schemes         http https
// @securityDefinitions.apikey  BearerAuth
// @in                          header
// @name                        Authorization
// @description                 Bearer 访问令牌，格式：Bearer {token}

package main

import (
	"auth-service/internal/router"

	"github.com/calmlax/aevons-framework/core"
	"github.com/calmlax/aevons-framework/xlog"
)

func main() {
	app, err := core.Bootstrap()
	if err != nil {
		xlog.Fatal("failed to bootstrap app: %v", err)
	}

	engine, err := router.Setup(app)
	if err != nil {
		xlog.Fatal("failed to setup router: %v", err)
	}

	if err := core.RunGin(app, engine); err != nil {
		xlog.Fatal("failed to run http server: %v", err)
	}
}
