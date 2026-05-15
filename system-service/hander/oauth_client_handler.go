/**
 * 终端应用 Handler
 *
 * @author
 * @date 2026-04-09 01:26:40.390618977 +0000 UTC
 * @website https://www.aevons.com
 * @email 1073602@qq.com
 * @copyright 2025-2026 Aevons
 */
package handler

import (
	"net/http"
	"system-service/dto"
	"system-service/model"
	"system-service/service"

	apperr "github.com/calmlax/aevons-framework/errors"
	"github.com/calmlax/aevons-framework/response"
	"github.com/calmlax/aevons-framework/utils"

	"github.com/calmlax/aevons-framework/core/base"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type OauthClientHandler struct {
	// 继承BaseHandler
	*base.BaseHandler[
		model.OauthClient,        // 模型
		*dto.OauthClientQuery,    // 查询 DTO
		dto.CreateOauthClientDTO, // 创建 DTO
		dto.UpdateOauthClientDTO, // 更新 DTO
	]
	svc service.OauthClientService
}

// 构造函数
func NewOauthClientHandler(svc service.OauthClientService) *OauthClientHandler {
	return &OauthClientHandler{
		BaseHandler: base.NewBaseHandler[
			model.OauthClient,
			*dto.OauthClientQuery,
			dto.CreateOauthClientDTO,
			dto.UpdateOauthClientDTO,
		](svc),
		svc: svc,
	}
}

// CreateOAuthClient 添加终端应用
func (h *OauthClientHandler) CreateOAuthClient(c *gin.Context) {
	var oauthClientDto dto.CreateOauthClientDTO
	if err := c.ShouldBindJSON(&oauthClientDto); err != nil {
		response.FailBy(c, apperr.ErrInvalidBody)
		return
	}
	isExist, err := h.svc.ExistField("client_id", oauthClientDto.ClientId)
	if err != nil || isExist {
		response.Fail(c, http.StatusInternalServerError, 5000, "err.sys.oauth_client.existing", map[string]any{"clientId": oauthClientDto.ClientId})
		return
	}

	var oc model.OauthClient
	utils.Copy(&oc, oauthClientDto)
	hashed, err := bcrypt.GenerateFromPassword([]byte(oauthClientDto.ClientSecret), bcrypt.DefaultCost)
	if err != nil {
		response.FailBy(c, apperr.ErrCreateFailed)
	}
	oc.ClientSecret = string(hashed)

	if err := h.svc.Create(&oc); err != nil {
		response.FailBy(c, apperr.ErrCreateFailed, map[string]any{"error": err.Error()})
		return
	}
	response.Success(c, oc.Id)
}

// UpdateOAuthClient 修改终端应用
func (h *OauthClientHandler) UpdateOAuthClient(c *gin.Context) {
	id, ok := h.GetId(c)
	if !ok {
		return
	}
	var oauthClientDto dto.UpdateOauthClientDTO
	if err := c.ShouldBindJSON(&oauthClientDto); err != nil {
		response.FailBy(c, apperr.ErrInvalidBody)
		return
	}
	isExist, err := h.svc.ExistFieldExcludeId("client_id", oauthClientDto.ClientId, id)
	if err != nil || isExist {
		response.Fail(c, http.StatusInternalServerError, 5000, "err.sys.oauth_client.existing", map[string]any{"clientId": oauthClientDto.ClientId})
		return
	}
	if utils.IsNotEmpty(oauthClientDto.ClientSecret) {
		clientSecret := *oauthClientDto.ClientSecret
		hashed, err := bcrypt.GenerateFromPassword([]byte(clientSecret), bcrypt.DefaultCost)
		if err != nil {
			response.FailBy(c, apperr.ErrUpdateFailed)
			return
		}
		hs := string(hashed)
		oauthClientDto.ClientSecret = &hs
	}

	mp := utils.StructToMapIgnoreNil(oauthClientDto)
	_, err2 := h.svc.Update(id, mp)
	if err2 != nil {
		h.Fail(c, apperr.ErrUpdateFailed)
		return
	}
	response.Success(c, nil)
}
