/**
 * 操作日志记录 Handler
 *
 * @author
 * @date 2026-04-14
 * @website https://www.aevons.com
 * @email 1073602@qq.com
 * @copyright 2025-2026 Aevons
 */
package handler

import (
	"log-service/dto"
	"log-service/model"
	"log-service/service"

	apperr "github.com/calmlax/aevons-framework/errors"

	"github.com/calmlax/aevons-framework/response"

	"github.com/calmlax/aevons-framework/core/base"

	"github.com/gin-gonic/gin"
)

type LoginLogHandler struct {
	*base.BaseHandler[model.LoginLog, *dto.LoginLogQuery, dto.CreateLoginLogDTO, dto.UpdateLoginLogDTO]
	svc service.LoginLogService
}

func NewLoginLogHandler(svc service.LoginLogService) *LoginLogHandler {
	return &LoginLogHandler{
		BaseHandler: base.NewBaseHandler[model.LoginLog, *dto.LoginLogQuery, dto.CreateLoginLogDTO, dto.UpdateLoginLogDTO](svc),
		svc:         svc,
	}
}

// Clear 清空所有登录日志
func (h *LoginLogHandler) Clear(c *gin.Context) {
	if err := h.svc.TruncateAll(); err != nil {
		response.FailBy(c, apperr.ErrDeleteFailed)
		return
	}
	response.Success(c, nil)
}
