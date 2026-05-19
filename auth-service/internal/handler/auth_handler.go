package handler

import (
	"encoding/base64"
	"net/http"
	"strings"

	"auth-service/internal/dto"
	"auth-service/internal/service"

	authctx "github.com/calmlax/aevons-framework/auth/context"
	authmodel "github.com/calmlax/aevons-framework/auth/model"
	"github.com/calmlax/aevons-framework/consts"
	"github.com/calmlax/aevons-framework/response"

	"github.com/gin-gonic/gin"
)

// AuthHandler 认证模块 HTTP 处理器。
type AuthHandler struct {
	svc service.AuthService
}

// NewAuthHandler 创建 AuthHandler 实例。
func NewAuthHandler(svc service.AuthService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

// Login 处理多模式 OAuth2 登录。
//
// @Summary      用户登录
// @Description  支持 password / email / authorization_code 多种 grant_type。客户端凭证通过 Authorization: Basic base64(client_id:client_secret) 传递。
// @Tags         认证
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string  false  "Basic base64(client_id:client_secret)"
// @Param        grant_type     body      string  true   "授权类型：password | email | authorization_code"
// @Param        username       body      string  false  "用户名（password 模式）"
// @Param        password       body      string  false  "密码（password 模式）"
// @Param        key_id         body      string  false  "RSA 加密时传入的 key_id"
// @Param        email          body      string  false  "邮箱（email 模式）"
// @Param        code           body      string  false  "邮箱验证码或授权码"
// @Success      200  {object}  response.Response
// @Failure      400  {object}  response.Response
// @Failure      401  {object}  response.Response
// @Router       /api/auth/v1/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req authmodel.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, http.StatusBadRequest, "err.sys.bad_request", map[string]any{"error": err.Error()})
		return
	}

	req.ClientIP = c.ClientIP()
	req.UserAgent = c.Request.UserAgent()

	if req.GrantType == "password" || req.GrantType == "email" {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Basic ") {
			token, err := h.svc.AuthorizeLogin(c.Request.Context(), &req)
			if err != nil {
				handleAuthError(c, err)
				return
			}
			response.Success(c, gin.H{"access_token": token})
			return
		}
	}

	clientId, clientSecret, ok := extractBasicAuth(c)
	if !ok {
		clientId = ""
		clientSecret = ""
	}
	req.ClientId = clientId
	req.ClientSecret = clientSecret

	pair, err := h.svc.Login(c.Request.Context(), &req)
	if err != nil {
		handleAuthError(c, err)
		return
	}
	respondWithTokenPair(c, pair)
}

// Refresh 使用 Refresh Token 换取新令牌对。
//
// @Summary      刷新令牌
// @Description  Refresh Token 从 HttpOnly Cookie（refreshToken）中读取，客户端凭证通过 Basic Auth 传递。
// @Tags         认证
// @Produce      json
// @Param        Authorization  header    string  true  "Basic base64(client_id:client_secret)"
// @Success      200  {object}  response.Response
// @Failure      401  {object}  response.Response
// @Router       /api/auth/v1/refresh [post]
func (h *AuthHandler) Refresh(c *gin.Context) {
	clientId, _, ok := extractBasicAuth(c)
	if !ok {
		response.Fail(c, http.StatusUnauthorized, http.StatusUnauthorized, consts.ErrOAuthInvalidClient)
		return
	}

	refreshToken, err := c.Cookie("refreshToken")
	if err != nil || refreshToken == "" {
		response.Fail(c, http.StatusUnauthorized, http.StatusUnauthorized, consts.ErrTokenMissing)
		return
	}

	pair, err := h.svc.Refresh(c.Request.Context(), refreshToken, clientId)
	if err != nil {
		handleAuthError(c, err)
		return
	}
	respondWithTokenPair(c, pair)
}

