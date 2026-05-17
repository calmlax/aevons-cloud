package handler

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"auth-service/internal/service"

	pkgauth "github.com/calmlax/aevons-framework/auth"
	"github.com/calmlax/aevons-framework/response"

	"github.com/gin-gonic/gin"
)

// PasskeyHandler WebAuthn Passkey 处理器。
type PasskeyHandler struct {
	svc     service.PasskeyService
	authSvc service.AuthService
}

func NewPasskeyHandler(svc service.PasskeyService, authSvc service.AuthService) *PasskeyHandler {
	return &PasskeyHandler{svc: svc, authSvc: authSvc}
}

// BeginRegistration 开始 Passkey 注册。
//
// @Summary      开始 Passkey 注册
// @Description  为当前登录用户生成 WebAuthn 注册 challenge，返回 PublicKeyCredentialCreationOptions。需要 Bearer Token 认证。
// @Tags         Passkey
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  response.Response
// @Failure      401  {object}  response.Response
// @Router       /auth/passkey/register/begin [post]
func (h *PasskeyHandler) BeginRegistration(c *gin.Context) {
	userId, err := pkgauth.GetCurrentUserId(c.Request.Context())
	if err != nil {
		response.Fail(c, http.StatusUnauthorized, http.StatusUnauthorized, "err.sys.unauthorized")
		return
	}

	optionsJSON, sessionKey, err := h.svc.BeginRegistration(c.Request.Context(), userId)
	if err != nil {
		response.FailServerError(c, "err.sys.server_error", map[string]any{"error": err.Error()})
		return
	}

	response.Success(c, gin.H{
		"options":     string(optionsJSON),
		"session_key": sessionKey,
	})
}

// FinishRegistration 完成 Passkey 注册。
//
// @Summary      完成 Passkey 注册
// @Description  验证浏览器返回的 PublicKeyCredential，保存凭据。需要 Bearer Token 认证。
// @Tags         Passkey
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        session_key  body  string  true  "BeginRegistration 返回的 session_key"
// @Param        response     body  string  true  "浏览器 navigator.credentials.create() 返回的 JSON"
// @Success      200  {object}  response.Response
// @Failure      400  {object}  response.Response
// @Failure      401  {object}  response.Response
// @Router       /auth/passkey/register/finish [post]
func (h *PasskeyHandler) FinishRegistration(c *gin.Context) {
	userId, err := pkgauth.GetCurrentUserId(c.Request.Context())
	if err != nil {
		response.Fail(c, http.StatusUnauthorized, http.StatusUnauthorized, "err.sys.unauthorized")
		return
	}

	var body struct {
		SessionKey   string          `json:"session_key" binding:"required"`
		ResponseJSON json.RawMessage `json:"response"    binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		response.Fail(c, http.StatusBadRequest, http.StatusBadRequest, "err.sys.bad_request", map[string]any{"error": err.Error()})
		return
	}

	if err := h.svc.FinishRegistration(c.Request.Context(), userId, body.SessionKey, []byte(body.ResponseJSON)); err != nil {
		response.Fail(c, http.StatusBadRequest, http.StatusBadRequest, "auth.passkey_register_failed", map[string]any{"error": err.Error()})
		return
	}

	response.Success(c, nil)
}

// BeginAuthentication 开始 Passkey 认证。
//
// @Summary      开始 Passkey 认证
// @Description  生成 WebAuthn 认证 challenge，返回 PublicKeyCredentialRequestOptions（支持 discoverable credential，username 可选）。
// @Tags         Passkey
// @Accept       json
// @Produce      json
// @Param        username  body  string  false  "用户名（可选，discoverable credential 无需传）"
// @Success      200  {object}  response.Response
// @Failure      500  {object}  response.Response
// @Router       /auth/passkey/login/begin [post]
func (h *PasskeyHandler) BeginAuthentication(c *gin.Context) {
	var body struct {
		Username string `json:"username"`
	}
	// 忽略绑定错误，username 可选
	_ = c.ShouldBindJSON(&body)

	optionsJSON, sessionKey, err := h.svc.BeginAuthentication(c.Request.Context(), body.Username)
	if err != nil {
		response.FailServerError(c, "err.sys.server_error", map[string]any{"error": err.Error()})
		return
	}

	response.Success(c, gin.H{
		"options":     string(optionsJSON),
		"session_key": sessionKey,
	})
}

// FinishAuthentication 完成 Passkey 认证。
//
// @Summary      完成 Passkey 认证
// @Description  验证浏览器返回的 PublicKeyCredential，颁发令牌对。
// @Tags         Passkey
// @Accept       json
// @Produce      json
// @Param        session_key  body  string  true  "BeginAuthentication 返回的 session_key"
// @Param        response     body  string  true  "浏览器 navigator.credentials.get() 返回的 JSON"
// @Success      200  {object}  response.Response
// @Failure      400  {object}  response.Response
// @Failure      401  {object}  response.Response
// @Router       /auth/passkey/login/finish [post]
func (h *PasskeyHandler) FinishAuthentication(c *gin.Context) {
	var body struct {
		SessionKey   string          `json:"session_key" binding:"required"`
		ResponseJSON json.RawMessage `json:"response"    binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		response.Fail(c, http.StatusBadRequest, http.StatusBadRequest, "err.sys.bad_request", map[string]any{"error": err.Error()})
		return
	}

	// 添加调试日志
	log.Printf("[Passkey Login] SessionKey: %s", body.SessionKey)
	log.Printf("[Passkey Login] Response JSON: %s", string(body.ResponseJSON))

	pair, err := h.svc.FinishAuthentication(
		c.Request.Context(),
		body.SessionKey,
		[]byte(body.ResponseJSON),
		c.ClientIP(),
		c.Request.UserAgent(),
	)
	if err != nil {
		log.Printf("[Passkey Login] Authentication failed: %v", err)
		if authErr, ok := err.(*service.AuthError); ok {
			response.Fail(c, authErr.HTTPStatus, authErr.HTTPStatus, authErr.Code)
			return
		}
		response.Fail(c, http.StatusUnauthorized, http.StatusUnauthorized, "auth.passkey_verify_failed")
		return
	}

	log.Printf("[Passkey Login] Authentication successful")
	respondWithTokenPair(c, pair)
}

