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
	crud *base.BaseHandler[model.OauthClient, *dto.OauthClientQuery, dto.CreateOauthClientDTO, dto.UpdateOauthClientDTO]
	svc service.OauthClientService
}

// 构造函数
func NewOauthClientHandler(svc service.OauthClientService) *OauthClientHandler {
	return &OauthClientHandler{
		crud: base.NewBaseHandler[model.OauthClient, *dto.OauthClientQuery, dto.CreateOauthClientDTO, dto.UpdateOauthClientDTO](svc),
		svc: svc,
	}
}

// List 查询终端应用列表。
//
// @Summary      查询终端应用列表
// @Tags         终端应用
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  response.Response
// @Router       /oauth/client/list [get]
func (h *OauthClientHandler) List(c *gin.Context) {
	h.crud.HandleList(c)
}

// Page 分页查询终端应用。
//
// @Summary      分页查询终端应用
// @Tags         终端应用
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  response.Response
// @Router       /oauth/client/page [get]
func (h *OauthClientHandler) Page(c *gin.Context) {
	h.crud.HandlePage(c)
}

// Get 获取终端应用详情。
//
// @Summary      获取终端应用详情
// @Tags         终端应用
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "终端应用ID"
// @Success      200  {object}  response.Response
// @Router       /oauth/client/{id} [get]
func (h *OauthClientHandler) Get(c *gin.Context) {
	h.crud.HandleGet(c)
}

// BatchDelete 批量删除终端应用。
//
// @Summary      批量删除终端应用
// @Tags         终端应用
// @Produce      json
// @Security     BearerAuth
// @Param        ids   path      string  true  "逗号分隔的终端应用ID"
// @Success      200   {object}  response.Response
// @Router       /oauth/client/{ids} [delete]
func (h *OauthClientHandler) BatchDelete(c *gin.Context) {
	h.crud.HandleBatchDelete(c)
}

// CreateOAuthClient 添加终端应用
//
// @Summary      新增终端应用
// @Tags         终端应用
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      dto.CreateOauthClientDTO  true  "新增终端应用"
// @Success      200      {object}  response.Response
// @Router       /oauth/client [post]
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
//
// @Summary      修改终端应用
// @Tags         终端应用
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      int                       true  "终端应用ID"
// @Param        request  body      dto.UpdateOauthClientDTO  true  "修改终端应用"
// @Success      200      {object}  response.Response
// @Router       /oauth/client/{id} [put]
func (h *OauthClientHandler) UpdateOAuthClient(c *gin.Context) {
	id, ok := base.GetId(c)
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
		response.FailBy(c, apperr.ErrUpdateFailed)
		return
	}
	response.Success(c, nil)
}
