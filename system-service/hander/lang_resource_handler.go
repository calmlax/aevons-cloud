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

type LangResourceHandler struct {
	*base.BaseHandler[
		model.LangResource,
		*dto.LangResourceQuery,
		dto.CreateLangResourceDTO,
		dto.UpdateLangResourceDTO,
	]
	svc service.LangResourceService
}

func NewLangResourceHandler(svc service.LangResourceService) *LangResourceHandler {
	return &LangResourceHandler{
		BaseHandler: base.NewBaseHandler[
			model.LangResource,
			*dto.LangResourceQuery,
			dto.CreateLangResourceDTO,
			dto.UpdateLangResourceDTO,
		](svc),
		svc: svc,
	}
}

// Create 新增语言资源（langCode+namespace+resourceKey 唯一）
func (h *LangResourceHandler) Create(c *gin.Context) {
	var d dto.CreateLangResourceDTO
	if err := c.ShouldBindJSON(&d); err != nil {
		response.FailBy(c, apperr.ErrInvalidBody)
		return
	}
	r, err := h.svc.CreateResource(d)
	if err != nil {
		response.FailBy(c, apperr.ErrCreateFailed, map[string]any{"error": err.Error()})
		return
	}
	response.Success(c, r)
}

// Update 修改语言资源（langCode+namespace+resourceKey 唯一）
func (h *LangResourceHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.FailBy(c, apperr.ErrInvalidId)
		return
	}
	var d dto.UpdateLangResourceDTO
	if err := c.ShouldBindJSON(&d); err != nil {
		response.FailBy(c, apperr.ErrInvalidBody)
		return
	}
	if err := h.svc.UpdateResource(id, d); err != nil {
		response.FailBy(c, apperr.ErrUpdateFailed, map[string]any{"error": err.Error()})
		return
	}
	response.Success(c, id)
}

// GetKeysByNamespace 获取命名空间下的所有 resourceKey
func (h *LangResourceHandler) GetKeysByNamespace(c *gin.Context) {
	ns := c.Query("namespace")
	if ns == "" {
		response.FailBy(c, apperr.ErrInvalidQuery)
		return
	}
	keys, err := h.svc.GetKeysByNamespace(ns)
	if err != nil {
		response.FailBy(c, apperr.ErrQueryFailed)
		return
	}
	response.Success(c, keys)
}

// GetTranslations 获取某个 namespace+resourceKey 的所有语言翻译
func (h *LangResourceHandler) GetTranslations(c *gin.Context) {
	ns := c.Query("namespace")
	key := c.Query("resourceKey")
	if ns == "" || key == "" {
		response.FailBy(c, apperr.ErrInvalidQuery)
		return
	}
	list, err := h.svc.GetTranslations(ns, key)
	if err != nil {
		response.FailBy(c, apperr.ErrQueryFailed)
		return
	}
	response.Success(c, list)
}

// PageKeys 去重分页查询 resourceKey
func (h *LangResourceHandler) PageKeys(c *gin.Context) {
	ns := c.Query("namespace")
	if ns == "" {
		response.FailBy(c, apperr.ErrInvalidQuery)
		return
	}
	resourceKey := c.Query("resourceKey")
	content := c.Query("content")
	pageNum := 1
	pageSize := 20
	if v, err := strconv.Atoi(c.DefaultQuery("pageNum", "1")); err == nil && v > 0 {
		pageNum = v
	}
	if v, err := strconv.Atoi(c.DefaultQuery("pageSize", "20")); err == nil && v > 0 {
		pageSize = v
	}
	offset := (pageNum - 1) * pageSize
	keys, total, err := h.svc.PageKeys(ns, resourceKey, content, offset, pageSize)
	if err != nil {
		response.FailBy(c, apperr.ErrQueryFailed)
		return
	}
	response.Success(c, map[string]any{"rows": keys, "total": total})
}

// SaveTranslations 批量保存翻译（upsert）
func (h *LangResourceHandler) SaveTranslations(c *gin.Context) {
	var d dto.SaveTranslationsDTO
	if err := c.ShouldBindJSON(&d); err != nil {
		response.FailBy(c, apperr.ErrInvalidBody)
		return
	}
	if err := h.svc.SaveTranslations(d.Namespace, d.ResourceKey, d.Items); err != nil {
		response.FailBy(c, apperr.ErrUpdateFailed, map[string]any{"error": err.Error()})
		return
	}
	response.Success(c, nil)
}
