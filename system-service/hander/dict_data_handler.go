/**
 * 字典数据表 Handler
 *
 * @author
 * @date 2026-04-09 01:08:50.443674979 +0000 UTC
 * @website https://www.aevons.com
 * @email 1073602@qq.com
 * @copyright 2025-2026 Aevons
 */
package handler

import (
	"errors"
	"net/http"
	"system-service/dto"
	"system-service/service"

	apperr "github.com/calmlax/aevons-framework/errors"

	"github.com/calmlax/aevons-framework/consts"
	"github.com/calmlax/aevons-framework/response"

	"github.com/calmlax/aevons-framework/core/base"

	"github.com/gin-gonic/gin"
)

type DictDataHandler struct {
	svc service.DictDataService
}

// 构造函数
func NewDictDataHandler(svc service.DictDataService) *DictDataHandler {
	return &DictDataHandler{
		svc: svc,
	}
}

// GetDictDataCache 获取字典数据
//
// @Summary      按字典类型查询字典数据
// @Tags         字典数据
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  response.Response
// @Router       /dict/data/list [get]
func (h *DictDataHandler) ListByDictType(c *gin.Context) {
	dictType := c.Query("dictType")
	langCode := c.GetHeader(consts.AcceptLanguage)
	dictData, err := h.svc.ListByDictType(dictType, langCode)
	if err != nil {
		response.Success(c, nil)
		return
	}
	response.Success(c, dictData)
}

// CreateDictData 添加字典数据
//
// @Summary      新增字典数据
// @Tags         字典数据
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      dto.CreateDictDataDTO  true  "新增字典数据"
// @Success      200      {object}  response.Response
// @Router       /dict/data [post]
func (h *DictDataHandler) CreateDictData(c *gin.Context) {
	var dictData dto.CreateDictDataDTO
	if err := c.ShouldBindJSON(&dictData); err != nil {
		response.FailBy(c, apperr.ErrInvalidBody)
		return
	}
	dictDataId, err := h.svc.CreateDictData(c, dictData)
	if err != nil {
		if errors.Is(err, apperr.ErrorExisting) {
			response.Fail(c, http.StatusInternalServerError, 5000, "err.sys.dict_data.existing", map[string]any{"dictType": dictData.DictType, "dictValue": dictData.DictValue})
			return
		}
		response.FailBy(c, apperr.ErrCreateFailed, map[string]any{"error": err.Error()})
		return
	}
	response.Success(c, dictDataId)
}

// UpdateDictData 修改字典数据
//
// @Summary      修改字典数据
// @Tags         字典数据
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      int                    true  "字典数据ID"
// @Param        request  body      dto.UpdateDictDataDTO  true  "修改字典数据"
// @Success      200      {object}  response.Response
// @Router       /dict/data/{id} [put]
func (h *DictDataHandler) UpdateDictData(c *gin.Context) {
	var dictData dto.UpdateDictDataDTO
	if err := c.ShouldBindJSON(&dictData); err != nil {
		response.FailBy(c, apperr.ErrInvalidBody)
		return
	}
	if err := h.svc.UpdateDictData(c, dictData); err != nil {
		response.FailBy(c, apperr.ErrUpdateFailed, map[string]any{"error": err.Error()})
		return
	}
	response.Success(c, nil)
}

// GetDictDataCache 获取字典数据
//
// @Summary      获取字典缓存
// @Tags         字典数据
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "字典类型"
// @Success      200  {object}  response.Response
// @Router       /dict/type/{id} [get]
func (h *DictDataHandler) GetDictDataCache(c *gin.Context) {
	dictType := c.Param("id")
	dictData, err := h.svc.GetDictDataCache(dictType)
	if err != nil {
		response.Success(c, nil)
		return
	}
	response.Success(c, dictData)
}

// RefreshCache 刷新字典数据
//
// @Summary      刷新字典缓存
// @Tags         字典数据
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  response.Response
// @Router       /dict/refresh-cache [delete]
func (h *DictDataHandler) RefreshCache(c *gin.Context) {
	if err := h.svc.RefreshCache(); err != nil {
		response.Fail(c, http.StatusInternalServerError, 1013, "err.api.failed_to_refresh_cache")
		return
	}
	response.Success(c, nil)
}

// GetDetail 获取字典数据详情
//
// @Summary      获取字典数据详情
// @Tags         字典数据
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "字典数据ID"
// @Success      200  {object}  response.Response
// @Router       /dict/data/{id} [get]
func (h *DictDataHandler) GetDetail(c *gin.Context) {
	id, ok := base.GetId(c)
	if !ok {
		return
	}
	dictData, err := h.svc.GetDetail(id)
	if err != nil {
		response.Success(c, nil)
		return
	}
	response.Success(c, dictData)
}

// Delete 删除
//
// @Summary      批量删除字典数据
// @Tags         字典数据
// @Produce      json
// @Security     BearerAuth
// @Param        ids   path      string  true  "逗号分隔的字典数据ID"
// @Success      200   {object}  response.Response
// @Router       /dict/data/{ids} [delete]
func (h *DictDataHandler) BatchDelete(c *gin.Context) {
	ids, ok := base.GetIds(c)
	if !ok {
		return
	}
	err := h.svc.DeleteByIds(c, ids)

	if err != nil {
		response.FailBy(c, apperr.ErrDeleteFailed)
		return
	}
	response.Success(c, ids)
}

// UpdateSort 批量更新排序
//
// @Summary      批量更新字典数据排序
// @Tags         字典数据
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      []dto.SortItemDTO  true  "排序项"
// @Success      200      {object}  response.Response
// @Router       /dict/data/sort [put]
func (h *DictDataHandler) UpdateSort(c *gin.Context) {
	var items []dto.SortItemDTO
	if err := c.ShouldBindJSON(&items); err != nil {
		response.FailBy(c, apperr.ErrInvalidBody)
		return
	}
	if err := h.svc.UpdateSort(c, items); err != nil {
		response.FailBy(c, apperr.ErrUpdateFailed, map[string]any{"error": err.Error()})
		return
	}
	response.Success(c, nil)
}
