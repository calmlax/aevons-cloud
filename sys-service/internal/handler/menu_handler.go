/**
 * 菜单权限表 Handler
 *
 * @author
 * @date 2026-04-09 01:17:32.070579992 +0000 UTC
 * @website https://www.aevons.com
 * @email 1073602@qq.com
 * @copyright 2025-2026 Aevons
 */
package handler

import (
	"errors"
	"net/http"
	"sys-service/internal/dto"
	"sys-service/internal/service"

	apperr "github.com/calmlax/aevons-framework/errors"

	"github.com/calmlax/aevons-framework/consts"
	"github.com/calmlax/aevons-framework/core/response"

	"github.com/calmlax/aevons-framework/core/base"

	"github.com/gin-gonic/gin"
)

type MenuHandler struct {
	svc service.MenuService
}

// 构造函数
func NewMenuHandler(svc service.MenuService) *MenuHandler {
	return &MenuHandler{
		svc: svc,
	}
}

// CreateMenu 添加菜单
//
// @Summary      新增菜单
// @Tags         菜单管理
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      dto.CreateMenuDTO  true  "新增菜单"
// @Success      200      {object}  response.Response
// @Router       /menu [post]
func (h *MenuHandler) CreateMenu(c *gin.Context) {
	var Menu dto.CreateMenuDTO
	if err := c.ShouldBindJSON(&Menu); err != nil {
		response.FailBy(c, apperr.ErrInvalidBody)
		return
	}
	MenuId, err := h.svc.CreateMenu(c, Menu)
	if err != nil {
		response.FailBy(c, apperr.ErrCreateFailed, map[string]any{"error": err.Error()})
		return
	}
	response.Success(c, MenuId)
}

// UpdateMenu 修改菜单。
//
// @Summary      修改菜单
// @Tags         菜单管理
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      int               true  "菜单ID"
// @Param        request  body      dto.UpdateMenuDTO true  "修改菜单"
// @Success      200      {object}  response.Response
// @Router       /menu/{id} [put]
func (h *MenuHandler) UpdateMenu(c *gin.Context) {
	var Menu dto.UpdateMenuDTO
	if err := c.ShouldBindJSON(&Menu); err != nil {
		response.FailBy(c, apperr.ErrInvalidBody)
		return
	}
	if err := h.svc.UpdateMenu(c, Menu); err != nil {
		response.FailBy(c, apperr.ErrUpdateFailed, map[string]any{"error": err.Error()})
		return
	}
	response.Success(c, nil)
}

// ListByLangCode 菜单列表
//
// @Summary      按语言查询菜单列表
// @Tags         菜单管理
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  response.Response
// @Router       /menu/list [get]
func (h *MenuHandler) ListByLangCode(c *gin.Context) {
	langCode := c.GetHeader(consts.AcceptLanguage)

	list, err := h.svc.ListByLangCode(langCode)

	if err != nil {
		response.Success(c, []dto.MenuDTO{})
		return
	}
	response.Success(c, list)
}

// GetDetail 获取详情
//
// @Summary      获取菜单详情
// @Tags         菜单管理
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "菜单ID"
// @Success      200  {object}  response.Response
// @Router       /menu/{id} [get]
func (h *MenuHandler) GetDetail(c *gin.Context) {
	id, ok := base.GetId(c)
	if !ok {
		return
	}
	dictData, err := h.svc.GetDetail(id)
	if err != nil {
		response.Success(c, nil)
		return
	}
	response.Success(c, dictData)
}

// Delete 删除
//
// @Summary      批量删除菜单
// @Tags         菜单管理
// @Produce      json
// @Security     BearerAuth
// @Param        ids   path      string  true  "逗号分隔的菜单ID"
// @Success      200   {object}  response.Response
// @Router       /menu/{ids} [delete]
func (h *MenuHandler) BatchDelete(c *gin.Context) {
	ids, ok := base.GetIds(c)
	if !ok {
		return
	}
	err := h.svc.DeleteByIds(c, ids)

	if err != nil {
		if errors.Is(err, apperr.ErrorExisting) {
			response.Fail(c, http.StatusInternalServerError, 5000, "err.sys.menu.exist_subordinate")
			return
		}
		response.FailBy(c, apperr.ErrDeleteFailed)
		return
	}
	response.Success(c, ids)
}
