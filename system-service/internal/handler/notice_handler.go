package handler

/**
 * 通知公告 Handler
 *
 * @author
 * @date 2026-04-21
 * @website https://www.aevons.com
 * @email 1073602@qq.com
 * @copyright 2025-2026 Aevons
 */

import (
	"system-service/internal/dto"
	"system-service/internal/model"
	"system-service/internal/service"

	"github.com/calmlax/aevons-framework/core/base"
	"github.com/gin-gonic/gin"
)

type NoticeHandler struct {
	crud *base.BaseHandler[model.Notice, *dto.NoticeQuery, dto.CreateNoticeDTO, dto.UpdateNoticeDTO]
	svc  service.NoticeService
}

// 构造函数
func NewNoticeHandler(svc service.NoticeService) *NoticeHandler {
	return &NoticeHandler{
		crud: base.NewBaseHandler[model.Notice, *dto.NoticeQuery, dto.CreateNoticeDTO, dto.UpdateNoticeDTO](svc),
		svc:  svc,
	}
}

// List 查询通知公告列表。
//
// @Summary      查询通知公告列表
// @Tags         通知公告
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  response.Response
// @Router       /notice/list [get]
func (h *NoticeHandler) List(c *gin.Context) {
	h.crud.HandleList(c)
}

// Page 分页查询通知公告。
//
// @Summary      分页查询通知公告
// @Tags         通知公告
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  response.Response
// @Router       /notice/page [get]
func (h *NoticeHandler) Page(c *gin.Context) {
	h.crud.HandlePage(c)
}

// Get 获取通知公告详情。
//
// @Summary      获取通知公告详情
// @Tags         通知公告
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "通知公告ID"
// @Success      200  {object}  response.Response
// @Router       /notice/{id} [get]
func (h *NoticeHandler) Get(c *gin.Context) {
	h.crud.HandleGet(c)
}

// Create 新增通知公告。
//
// @Summary      新增通知公告
// @Tags         通知公告
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      dto.CreateNoticeDTO  true  "新增通知公告"
// @Success      200      {object}  response.Response
// @Router       /notice [post]
func (h *NoticeHandler) Create(c *gin.Context) {
	h.crud.HandleCreate(c)
}

// Update 修改通知公告。
//
// @Summary      修改通知公告
// @Tags         通知公告
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      int                 true  "通知公告ID"
// @Param        request  body      dto.UpdateNoticeDTO true  "修改通知公告"
// @Success      200      {object}  response.Response
// @Router       /notice/{id} [put]
func (h *NoticeHandler) Update(c *gin.Context) {
	h.crud.HandleUpdate(c)
}

// BatchDelete 批量删除通知公告。
//
// @Summary      批量删除通知公告
// @Tags         通知公告
// @Produce      json
// @Security     BearerAuth
// @Param        ids   path      string  true  "逗号分隔的通知公告ID"
// @Success      200   {object}  response.Response
// @Router       /notice/{ids} [delete]
func (h *NoticeHandler) BatchDelete(c *gin.Context) {
	h.crud.HandleBatchDelete(c)
}
