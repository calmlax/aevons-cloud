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
	"aevo/internal/middleware"
	"aevo/internal/modules/{{.ModuleName}}/handler"
	"aevo/internal/modules/{{.ModuleName}}/repository"
	"aevo/internal/modules/{{.ModuleName}}/service"
	consts "aevo/pkg/const"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 在模块 router.go 的 RegisterRoutes 中最后添加 register{{.ClassName}}Routes(rg, db) 代码
func register{{.ClassName}}Routes(rg *gin.RouterGroup, db *gorm.DB) {
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
		g.POST("", middleware.HasPermission("{{.Permission}}$add"), middleware.OperLog(db, "{{.ClassName}}-[{{.Comment}}]", consts.INSERT), h.Create)
		g.PUT("/:id", middleware.HasPermission("{{.Permission}}$edit"), middleware.OperLog(db, "{{.ClassName}}-[{{.Comment}}]", consts.UPDATE), h.Update)
		g.DELETE("/:ids", middleware.HasPermission("{{.Permission}}$delete"), middleware.OperLog(db, "{{.ClassName}}-[{{.Comment}}]", consts.DELETE), h.BatchDelete)
	}
}
