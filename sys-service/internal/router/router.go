package router

import (
	"fmt"
	"sys-service/internal/handler"
	localmiddleware "sys-service/internal/middleware"
	"sys-service/internal/repository"
	"sys-service/internal/service"

	"internal-grpc/log_grpc"

	"github.com/calmlax/aevons-framework/auth/store"
	"github.com/calmlax/aevons-framework/config"
	"github.com/calmlax/aevons-framework/consts"
	"github.com/calmlax/aevons-framework/core"
	"github.com/calmlax/aevons-framework/core/scope"
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

	logWriter := log_grpc.OperLogWriter(log_grpc.NopOperLogWriter{})
	client, err := log_grpc.NewOperLogClient(cfg.Consul)
	if err != nil {
		xlog.Warn("init oper log grpc client from consul failed: %v", err)
	} else {
		logWriter = client
	}

	gin.SetMode(cfg.Server.Mode)

	r := gin.New()
	r.Use(middleware.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.RequestID())
	server.RegisterHealthRoute(r, cfg.Server.Name)
	server.RegisterOpenApiRoute(r, cfg)
	r.Use(middleware.AuthMiddleware(store.NewRedisTokenStore(redisClient), cfg.Auth.Excludes))

	v1 := r.Group("/api/sys/v1")
	{
		v1.GET("ping", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "pong"})
		})
		registerUserRoutes(v1, db, logWriter)
		registerRoleRoutes(v1, db, logWriter)
		registerDeptRoutes(v1, db, logWriter)
		registerPostRoutes(v1, db, logWriter)
		registerConfRoutes(v1, db, logWriter)
		registerDictRoutes(v1, db, logWriter)
		registerMenuRoutes(v1, db, logWriter)
		registerLangRoutes(v1, db, logWriter)
		registerOAuthClientRoutes(v1, db, logWriter)
		registerNoticeRoutes(v1, db, logWriter)
		registerMonitorRoutes(v1, db, redisClient, cfg, logWriter)
	}

	return r, nil
}

func registerConfRoutes(rg *gin.RouterGroup, db *gorm.DB, logWriter log_grpc.OperLogWriter) {
	h := handler.NewConfHandler(
		service.NewConfService(
			repository.NewConfRepository(db),
		),
	)
	g := rg.Group("/conf")
	{
		g.GET("/list", middleware.HasPermission("sys:conf$list"), h.List)
		g.GET("/page", middleware.HasPermission("sys:conf$list"), h.Page)
		g.GET("/:id", middleware.HasPermission("sys:conf$query"), h.Get)
		g.GET("/key/:key", h.GetConfByKey)
		g.POST("", middleware.HasPermission("sys:conf$add"), localmiddleware.OperLog(logWriter, "Conf-[参数配置]", consts.INSERT), h.CreateConf)
		g.POST("/refresh-cache", middleware.HasPermission("sys:conf$refresh"), localmiddleware.OperLog(logWriter, "Conf-[参数配置]", consts.CLEAN), h.RefreshCache)
		g.PUT("/:id", middleware.HasPermission("sys:conf$edit"), localmiddleware.OperLog(logWriter, "Conf-[参数配置]", consts.UPDATE), h.UpdateConf)
		g.DELETE("/:id", middleware.HasPermission("sys:conf$delete"), localmiddleware.OperLog(logWriter, "Conf-[参数配置]", consts.DELETE), h.Delete)
	}
}

func registerUserRoutes(rg *gin.RouterGroup, db *gorm.DB, logWriter log_grpc.OperLogWriter) {
	dictSvc := service.NewDictDataService(repository.NewDictDataRepository(db))
	h := handler.NewUserHandler(service.NewUserService(repository.NewUserRepository(db)), dictSvc)
	g := rg.Group("/user")
	{
		dataScope := middleware.DataScope(scope.DataScope{
			UserAlias:  "sys_user",
			UserIdName: "id",
			Resolver:   scope.DefaultDataScopeResolver(),
		})
		g.GET("/list", middleware.HasPermission("sys:user$list"), dataScope, h.List)
		g.GET("/page", middleware.HasPermission("sys:user$list"), dataScope, h.Page)
		g.GET("/export", middleware.HasPermission("sys:user$export"), dataScope, h.Export)
		g.GET("/import/template", middleware.HasPermission("sys:user$import"), h.ImportTemplate)
		g.POST("/import", middleware.HasPermission("sys:user$import"), localmiddleware.OperLog(logWriter, "User-[用户导入]", consts.IMPORT), h.Import)
		g.GET("/:id", middleware.HasPermission("sys:user$query"), dataScope, h.Get)
		g.GET("/:id/relations", middleware.HasPermission("sys:user$query"), h.GetRelations)
		g.POST("", middleware.HasPermission("sys:user$add"), localmiddleware.OperLog(logWriter, "User-[用户管理]", consts.INSERT), h.Create)
		g.PUT("/:id", middleware.HasPermission("sys:user$edit"), dataScope, localmiddleware.OperLog(logWriter, "User-[用户管理]", consts.UPDATE), h.Update)
		g.PUT("/:id/status", middleware.HasPermission("sys:user$edit"), dataScope, localmiddleware.OperLog(logWriter, "User-[用户管理]", consts.UPDATE), h.UpdateStatus)
		g.PUT("/:id/reset-password", middleware.HasPermission("sys:user$edit"), dataScope, localmiddleware.OperLog(logWriter, "User-[用户管理]", consts.UPDATE), h.ResetPassword)
		g.DELETE("/:ids", middleware.HasPermission("sys:user$delete"), dataScope, localmiddleware.OperLog(logWriter, "User-[用户管理]", consts.DELETE), h.BatchDelete)
	}
}

