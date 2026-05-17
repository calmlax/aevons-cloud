package handler

/**
 * {{.Comment}} Handler
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
	"aevo/internal/modules/{{.ModuleName}}/service"
	"aevo/pkg/base"
)

type {{.ClassName}}Handler struct {
	// 继承BaseHandler
	*base.BaseHandler[
		model.{{.ClassName}},        // 模型
		*dto.{{.ClassName}}Query,     // 查询 DTO
		dto.Create{{.ClassName}}DTO, // 创建 DTO
		dto.Update{{.ClassName}}DTO, // 更新 DTO
	]
	svc service.{{.ClassName}}Service
}

// 构造函数
func New{{.ClassName}}Handler(svc service.{{.ClassName}}Service) *{{.ClassName}}Handler {
	return &{{.ClassName}}Handler{
		BaseHandler: base.NewBaseHandler[
			model.{{.ClassName}},
			*dto.{{.ClassName}}Query,
			dto.Create{{.ClassName}}DTO,
			dto.Update{{.ClassName}}DTO,
		](svc),
		svc: svc,
	}
}