// Logout 撤销当前用户的令牌对。
//
// @Summary      退出登录
// @Description  撤销 Access Token 和 Refresh Token，清除 Cookie。需要 Bearer Token 认证。
// @Tags         认证
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  response.Response
// @Failure      401  {object}  response.Response
// @Router       /api/auth/v1/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	token := extractBearerToken(c)
	if token == "" {
		response.Fail(c, http.StatusUnauthorized, http.StatusUnauthorized, consts.ErrTokenMissing)
		return
	}

	if err := h.svc.Logout(c.Request.Context(), token); err != nil {
		handleAuthError(c, err)
		return
	}
	c.SetCookie("refreshToken", "", -1, "/api/auth/v1/refresh", "", true, true)
	response.Success(c, nil)
}

// SendEmailCode 发送邮箱验证码。
//
// @Summary      发送邮箱验证码
// @Description  生成 6 位验证码并发送到指定邮箱，purpose 区分用途（register / reset）。
// @Tags         认证
// @Accept       json
// @Produce      json
// @Param        email    body      string  true  "目标邮箱"
// @Param        purpose  body      string  true  "用途：register | reset"
// @Success      200  {object}  response.Response
// @Failure      400  {object}  response.Response
// @Router       /api/auth/v1/email/code [post]
func (h *AuthHandler) SendEmailCode(c *gin.Context) {
	var body struct {
		Email   string `json:"email" binding:"required,email"`
		Purpose string `json:"purpose" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		response.Fail(c, http.StatusBadRequest, http.StatusBadRequest, "err.sys.bad_request", map[string]any{"error": err.Error()})
		return
	}

	if err := h.svc.SendEmailCode(c.Request.Context(), body.Email, body.Purpose); err != nil {
		response.FailServerError(c, "err.sys.server_error", map[string]any{"error": err.Error()})
		return
	}
	response.Success(c, nil)
}

// Register 用户注册。
//
// @Summary      用户注册
// @Description  通过邮箱验证码完成注册，返回新建用户信息（密码字段已屏蔽）。
// @Tags         认证
// @Accept       json
// @Produce      json
// @Param        name      body      string  true  "姓名"
// @Param        email     body      string  true  "邮箱"
// @Param        password  body      string  true  "密码"
// @Param        code      body      string  true  "邮箱验证码"
// @Success      200  {object}  response.Response
// @Failure      400  {object}  response.Response
// @Router       /api/auth/v1/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req struct {
		Name     string `json:"name" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
		Code     string `json:"code" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, http.StatusBadRequest, "err.sys.bad_request", map[string]any{"error": err.Error()})
		return
	}

	user, err := h.svc.Register(c.Request.Context(), req.Name, req.Email, req.Password, req.Code)
	if err != nil {
		handleAuthError(c, err)
		return
	}
	user.Password = ""
	response.Success(c, user)
}

// ResetPassword 重置密码。
//
// @Summary      重置密码
// @Description  通过邮箱验证码验证身份后设置新密码。
// @Tags         认证
// @Accept       json
// @Produce      json
// @Param        email     body      string  true  "邮箱"
// @Param        password  body      string  true  "新密码"
// @Param        code      body      string  true  "邮箱验证码"
// @Success      200  {object}  response.Response
// @Failure      400  {object}  response.Response
// @Router       /api/auth/v1/reset-password [post]
func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
		Code     string `json:"code" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, http.StatusBadRequest, "err.sys.bad_request", map[string]any{"error": err.Error()})
		return
	}

	if err := h.svc.ResetPassword(c.Request.Context(), req.Email, req.Password, req.Code); err != nil {
		handleAuthError(c, err)
		return
	}
	response.Success(c, nil)
}

// GetLatestLoginLog 获取当前登录用户最近的登录日志。
//
// @Summary      查询最近登录日志
// @Description  返回当前登录用户最近 10 条登录日志记录。
// @Tags         认证
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  response.Response
// @Failure      401  {object}  response.Response
// @Failure      500  {object}  response.Response
// @Router       /api/auth/v1/user/login-logs [get]
func (h *AuthHandler) GetLatestLoginLog(c *gin.Context) {
	entries, err := h.svc.GetLatestLoginLog(c.Request.Context())
	if err != nil {
		handleAuthError(c, err)
		return
	}
	response.Success(c, entries)
}