func registerRoleRoutes(rg *gin.RouterGroup, db *gorm.DB, logWriter log_grpc.OperLogWriter) {
	h := handler.NewRoleHandler(
		service.NewRoleService(
			repository.NewRoleRepository(db),
		),
	)
	g := rg.Group("/role")
	{
		g.GET("/list", h.List)
		g.GET("/page", h.Page)
		g.GET("/:id", middleware.HasPermission("sys:role$query"), h.Get)
		g.GET("/:id/menu", middleware.HasPermission("sys:role$query"), h.GetMenuIds)
		g.GET("/:id/dept", middleware.HasPermission("sys:role$query"), h.GetDeptIds)
		g.POST("", middleware.HasPermission("sys:role$add"), localmiddleware.OperLog(logWriter, "Role-[角色信息表]", consts.INSERT), h.Create)
		g.PUT("/:id", middleware.HasPermission("sys:role$edit"), localmiddleware.OperLog(logWriter, "Role-[角色信息表]", consts.UPDATE), h.Update)
		g.DELETE("/:ids", middleware.HasPermission("sys:role$delete"), localmiddleware.OperLog(logWriter, "Role-[角色信息表]", consts.DELETE), h.BatchDelete)
	}
}

func registerDeptRoutes(rg *gin.RouterGroup, db *gorm.DB, logWriter log_grpc.OperLogWriter) {
	h := handler.NewDeptHandler(
		service.NewDeptService(
			repository.NewDeptRepository(db),
		),
	)
	g := rg.Group("/dept")
	{
		g.GET("/list", h.ListTree)
		g.GET("/:id", middleware.HasPermission("sys:dept$query"), h.Get)
		g.POST("", middleware.HasPermission("sys:dept$add"), localmiddleware.OperLog(logWriter, "Dept-[部门管理]", consts.INSERT), h.Create)
		g.PUT("/:id", middleware.HasPermission("sys:dept$edit"), localmiddleware.OperLog(logWriter, "Dept-[部门管理]", consts.UPDATE), h.Update)
		g.DELETE("/:id", middleware.HasPermission("sys:dept$delete"), localmiddleware.OperLog(logWriter, "Dept-[部门管理]", consts.DELETE), h.Delete)
	}
}

func registerPostRoutes(rg *gin.RouterGroup, db *gorm.DB, logWriter log_grpc.OperLogWriter) {
	h := handler.NewPostHandler(
		service.NewPostService(
			repository.NewPostRepository(db),
		),
	)
	g := rg.Group("/post")
	{
		g.GET("/list", h.List)
		g.GET("/page", h.Page)
		g.GET("/:id", middleware.HasPermission("sys:post$query"), h.Get)
		g.POST("", middleware.HasPermission("sys:post$add"), localmiddleware.OperLog(logWriter, "Post-[岗位信息表]", consts.INSERT), h.Create)
		g.PUT("/:id", middleware.HasPermission("sys:post$edit"), localmiddleware.OperLog(logWriter, "Post-[岗位信息表]", consts.UPDATE), h.Update)
		g.DELETE("/:ids", middleware.HasPermission("sys:post$delete"), localmiddleware.OperLog(logWriter, "Post-[岗位信息表]", consts.DELETE), h.BatchDelete)
	}
}

