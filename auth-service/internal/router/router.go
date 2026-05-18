package router

import (
	"fmt"

	"aevons-grpc/log_grpc"
	authHandler "auth-service/internal/handler"
	authRepo "auth-service/internal/repository"
	authService "auth-service/internal/service"

	authnotifier "github.com/calmlax/aevons-framework/auth/notifier"
	authstore "github.com/calmlax/aevons-framework/auth/store"
	"github.com/calmlax/aevons-framework/config"
	"github.com/calmlax/aevons-framework/core"
	"github.com/calmlax/aevons-framework/core/server"
	"github.com/calmlax/aevons-framework/middleware"
	"github.com/calmlax/aevons-framework/xlog"
	"github.com/redis/go-redis/v9"
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

	logStore := log_grpc.LoginLogStore(log_grpc.NopLoginLogStore{})
	client, err := log_grpc.NewLoginLogClient(cfg.Consul)
	if err != nil {
		xlog.Warn("init login log grpc client from consul failed: %v", err)
	} else {
		logStore = client
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
	r.Use(middleware.AuthMiddleware(authstore.NewRedisTokenStore(redisClient), cfg.Auth.Excludes))

	v1 := r.Group("/api/v1/auth")
	{
		v1.GET("ping", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "pong"})
		})
		RegisterRoutes(v1, db, redisClient, cfg, logStore)
	}

	return r, nil
}

// RegisterRoutes 将认证模块的所有路由注册到指定路由组。
func RegisterRoutes(rg *gin.RouterGroup, db *gorm.DB, client *redis.Client, cfg *config.Config, logStore log_grpc.LoginLogStore) {
	store := authstore.NewRedisTokenStore(client)
	repo := authRepo.NewAuthRepository(db)
	notifier := authnotifier.NewHttpSLONotifier()
	svc := authService.NewAuthService(store, repo, cfg.Auth, notifier, logStore)

	h := authHandler.NewAuthHandler(svc)

	// 公开路由，无需令牌
	rg.POST("/login", h.Login)
	rg.POST("/refresh", h.Refresh)
	rg.POST("/email/code", h.SendEmailCode)
	rg.POST("/register", h.Register)
	rg.POST("/reset-password", h.ResetPassword)
	rg.GET("/public-key", h.GetPublicKey)

	// OAuth2 标准授权码流程端点
	rg.GET("/authorize", h.Authorize)
	rg.POST("/authorize", h.ApproveAuthorize)
	rg.GET("/callback", h.Callback)

	// 受保护路由，由全局 AuthMiddleware 验证令牌
	rg.POST("/code", h.GenerateAuthCode)
	rg.POST("/logout", h.Logout)
	rg.GET("/routers", h.Routers)
	rg.GET("/user", h.GetUserInfo)
	rg.GET("/user/profile", h.GetProfile)
	rg.PUT("/user/profile", h.UpdateProfile)
	rg.PUT("/user/password", h.UpdatePassword)

	rpId := cfg.WebAuthn.RPID
	if rpId == "" {
		rpId = "localhost"
	}
	rpOrigins := cfg.WebAuthn.RPOrigins
	if len(rpOrigins) == 0 {
		rpOrigins = []string{"http://localhost:5173"}
	}
	rpName := cfg.WebAuthn.RPName
	if rpName == "" {
		rpName = "Aevons Admin"
	}
	passkeySvc, err := authService.NewPasskeyService(rpId, rpOrigins, rpName, store, repo, svc)
	if err != nil {
		xlog.Error("passkey service init failed: %v", err)
	}
	p := authHandler.NewPasskeyHandler(passkeySvc, svc)
	// Passkey 公开端点（认证流程无需 token）
	rg.POST("/passkey/login/begin", p.BeginAuthentication)
	rg.POST("/passkey/login/finish", p.FinishAuthentication)
	// Passkey 受保护端点（注册需要已登录）
	rg.POST("/passkey/register/begin", p.BeginRegistration)
	rg.POST("/passkey/register/finish", p.FinishRegistration)
	rg.GET("/passkey/credentials", p.ListCredentials)
	rg.DELETE("/passkey/credentials/:id", p.RevokeCredential)

	rg.GET("/user/login-logs", h.GetLatestLoginLog)
}
