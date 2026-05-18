package {{.ModuleName}}

/**
 * {{.Comment}} Router
 *
 * @author {{.Author}}
 * @date {{.Date}}
 * @website https://www.aevons.com
 * @email 1073602@qq.com
 * @copyright 2025-2026 Aevons
 */

import (
	localmiddleware "sys-service/internal/middleware"
	"github.com/calmlax/aevons-framework/middleware"
	"{{.ModuleName}}-service/internal/handler"
	"{{.ModuleName}}-service/internal/repository"
	"{{.ModuleName}}-service/internal/service"
	"github.com/calmlax/aevons-framework/consts"
	"internal-grpc/log_grpc"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 在模块 router.go 的 RegisterRoutes 中最后添加 register{{.ClassName}}Routes(rg, db) 代码
func register{{.ClassName}}Routes(rg *gin.RouterGroup, db *gorm.DB, logWriter log_grpc.OperLogWriter) {
	h := handler.New{{.ClassName}}Handler(
		service.New{{.ClassName}}Service(
			repository.New{{.ClassName}}Repository(db),
		),
	)
	g := rg.Group("/{{.Router}}")
	{
		g.GET("/list", middleware.HasPermission("{{.Permission}}$list"), h.List)
		g.GET("/page", middleware.HasPermission("{{.Permission}}$list"), h.Page)
		g.GET("/:id", middleware.HasPermission("{{.Permission}}$query"), h.Get)
		g.POST("", middleware.HasPermission("{{.Permission}}$add"), localmiddleware.OperLog(logWriter, "{{.ClassName}}-[{{.Comment}}]", consts.INSERT), h.Create)
		g.PUT("/:id", middleware.HasPermission("{{.Permission}}$edit"), localmiddleware.OperLog(logWriter, "{{.ClassName}}-[{{.Comment}}]", consts.UPDATE), h.Update)
		g.DELETE("/:ids", middleware.HasPermission("{{.Permission}}$delete"), localmiddleware.OperLog(logWriter, "{{.ClassName}}-[{{.Comment}}]", consts.DELETE), h.BatchDelete)
	}
}
