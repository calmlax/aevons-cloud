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
	"math/rand"
	"net/http"
	"strings"
	"sys-service/internal/dto"
	"sys-service/internal/model"
	"sys-service/internal/service"
	"time"

	apperr "github.com/calmlax/aevons-framework/errors"
	frameworkredis "github.com/calmlax/aevons-framework/redis"
	"github.com/calmlax/aevons-framework/response"
	"github.com/calmlax/aevons-framework/utils"

	"github.com/calmlax/aevons-framework/core/base"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

const (
	gatewayOauthClientCacheKey       = "gateway-service:oauth-client-rules:v1"
	gatewayOauthClientBackupCacheKey = "gateway-service:oauth-client-rules:backup:v1"
	gatewayOauthClientCacheTTL       = 30 * time.Second
	gatewayOauthClientBackupTTL      = 150 * time.Second
	gatewayOauthClientCacheJitter    = 10 * time.Second
)

type OauthClientHandler struct {
	crud *base.BaseHandler[model.OauthClient, *dto.OauthClientQuery, dto.CreateOauthClientDTO, dto.UpdateOauthClientDTO]
	svc  service.OauthClientService
}

type gatewayOauthClientRule struct {
	ClientID     string              `json:"ClientID"`
	Enabled      bool                `json:"Enabled"`
	AllowAll     bool                `json:"AllowAll"`
	ServiceNames map[string]struct{} `json:"ServiceNames"`
	ExactRules   map[string]struct{} `json:"ExactRules"`
	PrefixRules  []string            `json:"PrefixRules"`
}

// 构造函数
func NewOauthClientHandler(svc service.OauthClientService) *OauthClientHandler {
	return &OauthClientHandler{
		crud: base.NewBaseHandler[model.OauthClient, *dto.OauthClientQuery, dto.CreateOauthClientDTO, dto.UpdateOauthClientDTO](svc),
		svc:  svc,
	}
}

// List 查询终端应用列表。
//
// @Summary      查询终端应用列表
// @Tags         终端应用
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  response.Response
// @Router       /api/v1/sys/oauth/client/list [get]
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
// @Router       /api/v1/sys/oauth/client/page [get]
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
// @Router       /api/v1/sys/oauth/client/{id} [get]
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
// @Router       /api/v1/sys/oauth/client/{ids} [delete]
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
// @Router       /api/v1/sys/oauth/client [post]
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
// @Router       /api/v1/sys/oauth/client/{id} [put]
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

// RefreshGatewayCache 刷新网关 OAuth Client 资源缓存。
//
// @Summary      刷新网关 OAuth Client 资源缓存
// @Tags         终端应用
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  response.Response
// @Router       /api/v1/sys/oauth/client/refresh-cache [post]
func (h *OauthClientHandler) RefreshGatewayCache(c *gin.Context) {
	list, err := h.svc.List(&dto.OauthClientQuery{})
	if err != nil {
		response.FailBy(c, apperr.ErrQueryFailed, map[string]any{"error": err.Error()})
		return
	}

	rules := make(map[string]gatewayOauthClientRule, len(list))
	for _, item := range list {
		clientID := strings.TrimSpace(item.ClientId)
		if clientID == "" {
			continue
		}

		rule := gatewayOauthClientRule{
			ClientID:     clientID,
			Enabled:      true,
			ServiceNames: map[string]struct{}{},
			ExactRules:   map[string]struct{}{},
			PrefixRules:  []string{},
		}

		for _, resource := range splitGatewayResources(item.Resources) {
			if strings.EqualFold(resource, "ALL") {
				rule.AllowAll = true
				continue
			}
			rule.ServiceNames[resource] = struct{}{}
		}
		rules[clientID] = rule
	}

	if err := frameworkredis.SetJSON(c.Request.Context(), gatewayOauthClientCacheKey, rules, withGatewayCacheJitter(gatewayOauthClientCacheTTL)); err != nil {
		response.Fail(c, http.StatusInternalServerError, 1013, "err.api.failed_to_refresh_cache", map[string]any{"error": err.Error()})
		return
	}
	if err := frameworkredis.SetJSON(c.Request.Context(), gatewayOauthClientBackupCacheKey, rules, withGatewayCacheJitter(gatewayOauthClientBackupTTL)); err != nil {
		response.Fail(c, http.StatusInternalServerError, 1013, "err.api.failed_to_refresh_cache", map[string]any{"error": err.Error()})
		return
	}

	response.Success(c, gin.H{
		"count":        len(rules),
		"refreshed_at": time.Now().Format(time.RFC3339),
	})
}

func splitGatewayResources(raw string) []string {
	if strings.TrimSpace(raw) == "" {
		return nil
	}
	parts := strings.Split(raw, ",")
	result := make([]string, 0, len(parts))
	for _, item := range parts {
		item = strings.TrimSpace(item)
		if item != "" {
			result = append(result, item)
		}
	}
	return result
}

func withGatewayCacheJitter(base time.Duration) time.Duration {
	if base <= 0 || gatewayOauthClientCacheJitter <= 0 {
		return base
	}
	return base + time.Duration(rand.Int63n(int64(gatewayOauthClientCacheJitter)))
}