func registerDictRoutes(rg *gin.RouterGroup, db *gorm.DB, logWriter log_grpc.OperLogWriter) {
	dds := service.NewDictDataService(repository.NewDictDataRepository(db))
	h := handler.NewDictHandler(service.NewDictService(repository.NewDictRepository(db)), dds)
	ddh := handler.NewDictDataHandler(dds)
	g := rg.Group("/dict")
	{
		g.GET("/list", h.AvailableList)
		g.GET("/page", middleware.HasPermission("sys:dict$list"), h.Page)
		g.GET("/:id", middleware.HasPermission("sys:dict$query"), h.Get)
		g.POST("", middleware.HasPermission("sys:dict$add"), localmiddleware.OperLog(logWriter, "Dict-[字典类型]", consts.INSERT), h.CreateDict)
		g.PUT("/:id", middleware.HasPermission("sys:dict$edit"), localmiddleware.OperLog(logWriter, "Dict-[字典类型]", consts.UPDATE), h.UpdateDict)
		g.DELETE("/:id", middleware.HasPermission("sys:dict$delete"), localmiddleware.OperLog(logWriter, "Dict-[字典类型]", consts.DELETE), h.DeleteDict)
		g.DELETE("/refresh-cache", middleware.HasPermission("sys:dict$refresh"), ddh.RefreshCache)
		g.GET("/type/:id", ddh.GetDictDataCache)
		data := g.Group("/data", middleware.HasPermission("sys:dict$design"))
		{
			data.GET("/list", ddh.ListByDictType)
			data.GET("/:id", ddh.GetDetail)
			data.POST("", localmiddleware.OperLog(logWriter, "Dict-[字典数据]", consts.INSERT), ddh.CreateDictData)
			data.PUT("/sort", localmiddleware.OperLog(logWriter, "Dict-[字典数据]", consts.UPDATE), ddh.UpdateSort)
			data.PUT("/:id", localmiddleware.OperLog(logWriter, "Dict-[字典数据]", consts.UPDATE), ddh.UpdateDictData)
			data.DELETE("/:ids", localmiddleware.OperLog(logWriter, "Dict-[字典数据]", consts.DELETE), ddh.BatchDelete)
		}
	}
}

func registerMenuRoutes(rg *gin.RouterGroup, db *gorm.DB, logWriter log_grpc.OperLogWriter) {
	h := handler.NewMenuHandler(
		service.NewMenuService(
			repository.NewMenuRepository(db),
		),
	)
	g := rg.Group("/menu")
	{
		g.GET("/list", middleware.HasPermission("sys:menu$list"), h.ListByLangCode)
		g.GET("/:id", middleware.HasPermission("sys:menu$query"), h.GetDetail)
		g.POST("", middleware.HasPermission("sys:menu$add"), localmiddleware.OperLog(logWriter, "Menu-[菜单权限表]", consts.INSERT), h.CreateMenu)
		g.PUT("/:id", middleware.HasPermission("sys:menu$edit"), localmiddleware.OperLog(logWriter, "Menu-[菜单权限表]", consts.UPDATE), h.UpdateMenu)
		g.DELETE("/:ids", middleware.HasPermission("sys:menu$delete"), localmiddleware.OperLog(logWriter, "Menu-[菜单权限表]", consts.DELETE), h.BatchDelete)
	}
}

func registerLangRoutes(rg *gin.RouterGroup, db *gorm.DB, logWriter log_grpc.OperLogWriter) {
	h := handler.NewLangHandler(
		service.NewLangService(
			repository.NewLangRepository(db),
		),
	)
	rh := handler.NewLangResourceHandler(
		service.NewLangResourceService(
			repository.NewLangResourceRepository(db),
		),
	)
	g := rg.Group("/lang")
	{
		g.GET("", h.AvailableList)
		g.GET("/list", middleware.HasPermission("sys:lang$list"), h.List)
		g.GET("/page", middleware.HasPermission("sys:lang$list"), h.Page)
		g.GET("/:id", middleware.HasPermission("sys:lang$query"), h.Get)
		g.POST("", middleware.HasPermission("sys:lang$add"), localmiddleware.OperLog(logWriter, "Lang-[语言]", consts.INSERT), h.Create)
		g.PUT("/:id", middleware.HasPermission("sys:lang$edit"), localmiddleware.OperLog(logWriter, "Lang-[语言]", consts.UPDATE), h.Update)
		g.DELETE("/:ids", middleware.HasPermission("sys:lang$delete"), localmiddleware.OperLog(logWriter, "Lang-[语言]", consts.DELETE), h.BatchDelete)
		res := g.Group("/resource")
		{
			res.GET("/list", middleware.HasPermission("sys:lang:resource$list"), rh.List)
			res.GET("/page", middleware.HasPermission("sys:lang:resource$list"), rh.Page)
			res.GET("/keys", middleware.HasPermission("sys:lang:resource$list"), rh.GetKeysByNamespace)
			res.GET("/keys/page", middleware.HasPermission("sys:lang:resource$list"), rh.PageKeys)
			res.GET("/translations", middleware.HasPermission("sys:lang:resource$list"), rh.GetTranslations)
			res.GET("/:id", middleware.HasPermission("sys:lang:resource$query"), rh.Get)
			res.POST("", middleware.HasPermission("sys:lang:resource$add"), localmiddleware.OperLog(logWriter, "LangResource-[语言资源]", consts.INSERT), rh.Create)
			res.POST("/save-translations", middleware.HasPermission("sys:lang:resource$edit"), localmiddleware.OperLog(logWriter, "LangResource-[批量保存翻译]", consts.UPDATE), rh.SaveTranslations)
			res.PUT("/:id", middleware.HasPermission("sys:lang:resource$edit"), localmiddleware.OperLog(logWriter, "LangResource-[语言资源]", consts.UPDATE), rh.Update)
			res.DELETE("/:ids", middleware.HasPermission("sys:lang:resource$delete"), localmiddleware.OperLog(logWriter, "LangResource-[语言资源]", consts.DELETE), rh.BatchDelete)
		}
	}
}

