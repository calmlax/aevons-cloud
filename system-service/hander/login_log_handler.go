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
	"net/http"
	"system-service/dto"
	"system-service/model"
	"system-service/service"

	pkgauth "github.com/calmlax/aevons-framework/auth"
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

func (h *LoginLogHandler) GetProfileLoginLog(c *gin.Context) {
	user, err := pkgauth.GetCurrentUser(c.Request.Context())
	if err != nil {
		response.Fail(c, http.StatusUnauthorized, 401, "err.sys.unauthorized")
		return
	}
	result, err := h.svc.GetProfileLoginLogs(user.Username)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, 1001, "err.api.failed_to_query_login_logs")
		return
	}
	response.Success(c, result)
}

// Clear 清空所有登录日志
func (h *LoginLogHandler) Clear(c *gin.Context) {
	if err := h.svc.TruncateAll(); err != nil {
		response.FailBy(c, apperr.ErrDeleteFailed)
		return
	}
	response.Success(c, nil)
}
