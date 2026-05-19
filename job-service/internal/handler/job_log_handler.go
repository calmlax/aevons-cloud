package handler

import (
	"job-service/internal/dto"
	"job-service/internal/model"
	"job-service/internal/service"

	"github.com/calmlax/aevons-framework/core/base"

	"github.com/gin-gonic/gin"
)

type JobLogHandler struct {
	crud *base.BaseHandler[model.JobLog, *dto.JobLogQuery, dto.CreateJobLogDTO, dto.UpdateJobLogDTO]
	svc  service.JobLogService
}

func NewJobLogHandler(svc service.JobLogService) *JobLogHandler {
	return &JobLogHandler{
		crud: base.NewBaseHandler[model.JobLog, *dto.JobLogQuery, dto.CreateJobLogDTO, dto.UpdateJobLogDTO](svc),
		svc:  svc,
	}
}

// Page 分页查询任务执行日志。
//
// @Summary      分页查询任务执行日志
// @Tags         任务日志
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  response.Response
// @Router       /api/job/v1/log/page [get]
func (h *JobLogHandler) Page(c *gin.Context) {
	h.crud.HandlePage(c)
}

// Get 获取任务执行日志详情。
//
// @Summary      获取任务执行日志详情
// @Tags         任务日志
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "任务日志ID"
// @Success      200  {object}  response.Response
// @Router       /api/job/v1/log/{id} [get]
func (h *JobLogHandler) Get(c *gin.Context) {
	h.crud.HandleGet(c)
}

// BatchDelete 批量删除任务执行日志。
//
// @Summary      批量删除任务执行日志
// @Tags         任务日志
// @Produce      json
// @Security     BearerAuth
// @Param        ids   path      string  true  "逗号分隔的任务日志ID"
// @Success      200   {object}  response.Response
// @Router       /api/job/v1/log/{ids} [delete]
func (h *JobLogHandler) BatchDelete(c *gin.Context) {
	h.crud.HandleBatchDelete(c)
}
