/**
 * 登录日志记录 Handler
 *
 * @author
 * @date 2026-04-14
 * @website https://www.aevons.com
 * @email 1073602@qq.com
 * @copyright 2025-2026 Aevons
 */
package handler

import (
	"log-service/internal/dto"
	"log-service/internal/model"
	"log-service/internal/service"

	apperr "github.com/calmlax/aevons-framework/errors"

	"github.com/calmlax/aevons-framework/core/response"

	"github.com/calmlax/aevons-framework/core/base"

	"github.com/gin-gonic/gin"
)

type LoginLogHandler struct {
	crud *base.BaseHandler[model.LoginLog, *dto.LoginLogQuery, dto.CreateLoginLogDTO, dto.UpdateLoginLogDTO]
	svc  service.LoginLogService
}

func NewLoginLogHandler(svc service.LoginLogService) *LoginLogHandler {
	return &LoginLogHandler{
		crud: base.NewBaseHandler[model.LoginLog, *dto.LoginLogQuery, dto.CreateLoginLogDTO, dto.UpdateLoginLogDTO](svc),
		svc:  svc,
	}
}

// List 查询登录日志列表。
//
// @Summary      查询登录日志列表
// @Tags         登录日志
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  response.Response
// @Router       /login/log/list [get]
func (h *LoginLogHandler) List(c *gin.Context) {
	h.crud.HandleList(c)
}

// Page 分页查询登录日志。
//
// @Summary      分页查询登录日志
// @Tags         登录日志
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  response.Response
// @Router       /login/log/page [get]
func (h *LoginLogHandler) Page(c *gin.Context) {
	h.crud.HandlePage(c)
}

// Get 获取登录日志详情。
//
// @Summary      获取登录日志详情
// @Tags         登录日志
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "登录日志ID"
// @Success      200  {object}  response.Response
// @Router       /login/log/{id} [get]
func (h *LoginLogHandler) Get(c *gin.Context) {
	h.crud.HandleGet(c)
}

// BatchDelete 批量删除登录日志。
//
// @Summary      批量删除登录日志
// @Tags         登录日志
// @Produce      json
// @Security     BearerAuth
// @Param        ids   path      string  true  "逗号分隔的登录日志ID"
// @Success      200   {object}  response.Response
// @Router       /login/log/{ids} [delete]
func (h *LoginLogHandler) BatchDelete(c *gin.Context) {
	h.crud.HandleBatchDelete(c)
}

// Clear 清空所有登录日志
//
// @Summary      清空登录日志
// @Tags         登录日志
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  response.Response
// @Router       /login/log [delete]
func (h *LoginLogHandler) Clear(c *gin.Context) {
	if err := h.svc.TruncateAll(); err != nil {
		response.FailBy(c, apperr.ErrDeleteFailed)
		return
	}
	response.Success(c, nil)
}
