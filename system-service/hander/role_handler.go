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
	"system-service/dto"
	"system-service/model"
	"system-service/service"

	apperr "github.com/calmlax/aevons-framework/errors"
	"github.com/calmlax/aevons-framework/response"

	"github.com/calmlax/aevons-framework/core/base"

	"github.com/gin-gonic/gin"
)

type RoleHandler struct {
	*base.BaseHandler[
		model.Role,
		*dto.RoleQuery,
		dto.CreateRoleDTO,
		dto.UpdateRoleDTO,
	]
	svc service.RoleService
}

func NewRoleHandler(svc service.RoleService) *RoleHandler {
	return &RoleHandler{
		BaseHandler: base.NewBaseHandler[
			model.Role,
			*dto.RoleQuery,
			dto.CreateRoleDTO,
			dto.UpdateRoleDTO,
		](svc),
		svc: svc,
	}
}

// Create 创建角色（含菜单，事务）
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
