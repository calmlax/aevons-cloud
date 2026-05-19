/**
 * 角色信息表 Handler
 *
 * @author
 * @date 2026-04-09 01:49:40.926495422 +0000 UTC
 * @website https://www.aevons.com
 * @email 1073602@qq.com
 * @copyright 2025-2026 Aevons
 */
package handler

import (
	"strconv"
	"sys-service/internal/dto"
	"sys-service/internal/model"
	"sys-service/internal/service"

	apperr "github.com/calmlax/aevons-framework/errors"
	"github.com/calmlax/aevons-framework/response"

	"github.com/calmlax/aevons-framework/core/base"

	"github.com/gin-gonic/gin"
)

type RoleHandler struct {
	crud *base.BaseHandler[model.Role, *dto.RoleQuery, dto.CreateRoleDTO, dto.UpdateRoleDTO]
	svc  service.RoleService
}

func NewRoleHandler(svc service.RoleService) *RoleHandler {
	return &RoleHandler{
		crud: base.NewBaseHandler[model.Role, *dto.RoleQuery, dto.CreateRoleDTO, dto.UpdateRoleDTO](svc),
		svc:  svc,
	}
}

// List 查询角色列表。
//
// @Summary      查询角色列表
// @Tags         角色管理
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  response.Response
// @Router       /api/sys/v1/role/list [get]
func (h *RoleHandler) List(c *gin.Context) {
	h.crud.HandleList(c)
}

// Page 分页查询角色。
//
// @Summary      分页查询角色
// @Tags         角色管理
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  response.Response
// @Router       /api/sys/v1/role/page [get]
func (h *RoleHandler) Page(c *gin.Context) {
	h.crud.HandlePage(c)
}

// Get 获取角色详情。
//
// @Summary      获取角色详情
// @Tags         角色管理
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "角色ID"
// @Success      200  {object}  response.Response
// @Router       /api/sys/v1/role/{id} [get]
func (h *RoleHandler) Get(c *gin.Context) {
	h.crud.HandleGet(c)
}

// BatchDelete 批量删除角色。
//
// @Summary      批量删除角色
// @Tags         角色管理
// @Produce      json
// @Security     BearerAuth
// @Param        ids   path      string  true  "逗号分隔的角色ID"
// @Success      200   {object}  response.Response
// @Router       /api/sys/v1/role/{ids} [delete]
func (h *RoleHandler) BatchDelete(c *gin.Context) {
	h.crud.HandleBatchDelete(c)
}

// Create 创建角色（含菜单，事务）
//
// @Summary      新增角色
// @Tags         角色管理
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      dto.CreateRoleDTO  true  "新增角色"
// @Success      200      {object}  response.Response
// @Router       /api/sys/v1/role [post]
func (h *RoleHandler) Create(c *gin.Context) {
	var d dto.CreateRoleDTO
	if err := c.ShouldBindJSON(&d); err != nil {
		response.FailBy(c, apperr.ErrInvalidBody)
		return
	}
	role, err := h.svc.CreateWithMenus(c, d)
	if err != nil {
		response.FailBy(c, apperr.ErrCreateFailed, map[string]any{"error": err.Error()})
		return
	}
	response.Success(c, role)
}

// Update 更新角色（含菜单，事务）
//
// @Summary      修改角色
// @Tags         角色管理
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      int                true  "角色ID"
// @Param        request  body      dto.UpdateRoleDTO  true  "修改角色"
// @Success      200      {object}  response.Response
// @Router       /api/sys/v1/role/{id} [put]
func (h *RoleHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.FailBy(c, apperr.ErrInvalidId)
		return
	}
	var d dto.UpdateRoleDTO
	if err := c.ShouldBindJSON(&d); err != nil {
		response.FailBy(c, apperr.ErrInvalidBody)
		return
	}
	if err := h.svc.UpdateWithMenus(c, id, d); err != nil {
		response.FailBy(c, apperr.ErrUpdateFailed, map[string]any{"error": err.Error()})
		return
	}
	response.Success(c, id)
}

// GetMenuIds 获取角色已关联的菜单ID列表
//
// @Summary      获取角色菜单ID列表
// @Tags         角色管理
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "角色ID"
// @Success      200  {object}  response.Response
// @Router       /api/sys/v1/role/{id}/menu [get]
func (h *RoleHandler) GetMenuIds(c *gin.Context) {
	idStr := c.Param("id")
	roleId, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.FailBy(c, apperr.ErrInvalidId)
		return
	}
	ids, err := h.svc.GetMenuIds(roleId)
	if err != nil {
		response.FailBy(c, apperr.ErrQueryFailed)
		return
	}
	response.Success(c, ids)
}
