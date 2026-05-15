package handler

import (
	"strconv"
	"system-service/dto"
	"system-service/model"
	"system-service/service"

	apperr "github.com/calmlax/aevons-framework/errors"
	"github.com/calmlax/aevons-framework/response"

	"github.com/calmlax/aevons-framework/core/base"

	"github.com/gin-gonic/gin"
)

type LangHandler struct {
	*base.BaseHandler[
		model.Lang,
		*dto.LangQuery,
		dto.CreateLangDTO,
		dto.UpdateLangDTO,
	]
	svc service.LangService
}

func NewLangHandler(svc service.LangService) *LangHandler {
	return &LangHandler{
		BaseHandler: base.NewBaseHandler[
			model.Lang,
			*dto.LangQuery,
			dto.CreateLangDTO,
			dto.UpdateLangDTO,
		](svc),
		svc: svc,
	}
}

// AvailableList 获取可用语言列表
func (h *LangHandler) AvailableList(c *gin.Context) {
	status := int16(0)
	list, err := h.svc.List(&dto.LangQuery{Status: &status})
	if err != nil {
		response.Success(c, nil)
		return
	}
	response.Success(c, list)
}

// Create 新增语言（若设为默认则清除其他默认）
func (h *LangHandler) Create(c *gin.Context) {
	var d dto.CreateLangDTO
	if err := c.ShouldBindJSON(&d); err != nil {
		response.FailBy(c, apperr.ErrInvalidBody)
		return
	}
	lang, err := h.svc.CreateLang(d)
	if err != nil {
		response.FailBy(c, apperr.ErrCreateFailed, map[string]any{"error": err.Error()})
		return
	}
	response.Success(c, lang)
}

// Update 修改语言（若设为默认则清除其他默认）
func (h *LangHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.FailBy(c, apperr.ErrInvalidId)
		return
	}
	var d dto.UpdateLangDTO
	if err := c.ShouldBindJSON(&d); err != nil {
		response.FailBy(c, apperr.ErrInvalidBody)
		return
	}
	if err := h.svc.UpdateLang(id, d); err != nil {
		response.FailBy(c, apperr.ErrUpdateFailed, map[string]any{"error": err.Error()})
		return
	}
	response.Success(c, id)
}
