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
	"system-service/dto"
	"system-service/model"
	"system-service/service"

	"github.com/calmlax/aevons-framework/core/base"
)

type NoticeHandler struct {
	// 继承BaseHandler
	*base.BaseHandler[
		model.Notice,        // 模型
		*dto.NoticeQuery,    // 查询 DTO
		dto.CreateNoticeDTO, // 创建 DTO
		dto.UpdateNoticeDTO, // 更新 DTO
	]
	svc service.NoticeService
}

// 构造函数
func NewNoticeHandler(svc service.NoticeService) *NoticeHandler {
	return &NoticeHandler{
		BaseHandler: base.NewBaseHandler[
			model.Notice,
			*dto.NoticeQuery,
			dto.CreateNoticeDTO,
			dto.UpdateNoticeDTO,
		](svc),
		svc: svc,
	}
}
