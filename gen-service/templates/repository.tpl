package repository

/**
 * {{.Comment}} Repository
 *
 * @author {{.Author}}
 * @date {{.Date}}
 * @website https://www.aevons.com
 * @email 1073602@qq.com
 * @copyright 2025-2026 Aevons
 */

import (
	"{{.ModuleName}}-service/internal/model"
	"github.com/calmlax/aevons-framework/core/base"

	"gorm.io/gorm"
)

// 继承BaseRepository
type {{.ClassName}}Repository interface {
	base.BaseRepository[model.{{.ClassName}}]
}

// {{.ClassName | toLowerCamel}}Repository 结构体
type {{.ClassName | toLowerCamel}}Repository struct {
	base.BaseRepository[model.{{.ClassName}}]
	db *gorm.DB
}

// 创建 {{.ClassName}}Repository 实例
func New{{.ClassName}}Repository(db *gorm.DB) {{.ClassName}}Repository {
	return &{{.ClassName | toLowerCamel}}Repository{
		BaseRepository: base.NewBaseRepository[model.{{.ClassName}}](db),
		db:             db,
	}
}
