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
	"{{.ModuleName}}-service/internal/dto"
	"{{.ModuleName}}-service/internal/model"
	"{{.ModuleName}}-service/internal/service"
	"github.com/calmlax/aevons-framework/core/base"
)

type {{.ClassName}}Handler struct {
	crud *base.BaseHandler[model.{{.ClassName}}, *dto.{{.ClassName}}Query, dto.Create{{.ClassName}}DTO, dto.Update{{.ClassName}}DTO]
	svc service.{{.ClassName}}Service
}

// 构造函数
func New{{.ClassName}}Handler(svc service.{{.ClassName}}Service) *{{.ClassName}}Handler {
	return &{{.ClassName}}Handler{
		crud: base.NewBaseHandler[model.{{.ClassName}},, *dto.{{.ClassName}},Query, dto.Create{{.ClassName}},DTO, dto.Update{{.ClassName}},DTO](svc),
		svc: svc,
	}
}


// List 查询{{.Comment}}列表。
//
// @Summary      查询{{.Comment}}列表
// @Tags         {{.Comment}}
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  response.Response
// @Router       /api/{{.ModuleName}}/v1/{{.Router}}/list [get]
func (h *{{.ClassName}}Handler) List(c *gin.Context) {
	h.crud.HandleList(c)
}

// Page 分页查询{{.Comment}}。
//
// @Summary      分页查询{{.Comment}}
// @Tags         {{.Comment}}
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  response.Response
// @Router       /api/{{.ModuleName}}/v1/{{.Router}}/page [get]
func (h *{{.ClassName}}Handler) Page(c *gin.Context) {
	h.crud.HandlePage(c)
}

// Get 获取{{.Comment}}详情。
//
// @Summary      获取{{.Comment}}详情
// @Tags         {{.Comment}}
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "{{.Comment}}ID"
// @Success      200  {object}  response.Response
// @Router       /api/{{.ModuleName}}/v1/{{.Router}}/{id} [get]
func (h *{{.ClassName}}Handler) Get(c *gin.Context) {
	h.crud.HandleGet(c)
}

// Create 新增{{.Comment}}
//
// @Summary      新增{{.Comment}}
// @Tags         {{.Comment}}
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      dto.Create{{.ClassName}}DTO  true  "新增{{.Comment}}"
// @Success      200      {object}  response.Response
// @Router       /api/{{.ModuleName}}/v1/{{.Router}} [post]
func (h *{{.ClassName}}Handler) Create(c *gin.Context) {
	h.crud.HandleCreate(c)
}

// Update 修改{{.Comment}}
//
// @Summary      修改{{.Comment}}
// @Tags         {{.Comment}}
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      int                        true  "{{.Comment}}ID"
// @Param        request  body      dto.Update{{.ClassName}}DTO  true  "修改{{.Comment}}"
// @Success      200      {object}  response.Response
// @Router       /api/{{.ModuleName}}/v1/{{.Router}}/{id} [put]
func (h *{{.ClassName}}Handler) Update(c *gin.Context) {
	h.crud.HandleUpdate(c)
}

// BatchDelete 批量删除{{.Comment}}。
//
// @Summary      批量删除{{.Comment}}
// @Tags         {{.Comment}}
// @Produce      json
// @Security     BearerAuth
// @Param        ids   path      string  true  "逗号分隔的{{.Comment}}ID"
// @Success      200   {object}  response.Response
// @Router       /api/{{.ModuleName}}/v1/{{.Router}}/{ids} [delete]
func (h *{{.ClassName}}Handler) BatchDelete(c *gin.Context) {
	h.crud.HandleBatchDelete(c)
}
