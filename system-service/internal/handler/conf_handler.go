/**
 * 参数配置表 Handler
 *
 * @author
 * @date 2026-04-09 00:38:25.504785055 +0000 UTC
 * @website https://www.aevons.com
 * @email 1073602@qq.com
 * @copyright 2025-2026 Aevons
 */
package handler

import (
	"errors"
	"fmt"
	"net/http"
	"system-service/internal/dto"
	"system-service/internal/model"
	"system-service/internal/service"

	"github.com/calmlax/aevons-framework/response"
	"github.com/calmlax/aevons-framework/utils"

	"github.com/calmlax/aevons-framework/core/base"
	apperr "github.com/calmlax/aevons-framework/errors"

	"github.com/gin-gonic/gin"
)

type ConfHandler struct {
	crud *base.BaseHandler[model.Conf, *dto.ConfQuery, dto.CreateConfDTO, dto.UpdateConfDTO]
	svc  service.ConfService
}

// 构造函数
func NewConfHandler(svc service.ConfService) *ConfHandler {
	return &ConfHandler{
		crud: base.NewBaseHandler[model.Conf, *dto.ConfQuery, dto.CreateConfDTO, dto.UpdateConfDTO](svc),
		svc:  svc,
	}
}

// List 查询参数配置列表。
//
// @Summary      查询参数配置列表
// @Tags         参数配置
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  response.Response
// @Router       /conf/list [get]
func (h *ConfHandler) List(c *gin.Context) {
	h.crud.HandleList(c)
}

// Page 分页查询参数配置。
//
// @Summary      分页查询参数配置
// @Tags         参数配置
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  response.Response
// @Router       /conf/page [get]
func (h *ConfHandler) Page(c *gin.Context) {
	h.crud.HandlePage(c)
}

// GetConfByKey 提供从 URL 动态参数提取由服务层按 key 读取（自动解密）的配置获取接口。
//
// @Summary      按配置键获取参数配置
// @Tags         参数配置
// @Produce      json
// @Param        key   path      string  true  "配置键"
// @Success      200   {object}  response.Response
// @Failure      400   {object}  response.Response
// @Router       /conf/key/{key} [get]
func (h *ConfHandler) GetConfByKey(c *gin.Context) {
	key := c.Param("key")
	if key == "" {
		response.Fail(c, http.StatusBadRequest, 1002, "err.api.invalid_conf_key")
		return
	}
	conf, err := h.svc.GetConfByKey(key)
	if err != nil {
		if errors.Is(err, apperr.ErrorNotFound) {
			response.FailByErr(c, http.StatusNotFound, apperr.ErrDataNotFound)
			return
		}
		response.Fail(c, http.StatusInternalServerError, 1003, "err.api.failed_to_fetch_conf")
		return
	}
	response.Success(c, conf)
}

// RefreshCache 后台管理开放的清理配置缓存接口（全量清除关联键前缀缓存）。
//
// @Summary      刷新参数配置缓存
// @Tags         参数配置
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  response.Response
// @Router       /conf/refresh-cache [post]
func (h *ConfHandler) RefreshCache(c *gin.Context) {
	if err := h.svc.RefreshCache(); err != nil {
		response.Fail(c, http.StatusInternalServerError, 1013, "err.api.failed_to_refresh_cache")
		return
	}
	response.Success(c, nil)
}

// Get 获取配置
//
// @Summary      获取参数配置详情
// @Tags         参数配置
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "配置ID"
// @Success      200  {object}  response.Response
// @Router       /conf/{id} [get]
func (h *ConfHandler) Get(c *gin.Context) {
	id, ok := base.GetId(c)
	if !ok {
		return
	}

	conf, err := h.svc.GetById(id)
	if err != nil {
		response.FailBy(c, apperr.ErrQueryFailed)
		return
	}
	conf.ConfValue = h.svc.DecryptIfNeeded(conf.ConfValue, conf.IsSecret == 1)
	response.Success(c, conf)
}

// Delete 删除参数配置。
//
// @Summary      删除参数配置
// @Tags         参数配置
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "配置ID"
// @Success      200  {object}  response.Response
// @Router       /conf/{id} [delete]
func (h *ConfHandler) Delete(c *gin.Context) {
	h.crud.HandleDelete(c)
}

// CreateConf 添加配置
//
// @Summary      新增参数配置
// @Tags         参数配置
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      dto.CreateConfDTO  true  "新增参数配置"
// @Success      200      {object}  response.Response
// @Router       /conf [post]
func (h *ConfHandler) CreateConf(c *gin.Context) {
	var confDto dto.CreateConfDTO
	if err := c.ShouldBindJSON(&confDto); err != nil {
		response.FailBy(c, apperr.ErrInvalidBody)
		return
	}
	isExist, err := h.svc.ExistField("conf_key", confDto.ConfKey)
	if err != nil || isExist {
		response.Fail(c, http.StatusInternalServerError, 5000, "err.sys.conf.existing", map[string]any{"confKey": confDto.ConfKey})
		return
	}
	confDto.ConfValue = h.svc.EncryptIfNeeded(confDto.ConfValue, confDto.IsSecret == 1)
	var conf model.Conf
	utils.Copy(&conf, confDto)

	if err := h.svc.Create(&conf); err != nil {
		response.FailBy(c, apperr.ErrCreateFailed, map[string]any{"error": err.Error()})
		return
	}
	response.Success(c, conf.Id)
}

// UpdateConf 修改配置
//
// @Summary      修改参数配置
// @Tags         参数配置
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      int                true  "配置ID"
// @Param        request  body      dto.UpdateConfDTO  true  "修改参数配置"
// @Success      200      {object}  response.Response
// @Router       /conf/{id} [put]
func (h *ConfHandler) UpdateConf(c *gin.Context) {
	id, ok := base.GetId(c)
	if !ok {
		return
	}
	var confDto dto.UpdateConfDTO
	if err := c.ShouldBindJSON(&confDto); err != nil {
		response.FailBy(c, apperr.ErrInvalidBody)
		return
	}
	isExist, err := h.svc.ExistFieldExcludeId("conf_key", confDto.ConfKey, id)
	if err != nil || isExist {
		response.Fail(c, http.StatusInternalServerError, 5000, "err.sys.conf.existing", map[string]any{"confKey": confDto.ConfKey})
		return
	}
	fmt.Printf("---------UpdateConf--------IsSecret %v", *confDto.IsSecret == 1)
	value := h.svc.EncryptIfNeeded(*confDto.ConfValue, *confDto.IsSecret == 1)
	confDto.ConfValue = &value
	fmt.Printf("---------UpdateConf--------value %v", value)

	mp := utils.StructToMapIgnoreNil(confDto)
	_, err2 := h.svc.Update(id, mp)
	if err2 != nil {
		response.FailBy(c, apperr.ErrUpdateFailed)
		return
	}
	response.Success(c, nil)
}
