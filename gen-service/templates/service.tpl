package service

/**
 * {{.Comment}} Service
 *
 * @author {{.Author}}
 * @date {{.Date}}
 * @website https://www.aevons.com
 * @email 1073602@qq.com
 * @copyright 2025-2026 Aevons
 */

import (
	"aevo/internal/modules/{{.ModuleName}}/dto"
	"aevo/internal/modules/{{.ModuleName}}/model"
	"aevo/internal/modules/{{.ModuleName}}/repository"
	"aevo/pkg/base"
)

// 继承BaseService
type {{.ClassName}}Service interface {
	base.BaseService[model.{{.ClassName}}, *dto.{{.ClassName}}Query]
}

type {{.ClassName | toLowerCamel}}Service struct {
	base.BaseService[model.{{.ClassName}}, *dto.{{.ClassName}}Query]
	repo repository.{{.ClassName}}Repository
}

func New{{.ClassName}}Service(repo repository.{{.ClassName}}Repository) {{.ClassName}}Service {
	baseSrv := base.NewBaseService[model.{{.ClassName}}, *dto.{{.ClassName}}Query](repo)
	return &{{.ClassName | toLowerCamel}}Service{
		BaseService: baseSrv,
		repo:        repo,
	}
}
