package handler

import (
	"strconv"

	"job-service/internal/dto"
	"job-service/internal/model"
	"job-service/internal/service"

	"github.com/calmlax/aevons-framework/core/base"
	apperr "github.com/calmlax/aevons-framework/errors"
	"github.com/calmlax/aevons-framework/response"

	"github.com/gin-gonic/gin"
)

type JobHandler struct {
	crud *base.BaseHandler[model.Job, *dto.JobQuery, dto.CreateJobDTO, dto.UpdateJobDTO]
	svc  service.JobService
}

func NewJobHandler(svc service.JobService) *JobHandler {
	return &JobHandler{
		crud: base.NewBaseHandler[model.Job, *dto.JobQuery, dto.CreateJobDTO, dto.UpdateJobDTO](svc),
		svc:  svc,
	}
}

// Page 分页查询定时任务。
//
// @Summary      分页查询定时任务
// @Tags         定时任务
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  response.Response
// @Router       /job/page [get]
func (h *JobHandler) Page(c *gin.Context) {
	h.crud.HandlePage(c)
}

// Get 获取定时任务详情。
//
// @Summary      获取定时任务详情
// @Tags         定时任务
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "定时任务ID"
// @Success      200  {object}  response.Response
// @Router       /job/{id} [get]
func (h *JobHandler) Get(c *gin.Context) {
	h.crud.HandleGet(c)
}

// Create 新增定时任务。
//
// @Summary      新增定时任务
// @Tags         定时任务
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      dto.CreateJobDTO  true  "新增定时任务"
// @Success      200      {object}  response.Response
// @Router       /job [post]
func (h *JobHandler) Create(c *gin.Context) {
	h.crud.HandleCreate(c)
}

// Update 修改定时任务。
//
// @Summary      修改定时任务
// @Tags         定时任务
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      int               true  "定时任务ID"
// @Param        request  body      dto.UpdateJobDTO  true  "修改定时任务"
// @Success      200      {object}  response.Response
// @Router       /job/{id} [put]
func (h *JobHandler) Update(c *gin.Context) {
	h.crud.HandleUpdate(c)
}

// BatchDelete 批量删除定时任务。
//
// @Summary      批量删除定时任务
// @Tags         定时任务
// @Produce      json
// @Security     BearerAuth
// @Param        ids   path      string  true  "逗号分隔的定时任务ID"
// @Success      200   {object}  response.Response
// @Router       /job/{ids} [delete]
func (h *JobHandler) BatchDelete(c *gin.Context) {
	h.crud.HandleBatchDelete(c)
}

// Trigger 手动触发任务。
//
// @Summary      手动触发任务
// @Tags         定时任务
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "定时任务ID"
// @Success      200  {object}  response.Response
// @Router       /job/{id}/trigger [post]
func (h *JobHandler) Trigger(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.FailBy(c, apperr.ErrInvalidId)
		return
	}
	if err := h.svc.TriggerJob(id); err != nil {
		response.FailBy(c, apperr.ErrUpdateFailed, map[string]any{"error": err.Error()})
		return
	}
	response.Success(c, nil)
}

// ChangeStatus 启动或暂停任务。
//
// @Summary      修改任务状态
// @Tags         定时任务
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      int     true  "定时任务ID"
// @Param        request  body      object  true  "任务状态"
// @Success      200      {object}  response.Response
// @Router       /job/{id}/status [put]
func (h *JobHandler) ChangeStatus(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.FailBy(c, apperr.ErrInvalidId)
		return
	}

	var body struct {
		Status int16 `json:"status" binding:"oneof=0 1"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		response.FailBy(c, apperr.ErrInvalidBody)
		return
	}

	if err := h.svc.ChangeStatus(id, body.Status); err != nil {
		response.FailBy(c, apperr.ErrUpdateFailed, map[string]any{"error": err.Error()})
		return
	}
	response.Success(c, nil)
}
