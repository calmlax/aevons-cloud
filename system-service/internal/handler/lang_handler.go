package handler

import (
	"strconv"
	"system-service/internal/dto"
	"system-service/internal/model"
	"system-service/internal/service"

	apperr "github.com/calmlax/aevons-framework/errors"
	"github.com/calmlax/aevons-framework/response"

	"github.com/calmlax/aevons-framework/core/base"

	"github.com/gin-gonic/gin"
)

type LangHandler struct {
	crud *base.BaseHandler[model.Lang, *dto.LangQuery, dto.CreateLangDTO, dto.UpdateLangDTO]
	svc  service.LangService
}

func NewLangHandler(svc service.LangService) *LangHandler {
	return &LangHandler{
		crud: base.NewBaseHandler[model.Lang, *dto.LangQuery, dto.CreateLangDTO, dto.UpdateLangDTO](svc),
		svc:  svc,
	}
}

// List 查询语言列表。
//
// @Summary      查询语言列表
// @Tags         语言管理
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  response.Response
// @Router       /lang/list [get]
func (h *LangHandler) List(c *gin.Context) {
	h.crud.HandleList(c)
}

// Page 分页查询语言。
//
// @Summary      分页查询语言
// @Tags         语言管理
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  response.Response
// @Router       /lang/page [get]
func (h *LangHandler) Page(c *gin.Context) {
	h.crud.HandlePage(c)
}

// Get 获取语言详情。
//
// @Summary      获取语言详情
// @Tags         语言管理
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "语言ID"
// @Success      200  {object}  response.Response
// @Router       /lang/{id} [get]
func (h *LangHandler) Get(c *gin.Context) {
	h.crud.HandleGet(c)
}

// BatchDelete 批量删除语言。
//
// @Summary      批量删除语言
// @Tags         语言管理
// @Produce      json
// @Security     BearerAuth
// @Param        ids   path      string  true  "逗号分隔的语言ID"
// @Success      200   {object}  response.Response
// @Router       /lang/{ids} [delete]
func (h *LangHandler) BatchDelete(c *gin.Context) {
	h.crud.HandleBatchDelete(c)
}

// AvailableList 获取可用语言列表
//
// @Summary      查询可用语言列表
// @Tags         语言管理
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  response.Response
// @Router       /lang [get]
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
//
// @Summary      新增语言
// @Tags         语言管理
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      dto.CreateLangDTO  true  "新增语言"
// @Success      200      {object}  response.Response
// @Router       /lang [post]
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
//
// @Summary      修改语言
// @Tags         语言管理
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      int               true  "语言ID"
// @Param        request  body      dto.UpdateLangDTO true  "修改语言"
// @Success      200      {object}  response.Response
// @Router       /lang/{id} [put]
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
