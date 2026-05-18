/**
 * 字典类型表 Handler
 *
 * @author
 * @date 2026-04-09 01:08:50.442643548 +0000 UTC
 * @website https://www.aevons.com
 * @email 1073602@qq.com
 * @copyright 2025-2026 Aevons
 */
package handler

import (
	"fmt"
	"net/http"
	"sys-service/internal/dto"
	"sys-service/internal/model"
	"sys-service/internal/service"

	"github.com/calmlax/aevons-framework/core/response"
	apperr "github.com/calmlax/aevons-framework/errors"
	"github.com/calmlax/aevons-framework/utils"

	"github.com/calmlax/aevons-framework/core/base"

	"github.com/gin-gonic/gin"
)

type DictHandler struct {
	crud  *base.BaseHandler[model.Dict, *dto.DictQuery, dto.CreateDictDTO, dto.UpdateDictDTO]
	svc   service.DictService
	ddsve service.DictDataService
}

// 构造函数
func NewDictHandler(svc service.DictService, ddsve service.DictDataService) *DictHandler {
	return &DictHandler{
		crud:  base.NewBaseHandler[model.Dict, *dto.DictQuery, dto.CreateDictDTO, dto.UpdateDictDTO](svc),
		svc:   svc,
		ddsve: ddsve,
	}
}

// AvailableList 获取可用列表
//
// @Summary      查询可用字典类型列表
// @Tags         字典管理
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  response.Response
// @Router       /dict/list [get]
func (h *DictHandler) AvailableList(c *gin.Context) {
	status := int16(0)
	list, err := h.svc.List(&dto.DictQuery{
		Status: &status,
	})
	if err != nil {
		response.Success(c, nil)
		return
	}
	response.Success(c, list)
}

// Page 分页查询字典类型。
//
// @Summary      分页查询字典类型
// @Tags         字典管理
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  response.Response
// @Router       /dict/page [get]
func (h *DictHandler) Page(c *gin.Context) {
	h.crud.HandlePage(c)
}

// Get 获取字典类型详情。
//
// @Summary      获取字典类型详情
// @Tags         字典管理
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "字典类型ID"
// @Success      200  {object}  response.Response
// @Router       /dict/{id} [get]
func (h *DictHandler) Get(c *gin.Context) {
	h.crud.HandleGet(c)
}

// CreateDict 新增字典类型。
//
// @Summary      新增字典类型
// @Tags         字典管理
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      dto.CreateDictDTO  true  "新增字典类型"
// @Success      200      {object}  response.Response
// @Router       /dict [post]
func (h *DictHandler) CreateDict(c *gin.Context) {
	var dictDto dto.CreateDictDTO
	if !base.BindJSON(c, &dictDto) {
		return
	}
	fmt.Printf("-----CreateDict----- DictType: %v \n", dictDto.DictType)
	isExist, _ := h.svc.ExistField("dict_type", dictDto.DictType)
	if isExist {
		response.Fail(c, http.StatusInternalServerError, 5000, "err.sys.dict.existing", map[string]any{"dictType": dictDto.DictType})
		return
	}

	var m model.Dict
	utils.Copy(&m, dictDto)

	err := h.svc.Create(&m)
	if err != nil {
		response.FailBy(c, apperr.ErrCreateFailed, map[string]any{"error": err.Error()})
		return
	}
	response.Success(c, m.Id)
}

// UpdateDict 修改字典类型。
//
// @Summary      修改字典类型
// @Tags         字典管理
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      int               true  "字典类型ID"
// @Param        request  body      dto.UpdateDictDTO true  "修改字典类型"
// @Success      200      {object}  response.Response
// @Router       /dict/{id} [put]
func (h *DictHandler) UpdateDict(c *gin.Context) {
	id, ok := base.GetId(c)
	if !ok {
		return
	}
	var dictDto dto.UpdateDictDTO
	if err := c.ShouldBindJSON(&dictDto); err != nil {
		response.FailBy(c, apperr.ErrInvalidBody)
		return
	}
	mp := utils.StructToMapIgnoreNil(dictDto)
	err := h.svc.UpdateDict(c, id, mp)
	if err != nil {
		response.FailBy(c, apperr.ErrUpdateFailed)
		return
	}
	response.Success(c, nil)
}

// DeleteDict 删除字典类型。
//
// @Summary      删除字典类型
// @Tags         字典管理
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "字典类型ID"
// @Success      200  {object}  response.Response
// @Router       /dict/{id} [delete]
func (h *DictHandler) DeleteDict(c *gin.Context) {
	id, ok := base.GetId(c)
	if !ok {
		return
	}
	dict, err := h.svc.GetById(id)
	if err != nil || utils.IsEmpty(dict) {
		response.FailBy(c, apperr.ErrDeleteFailed)
		return
	}
	isExist, _ := h.ddsve.ExistField("dict_type", dict.DictType)
	if isExist {
		response.Fail(c, http.StatusInternalServerError, 5000, "err.sys.dict.dict_data_existing", map[string]any{"dictType": dict.DictType})
		return
	}
	err2 := h.svc.Delete(id)

	if err2 != nil {
		response.FailBy(c, apperr.ErrDeleteFailed)
		return
	}
	response.Success(c, id)
}
