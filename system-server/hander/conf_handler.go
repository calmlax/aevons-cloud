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
	"system-server/dto"
	"system-server/model"
	"system-server/service"

	"github.com/calmlax/aevons-framework/response"
	"github.com/calmlax/aevons-framework/utils"

	"github.com/calmlax/aevons-framework/core/base"
	errdef "github.com/calmlax/aevons-framework/err"

	"github.com/gin-gonic/gin"
)

type ConfHandler struct {
	// 继承BaseHandler
	*base.BaseHandler[
		model.Conf,        // 模型
		*dto.ConfQuery,    // 查询 DTO
		dto.CreateConfDTO, // 创建 DTO
		dto.UpdateConfDTO, // 更新 DTO
	]
	svc service.ConfService
}

// 构造函数
func NewConfHandler(svc service.ConfService) *ConfHandler {
	return &ConfHandler{
		BaseHandler: base.NewBaseHandler[
			model.Conf,
			*dto.ConfQuery,
			dto.CreateConfDTO,
			dto.UpdateConfDTO,
		](svc),
		svc: svc,
	}
}

// GetConfByKey 提供从 URL 动态参数提取由服务层按 key 读取（自动解密）的配置获取接口。
func (h *ConfHandler) GetConfByKey(c *gin.Context) {
	key := c.Param("key")
	if key == "" {
		response.Fail(c, http.StatusBadRequest, 1002, "err.api.invalid_conf_key")
		return
	}
	conf, err := h.svc.GetConfByKey(key)
	if err != nil {
		if errors.Is(err, errdef.ErrorNotFound) {
			response.FailByErr(c, http.StatusNotFound, errdef.ErrDataNotFound)
			return
		}
		response.Fail(c, http.StatusInternalServerError, 1003, "err.api.failed_to_fetch_conf")
		return
	}
	response.Success(c, conf)
}

// RefreshCache 后台管理开放的清理配置缓存接口（全量清除关联键前缀缓存）。
func (h *ConfHandler) RefreshCache(c *gin.Context) {
	if err := h.svc.RefreshCache(); err != nil {
		response.Fail(c, http.StatusInternalServerError, 1013, "err.api.failed_to_refresh_cache")
		return
	}
	response.Success(c, nil)
}

// Get 获取配置
func (h *ConfHandler) Get(c *gin.Context) {
	id, ok := h.GetId(c)
	if !ok {
		return
	}

	conf, err := h.svc.GetById(id)
	if err != nil {
		h.Fail(c, errdef.ErrQueryFailed)
		return
	}
	conf.ConfValue = h.svc.DecryptIfNeeded(conf.ConfValue, conf.IsSecret == 1)
	response.Success(c, conf)
}

// CreateConf 添加配置
func (h *ConfHandler) CreateConf(c *gin.Context) {
	var confDto dto.CreateConfDTO
	if err := c.ShouldBindJSON(&confDto); err != nil {
		response.FailBy(c, errdef.ErrInvalidBody)
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
		response.FailBy(c, errdef.ErrCreateFailed, map[string]any{"error": err.Error()})
		return
	}
	response.Success(c, conf.Id)
}

// UpdateConf 修改配置
func (h *ConfHandler) UpdateConf(c *gin.Context) {
	id, ok := h.GetId(c)
	if !ok {
		return
	}
	var confDto dto.UpdateConfDTO
	if err := c.ShouldBindJSON(&confDto); err != nil {
		response.FailBy(c, errdef.ErrInvalidBody)
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
		h.Fail(c, errdef.ErrUpdateFailed)
		return
	}
	response.Success(c, nil)
}
