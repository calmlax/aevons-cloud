package router

import (
	"fmt"

	"aevo/internal/modules/system/handler"
	"aevo/internal/modules/system/repository"
	"aevo/internal/modules/system/service"
	"aevons-grpc/log_grpc"
	authHandler "auth-service/handler"
	credRepo "auth-service/repository"
	authService "auth-service/service"

	"github.com/calmlax/aevons-framework/auth"
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

	logWriter := log_grpc.LoginLogWriter(log_grpc.NopLoginLogWriter{})
	client, err := log_grpc.NewLoginLogClient(cfg.Consul)
	if err != nil {
		xlog.Warn("init login log grpc client from consul failed: %v", err)
	} else {
		logWriter = client
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
	r.Use(middleware.AuthMiddleware(auth.NewRedisTokenStore(redisClient), cfg.Auth.Excludes))

	v1 := r.Group("/api/v1")
	{
		RegisterRoutes(v1, db, redisClient, cfg, logWriter)
	}

	return r, nil
}

// RegisterRoutes 将认证模块的所有路由注册到指定路由组。
func RegisterRoutes(rg *gin.RouterGroup, db *gorm.DB, client *redis.Client, cfg *config.Config, logWriter log_grpc.LoginLogWriter) {
	store := auth.NewRedisTokenStore(client)
	userRepo := repository.NewUserRepository(db)
	clientRepo := repository.NewOauthClientRepository(db)
	clientSvc := service.NewOauthClientService(clientRepo)
	notifier := auth.NewHttpSLONotifier()
	svc := authService.NewAuthService(store, userRepo, clientSvc, cfg.Auth, notifier, logWriter)

	h := authHandler.NewAuthHandler(svc)

	// 公开路由，无需令牌
	rg.POST("/auth/login", h.Login)
	rg.POST("/auth/refresh", h.Refresh)
	rg.POST("/auth/email/code", h.SendEmailCode)
	rg.POST("/auth/register", h.Register)
	rg.POST("/auth/reset-password", h.ResetPassword)
	rg.GET("/auth/public-key", h.GetPublicKey)

	// OAuth2 标准授权码流程端点
	rg.GET("/auth/authorize", h.Authorize)
	rg.POST("/auth/authorize", h.ApproveAuthorize)
	rg.GET("/auth/callback", h.Callback)

	// 受保护路由，由全局 AuthMiddleware 验证令牌
	rg.POST("/auth/code", h.GenerateAuthCode)
	rg.POST("/auth/logout", h.Logout)
	rg.GET("/auth/routers", h.Routers)
	rg.GET("/auth/user", h.GetUserInfo)
	rg.GET("/auth/user/profile", h.GetProfile)
	rg.PUT("/auth/user/profile", h.UpdateProfile)
	rg.PUT("/auth/user/password", h.UpdatePassword)

	credRepository := credRepo.NewCredentialRepository(db)
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
	passkeySvc, err := authService.NewPasskeyService(rpId, rpOrigins, rpName, store, userRepo, credRepository, svc)
	if err != nil {
		xlog.Error("passkey service init failed: %v", err)
	}
	p := authHandler.NewPasskeyHandler(passkeySvc, svc)
	// Passkey 公开端点（认证流程无需 token）
	rg.POST("/auth/passkey/login/begin", p.BeginAuthentication)
	rg.POST("/auth/passkey/login/finish", p.FinishAuthentication)
	// Passkey 受保护端点（注册需要已登录）
	rg.POST("/auth/passkey/register/begin", p.BeginRegistration)
	rg.POST("/auth/passkey/register/finish", p.FinishRegistration)
	rg.GET("/auth/passkey/credentials", p.ListCredentials)
	rg.DELETE("/auth/passkey/credentials/:id", p.RevokeCredential)

	loginLogHandler := handler.NewLoginLogHandler(service.NewLoginLogService(repository.NewLoginLogRepository(db)))
	rg.GET("/auth/user/login-logs", loginLogHandler.GetProfileLoginLog)
}
