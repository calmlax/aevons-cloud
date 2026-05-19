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
	"log-service/internal/dto"
	"log-service/internal/model"
	"log-service/internal/service"

	apperr "github.com/calmlax/aevons-framework/errors"

	"github.com/calmlax/aevons-framework/response"

	"github.com/calmlax/aevons-framework/core/base"

	"github.com/gin-gonic/gin"
)

type OperLogHandler struct {
	crud *base.BaseHandler[model.OperLog, *dto.OperLogQuery, dto.CreateOperLogDTO, dto.UpdateOperLogDTO]
	svc  service.OperLogService
}

func NewOperLogHandler(svc service.OperLogService) *OperLogHandler {
	return &OperLogHandler{
		crud: base.NewBaseHandler[model.OperLog, *dto.OperLogQuery, dto.CreateOperLogDTO, dto.UpdateOperLogDTO](svc),
		svc:  svc,
	}
}

// List 查询操作日志列表。
//
// @Summary      查询操作日志列表
// @Tags         操作日志
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  response.Response
// @Router       /api/log/v1/oper/log/list [get]
func (h *OperLogHandler) List(c *gin.Context) {
	h.crud.HandleList(c)
}

// Page 分页查询操作日志。
//
// @Summary      分页查询操作日志
// @Tags         操作日志
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  response.Response
// @Router       /api/log/v1/oper/log/page [get]
func (h *OperLogHandler) Page(c *gin.Context) {
	h.crud.HandlePage(c)
}

// Get 获取操作日志详情。
//
// @Summary      获取操作日志详情
// @Tags         操作日志
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "操作日志ID"
// @Success      200  {object}  response.Response
// @Router       /api/log/v1/oper/log/{id} [get]
func (h *OperLogHandler) Get(c *gin.Context) {
	h.crud.HandleGet(c)
}

// BatchDelete 批量删除操作日志。
//
// @Summary      批量删除操作日志
// @Tags         操作日志
// @Produce      json
// @Security     BearerAuth
// @Param        ids   path      string  true  "逗号分隔的操作日志ID"
// @Success      200   {object}  response.Response
// @Router       /api/log/v1/oper/log/{ids} [delete]
func (h *OperLogHandler) BatchDelete(c *gin.Context) {
	h.crud.HandleBatchDelete(c)
}

// Clear 清空所有操作日志
//
// @Summary      清空操作日志
// @Tags         操作日志
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  response.Response
// @Router       /api/log/v1/oper/log [delete]
func (h *OperLogHandler) Clear(c *gin.Context) {
	if err := h.svc.TruncateAll(); err != nil {
		response.FailBy(c, apperr.ErrDeleteFailed)
		return
	}
	response.Success(c, nil)
}
