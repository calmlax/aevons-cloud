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
	"system-service/dto"
	"system-service/model"
	"system-service/service"

	apperr "github.com/calmlax/aevons-framework/errors"

	"github.com/calmlax/aevons-framework/consts"
	"github.com/calmlax/aevons-framework/response"

	"github.com/calmlax/aevons-framework/core/base"

	"github.com/gin-gonic/gin"
)

type MenuHandler struct {
	// 继承BaseHandler
	*base.BaseHandler[
		model.Menu,        // 模型
		*dto.MenuQuery,    // 查询 DTO
		dto.CreateMenuDTO, // 创建 DTO
		dto.UpdateMenuDTO, // 更新 DTO
	]
	svc service.MenuService
}

// 构造函数
func NewMenuHandler(svc service.MenuService) *MenuHandler {
	return &MenuHandler{
		BaseHandler: base.NewBaseHandler[
			model.Menu,
			*dto.MenuQuery,
			dto.CreateMenuDTO,
			dto.UpdateMenuDTO,
		](svc),
		svc: svc,
	}
}

// CreateMenu 添加菜单
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

// CreateMenu 添加菜单
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
func (h *MenuHandler) GetDetail(c *gin.Context) {
	id, ok := h.GetId(c)
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
func (h *MenuHandler) BatchDelete(c *gin.Context) {
	ids, ok := h.GetIds(c)
	if !ok {
		return
	}
	err := h.svc.DeleteByIds(c, ids)

	if err != nil {
		if errors.Is(err, apperr.ErrorExisting) {
			response.Fail(c, http.StatusInternalServerError, 5000, "err.sys.menu.exist_subordinate")
			return
		}
		h.Fail(c, apperr.ErrDeleteFailed)
		return
	}
	h.Success(c, ids)
}
