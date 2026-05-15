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
	"system-service/dto"
	"system-service/model"
	"system-service/service"

	apperr "github.com/calmlax/aevons-framework/errors"
	"github.com/calmlax/aevons-framework/response"

	"github.com/calmlax/aevons-framework/core/base"

	"github.com/gin-gonic/gin"
)

type OperLogHandler struct {
	*base.BaseHandler[model.OperLog, *dto.OperLogQuery, dto.CreateOperLogDTO, dto.UpdateOperLogDTO]
	svc service.OperLogService
}

func NewOperLogHandler(svc service.OperLogService) *OperLogHandler {
	return &OperLogHandler{
		BaseHandler: base.NewBaseHandler[model.OperLog, *dto.OperLogQuery, dto.CreateOperLogDTO, dto.UpdateOperLogDTO](svc),
		svc:         svc,
	}
}

// Clear 清空所有操作日志
func (h *OperLogHandler) Clear(c *gin.Context) {
	if err := h.svc.TruncateAll(); err != nil {
		response.FailBy(c, apperr.ErrDeleteFailed)
		return
	}
	response.Success(c, nil)
}