// GenerateAuthCode 为当前用户生成 SSO 授权码。
//
// @Summary      生成授权码
// @Description  为当前登录用户和指定客户端生成一次性授权码（SSO 入口）。需要 Bearer Token 认证。
// @Tags         认证
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        client_id  body      string  true  "客户端 ID"
// @Success      200  {object}  response.Response
// @Failure      400  {object}  response.Response
// @Failure      401  {object}  response.Response
// @Router       /api/auth/v1/code [post]
func (h *AuthHandler) GenerateAuthCode(c *gin.Context) {
	userId, err := authctx.GetCurrentUserId(c.Request.Context())
	if err != nil {
		response.Fail(c, http.StatusUnauthorized, http.StatusUnauthorized, consts.ErrTokenMissing)
		return
	}

	var body struct {
		ClientId string `json:"client_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		response.Fail(c, http.StatusBadRequest, http.StatusBadRequest, "err.sys.bad_request", map[string]any{"error": err.Error()})
		return
	}

	code, err := h.svc.GenerateAuthCode(c.Request.Context(), userId, body.ClientId)
	if err != nil {
		response.FailServerError(c, "err.sys.server_error", map[string]any{"error": err.Error()})
		return
	}
	response.Success(c, gin.H{"code": code})
}

// Authorize 发起 OAuth2 授权码流程。
//
// @Summary      OAuth2 授权页信息
// @Description  校验客户端合法性，生成 state，返回授权页所需的客户端展示信息。
// @Tags         OAuth2
// @Produce      json
// @Param        client_id      query     string  true   "客户端 ID"
// @Param        redirect_uri   query     string  false  "回调地址"
// @Param        response_type  query     string  false  "固定为 code"
// @Param        state          query     string  false  "客户端随机状态值"
// @Success      200  {object}  response.Response
// @Failure      400  {object}  response.Response
// @Router       /api/auth/v1/authorize [get]
func (h *AuthHandler) Authorize(c *gin.Context) {
	clientId := c.Query("client_id")
	if clientId == "" {
		response.Fail(c, http.StatusBadRequest, http.StatusBadRequest, consts.ErrOAuthInvalidClient)
		return
	}
	redirectURI := c.Query("redirect_uri")

	info, err := h.svc.Authorize(c.Request.Context(), clientId, redirectURI)
	if err != nil {
		handleAuthError(c, err)
		return
	}

	response.Success(c, gin.H{
		"client_id":    info.ClientId,
		"client_name":  info.ClientName,
		"logo_uri":     info.LogoURI,
		"scope":        info.Scope,
		"redirect_uri": info.RedirectURI,
		"state":        info.State,
		"autoapprove":  info.AutoApprove,
	})
}

// ApproveAuthorize 用户确认授权，生成授权码。
//
// @Summary      用户确认 OAuth2 授权
// @Description  用户在授权页点击同意后调用，生成授权码并返回回调 URL。
// @Tags         OAuth2
// @Accept       json
// @Produce      json
// @Param        access_token  body      string    true   "用户 Access Token"
// @Param        state         body      string    true   "state 值"
// @Param        scopes        body      []string  false  "授权范围列表"
// @Success      200  {object}  response.Response
// @Failure      400  {object}  response.Response
// @Failure      401  {object}  response.Response
// @Router       /api/auth/v1/authorize [post]
func (h *AuthHandler) ApproveAuthorize(c *gin.Context) {
	var body struct {
		AccessToken string   `json:"access_token" binding:"required"`
		State       string   `json:"state" binding:"required"`
		Scopes      []string `json:"scopes"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		response.Fail(c, http.StatusBadRequest, http.StatusBadRequest, "err.sys.bad_request", map[string]any{"error": err.Error()})
		return
	}

	loginUser, err := h.svc.GetLoginUser(c.Request.Context(), body.AccessToken)
	if err != nil {
		response.Fail(c, http.StatusUnauthorized, http.StatusUnauthorized, consts.ErrTokenExpired)
		return
	}

	callbackURL, err := h.svc.ApproveAuthorize(c.Request.Context(), loginUser.UserId, body.State, body.Scopes)
	if err != nil {
		handleAuthError(c, err)
		return
	}
	response.Success(c, gin.H{"redirect_uri": callbackURL})
}