// ListCredentials 列出当前用户的所有 Passkey 凭据。
//
// @Summary      列出 Passkey 凭据
// @Description  返回当前用户注册的所有 Passkey 凭据列表。需要 Bearer Token 认证。
// @Tags         Passkey
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  response.Response
// @Failure      401  {object}  response.Response
// @Router       /auth/passkey/credentials [get]
func (h *PasskeyHandler) ListCredentials(c *gin.Context) {
	userId, err := pkgauth.GetCurrentUserId(c.Request.Context())
	if err != nil {
		response.Fail(c, http.StatusUnauthorized, http.StatusUnauthorized, "err.sys.unauthorized")
		return
	}

	list, err := h.svc.ListCredentials(c.Request.Context(), userId)
	if err != nil {
		response.FailServerError(c, "err.sys.server_error", map[string]any{"error": err.Error()})
		return
	}
	response.Success(c, list)
}

// RevokeCredential 吊销指定 Passkey 凭据。
//
// @Summary      吊销 Passkey 凭据
// @Description  吊销当前用户的指定凭据，吊销后该凭据无法再用于登录。需要 Bearer Token 认证。
// @Tags         Passkey
// @Produce      json
// @Security     BearerAuth
// @Param        id  path  int  true  "凭据 ID"
// @Success      200  {object}  response.Response
// @Failure      400  {object}  response.Response
// @Failure      401  {object}  response.Response
// @Router       /auth/passkey/credentials/{id} [delete]
func (h *PasskeyHandler) RevokeCredential(c *gin.Context) {
	userId, err := pkgauth.GetCurrentUserId(c.Request.Context())
	if err != nil {
		response.Fail(c, http.StatusUnauthorized, http.StatusUnauthorized, "err.sys.unauthorized")
		return
	}

	var uri struct {
		Id int64 `uri:"id" binding:"required"`
	}
	if err := c.ShouldBindUri(&uri); err != nil {
		response.Fail(c, http.StatusBadRequest, http.StatusBadRequest, "err.sys.bad_request", map[string]any{"error": err.Error()})
		return
	}

	if err := h.svc.RevokeCredential(c.Request.Context(), userId, uri.Id); err != nil {
		response.Fail(c, http.StatusBadRequest, http.StatusBadRequest, "auth.passkey_revoke_failed", map[string]any{"error": err.Error()})
		return
	}
	response.Success(c, nil)
}

// 确保 io 包被使用（避免 unused import）
var _ = io.EOF