func registerOAuthClientRoutes(rg *gin.RouterGroup, db *gorm.DB, logWriter log_grpc.OperLogWriter) {
	h := handler.NewOauthClientHandler(
		service.NewOauthClientService(
			repository.NewOauthClientRepository(db),
		),
	)
	g := rg.Group("/oauth/client")
	{
		g.GET("/list", middleware.HasPermission("sys:oauth:client$list"), h.List)
		g.GET("/page", middleware.HasPermission("sys:oauth:client$list"), h.Page)
		g.GET("/:id", middleware.HasPermission("sys:oauth:client$query"), h.Get)
		g.POST("", middleware.HasPermission("sys:oauth:client$add"), localmiddleware.OperLog(logWriter, "OAuthClient-[终端应用]", consts.INSERT), h.CreateOAuthClient)
		g.POST("/refresh-cache", middleware.HasPermission("sys:oauth:client$edit"), localmiddleware.OperLog(logWriter, "OAuthClient-[终端应用缓存刷新]", consts.CLEAN), h.RefreshGatewayCache)
		g.PUT("/:id", middleware.HasPermission("sys:oauth:client$edit"), localmiddleware.OperLog(logWriter, "OAuthClient-[终端应用]", consts.UPDATE), h.UpdateOAuthClient)
		g.DELETE("/:ids", middleware.HasPermission("sys:oauth:client$delete"), localmiddleware.OperLog(logWriter, "OAuthClient-[终端应用]", consts.DELETE), h.BatchDelete)
	}
}

// 在模块 router.go 的 RegisterRoutes 中最后添加 registerNoticeRoutes(rg, db, logWriter) 代码
func registerNoticeRoutes(rg *gin.RouterGroup, db *gorm.DB, logWriter log_grpc.OperLogWriter) {
	h := handler.NewNoticeHandler(
		service.NewNoticeService(
			repository.NewNoticeRepository(db),
		),
	)
	g := rg.Group("/notice")
	{
		g.GET("/list", middleware.HasPermission("sys:notice$list"), h.List)
		g.GET("/page", middleware.HasPermission("sys:notice$list"), h.Page)
		g.GET("/:id", middleware.HasPermission("sys:notice$query"), h.Get)
		g.POST("", middleware.HasPermission("sys:notice$add"), localmiddleware.OperLog(logWriter, "Notice-[通知公告]", consts.INSERT), h.Create)
		g.PUT("/:id", middleware.HasPermission("sys:notice$edit"), localmiddleware.OperLog(logWriter, "Notice-[通知公告]", consts.UPDATE), h.Update)
		g.DELETE("/:ids", middleware.HasPermission("sys:notice$delete"), localmiddleware.OperLog(logWriter, "Notice-[通知公告]", consts.DELETE), h.BatchDelete)
	}
}

func registerMonitorRoutes(rg *gin.RouterGroup, db *gorm.DB, redisClient *redis.Client, cfg *config.Config, logWriter log_grpc.OperLogWriter) {
	m := rg.Group("monitor")
	{
		monitorH := handler.NewMonitorHandler(db, redisClient, cfg)
		m.GET("/server", middleware.HasPermission("monitor:server$query"), monitorH.GetServerInfo)
		m.GET("/online", middleware.HasPermission("monitor:online$list"), monitorH.ListOnline)
		m.DELETE("/online/:token", middleware.HasPermission("monitor:online$forceLogout"), localmiddleware.OperLog(logWriter, "Online-[在线用户]", consts.DELETE), monitorH.ForceLogout)
		m.GET("/cache", middleware.HasPermission("cache$list"), monitorH.List)
		m.GET("/cache/detail", middleware.HasPermission("cache$list"), monitorH.Detail)
		m.DELETE("/cache", middleware.HasPermission("cache$delete"), localmiddleware.OperLog(logWriter, "Redis-[缓存管理]", consts.DELETE), monitorH.Delete)
		m.DELETE("/cache/prefix", middleware.HasPermission("cache$delete"), localmiddleware.OperLog(logWriter, "Redis-[缓存管理]", consts.DELETE), monitorH.DeleteByPrefix)
	}
}