// Callback 处理 OAuth2 授权回调。
//
// @Summary      OAuth2 授权回调
// @Description  校验 state 和授权码，换取令牌对。
// @Tags         OAuth2
// @Produce      json
// @Param        code   query     string  true  "授权码"
// @Param        state  query     string  true  "state 值"
// @Success      200  {object}  response.Response
// @Failure      400  {object}  response.Response
// @Router       /api/auth/v1/callback [get]
func (h *AuthHandler) Callback(c *gin.Context) {
	code := c.Query("code")
	state := c.Query("state")

	if code == "" || state == "" {
		response.Fail(c, http.StatusBadRequest, http.StatusBadRequest, consts.ErrInvalidAuthCode)
		return
	}

	pair, err := h.svc.Callback(c.Request.Context(), code, state)
	if err != nil {
		handleAuthError(c, err)
		return
	}
	respondWithTokenPair(c, pair)
}

// GetPublicKey 获取 RSA 公钥。
//
// @Summary      获取 RSA 公钥
// @Description  返回用于前端加密密码的 RSA 公钥（PEM 格式）和对应的 key_id。
// @Tags         认证
// @Produce      json
// @Success      200  {object}  response.Response
// @Failure      500  {object}  response.Response
// @Router       /api/auth/v1/public-key [get]
func (h *AuthHandler) GetPublicKey(c *gin.Context) {
	resp, err := h.svc.GetPublicKey(c.Request.Context())
	if err != nil {
		response.FailServerError(c, "err.sys.server_error", map[string]any{"error": err.Error()})
		return
	}
	response.Success(c, resp)
}

// Routers 获取当前用户的动态路由菜单。
//
// @Summary      获取动态路由菜单
// @Description  根据当前登录用户的角色权限返回前端路由菜单树。需要 Bearer Token 认证。
// @Tags         认证
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  response.Response
// @Failure      401  {object}  response.Response
// @Router       /api/auth/v1/routers [get]
func (h *AuthHandler) Routers(c *gin.Context) {
	langCode := c.GetHeader(consts.AcceptLanguage)
	routers, err := h.svc.GetRouters(c.Request.Context(), langCode)
	if err != nil {
		handleAuthError(c, err)
		return
	}
	response.Success(c, routers)
}

// GetUserInfo 获取当前登录用户信息。
//
// @Summary      获取当前用户信息
// @Description  返回当前 Token 对应的用户基本信息（角色、权限等）。需要 Bearer Token 认证。
// @Tags         认证
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  response.Response
// @Failure      401  {object}  response.Response
// @Router       /api/auth/v1/user [get]
func (h *AuthHandler) GetUserInfo(c *gin.Context) {
	user, err := authctx.GetCurrentUser(c.Request.Context())
	if err != nil {
		handleAuthError(c, err)
		return
	}
	response.Success(c, user)
}

// GetProfile 获取当前登录用户的完整个人档案。
//
// @Summary      获取个人档案
// @Description  返回当前用户的完整档案，包含部门、岗位、登录日志等。需要 Bearer Token 认证。
// @Tags         认证
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  response.Response
// @Failure      401  {object}  response.Response
// @Router       /api/auth/v1/user/profile [get]
func (h *AuthHandler) GetProfile(c *gin.Context) {
	userId, err := authctx.GetCurrentUserId(c.Request.Context())
	if err != nil {
		response.Fail(c, http.StatusUnauthorized, 401, "err.sys.unauthorized")
		return
	}

	profile, err := h.svc.GetProfile(userId)
	if err != nil {
		response.Fail(c, http.StatusInternalServerError, 1003, "err.api.failed_to_fetch_user")
		return
	}
	response.Success(c, profile)
}

// UpdateProfile 更新个人资料。
//
// @Summary      更新个人资料
// @Description  更新当前登录用户的昵称、手机号、性别等基本信息。需要 Bearer Token 认证。
// @Tags         认证
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body  body      dto.UpdateProfileDTO  true  "更新内容"
// @Success      200  {object}  response.Response
// @Failure      400  {object}  response.Response
// @Failure      401  {object}  response.Response
// @Router       /api/auth/v1/user/profile [put]
func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	userId, err := authctx.GetCurrentUserId(c.Request.Context())
	if err != nil {
		handleAuthError(c, err)
		return
	}

	accessToken := extractBearerToken(c)

	var req dto.UpdateProfileDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, http.StatusBadRequest, "err.sys.bad_request", map[string]any{"error": err.Error()})
		return
	}

	if err := h.svc.UpdateProfile(c.Request.Context(), userId, accessToken, &req); err != nil {
		handleAuthError(c, err)
		return
	}
	response.Success(c, nil)
}

// UpdatePassword 修改密码。
//
// @Summary      修改密码
// @Description  登录用户修改自己的密码，成功后所有端的 Token 立即失效。需要 Bearer Token 认证。
// @Tags         认证
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        currentPassword  body      string  true  "当前密码"
// @Param        nextPassword     body      string  true  "新密码"
// @Success      200  {object}  response.Response
// @Failure      400  {object}  response.Response
// @Failure      401  {object}  response.Response
// @Router       /api/auth/v1/user/password [put]
func (h *AuthHandler) UpdatePassword(c *gin.Context) {
	userId, err := authctx.GetCurrentUserId(c.Request.Context())
	if err != nil {
		handleAuthError(c, err)
		return
	}

	var req authmodel.UpdatePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, http.StatusBadRequest, http.StatusBadRequest, "err.sys.bad_request", map[string]any{"error": err.Error()})
		return
	}

	if err := h.svc.UpdatePassword(c.Request.Context(), userId, &req); err != nil {
		handleAuthError(c, err)
		return
	}
	response.Success(c, nil)
}

// ── 内部工具函数 ──────────────────────────────────────────────────────────────

// extractBasicAuth 从 "Authorization: Basic base64(client_id:client_secret)" 头中解析客户端凭证。
func extractBasicAuth(c *gin.Context) (clientId, clientSecret string, ok bool) {
	header := c.GetHeader("Authorization")
	if !strings.HasPrefix(header, "Basic ") {
		return "", "", false
	}
	encoded := strings.TrimPrefix(header, "Basic ")
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", "", false
	}
	parts := strings.SplitN(string(decoded), ":", 2)
	if len(parts) == 0 || parts[0] == "" {
		return "", "", false
	}
	clientId = parts[0]
	if len(parts) == 2 {
		clientSecret = parts[1]
	}
	return clientId, clientSecret, true
}

// extractBearerToken 从 "Authorization: Bearer <token>" 请求头中提取令牌。
func extractBearerToken(c *gin.Context) string {
	header := c.GetHeader("Authorization")
	if !strings.HasPrefix(header, "Bearer ") {
		return ""
	}
	return strings.TrimPrefix(header, "Bearer ")
}

// handleAuthError 将 service.AuthError 映射为对应的 HTTP 响应。
func handleAuthError(c *gin.Context, err error) {
	if authErr, ok := err.(*service.AuthError); ok {
		response.Fail(c, authErr.HTTPStatus, authErr.HTTPStatus, authErr.Code)
		return
	}
	response.FailServerError(c, "err.sys.server_error", map[string]any{"error": err.Error()})
}

// respondWithTokenPair 安全分发双 Token 机制响应。
// Access Token 放 JSON 返回，Refresh Token 设为 HttpOnly, Secure Cookie。
func respondWithTokenPair(c *gin.Context, pair *authmodel.TokenPair) {
	if pair.RefreshToken != "" {
		c.SetSameSite(http.SameSiteLaxMode)
		c.SetCookie("refreshToken", pair.RefreshToken, int(pair.RefreshExpiresIn), "/api/auth/v1/refresh", "", true, true)
	}
	response.Success(c, gin.H{
		"access_token": pair.AccessToken,
		"expires_in":   pair.ExpiresIn,
		"token_type":   pair.TokenType,
		"scope":        pair.Scope,
	})
}
