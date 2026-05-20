// Package service 实现认证模块的业务逻辑。
package service

import (
	"auth-service/internal/dto"
	authRepo "auth-service/internal/repository"
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"internal-grpc/log_grpc"
	"math/big"
	"strings"
	"time"

	"auth-service/internal/model"

	authctx "github.com/calmlax/aevons-framework/auth/context"
	authmodel "github.com/calmlax/aevons-framework/auth/model"
	authnotifier "github.com/calmlax/aevons-framework/auth/notifier"
	authstore "github.com/calmlax/aevons-framework/auth/store"
	"github.com/calmlax/aevons-framework/config"
	"github.com/calmlax/aevons-framework/consts"
	"github.com/calmlax/aevons-framework/utils"

	"github.com/google/uuid"
	"github.com/mileusna/useragent"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// AuthorizeInfo 授权页展示所需的客户端信息。
type AuthorizeInfo struct {
	ClientId    string
	ClientName  string
	LogoURI     string
	Scope       string
	RedirectURI string
	State       string
	AutoApprove bool // true 时前端直接授权，无需用户手动确认
}

// AuthService 定义认证业务逻辑接口。
type AuthService interface {
	AuthorizeLogin(ctx context.Context, req *authmodel.LoginRequest) (string, error)
	Login(ctx context.Context, req *authmodel.LoginRequest) (*authmodel.TokenPair, error)
	Refresh(ctx context.Context, refreshToken, clientId string) (*authmodel.TokenPair, error)
	Logout(ctx context.Context, accessToken string) error
	// GlobalLogout 全局登出：清除该用户所有客户端会话并发送 SLO 通知。
	GlobalLogout(ctx context.Context, userId int64) error
	// UpdateProfile 修改当前登录用户的档案信息
	UpdateProfile(ctx context.Context, userId int64, accessToken string, req *dto.UpdateProfileDTO) error
	// UpdatePassword 修改当前登录用户密码
	UpdatePassword(ctx context.Context, userId int64, req *authmodel.UpdatePasswordRequest) error

	SendEmailCode(ctx context.Context, email, purpose string) error
	Register(ctx context.Context, name, email, password, code string) (*model.User, error)
	ResetPassword(ctx context.Context, email, password, code string) error
	GenerateAuthCode(ctx context.Context, userId int64, clientId string) (string, error)
	// Authorize 校验客户端参数，生成 state 存入 Redis，返回授权页展示信息。
	Authorize(ctx context.Context, clientId, redirectURI string) (*AuthorizeInfo, error)
	// ApproveAuthorize 用户确认授权：校验 state，生成授权码，返回回调 URL。
	ApproveAuthorize(ctx context.Context, userId int64, state string, scopes []string) (callbackURL string, err error)
	// GetLoginUser 根据 access token 获取登录用户信息。
	GetLoginUser(ctx context.Context, accessToken string) (*authmodel.LoginUser, error)
	// ValidateUserPassword 校验用户名密码，返回 userId。
	ValidateUserPassword(ctx context.Context, username, password, keyId string) (int64, error)
	// IssueShortToken 为指定用户颁发短期 access token（授权页登录用）。
	IssueShortToken(ctx context.Context, userId int64) (string, error)
	// Callback 处理授权回调：校验 state，使用 oauth2.Config.Exchange 换取令牌（Requirements 1.2, 1.4, 6.1, 6.2）。
	Callback(ctx context.Context, code, state string) (*authmodel.TokenPair, error)
	// GetRouters 获取当前用户的动态路由菜单。
	GetRouters(ctx context.Context, langCode string) ([]authmodel.Menu, error)
	GetPublicKey(ctx context.Context) (*authmodel.PublicKeyResponse, error)
	GetProfile(id int64) (*dto.UserProfile, error)
	// LoginByUserId 直接为指定用户颁发令牌对（Passkey 等无密码登录使用）
	LoginByUserId(ctx context.Context, userId int64, clientId string) (*authmodel.TokenPair, error)
	// RecordLoginLog 记录登录日志（供外部服务调用）
	RecordLoginLog(ctx context.Context, username string, clientId string, grantType string, status int16, msg, userAgent, ip string)
	// GetLatestLoginLog 查询指定用户最近的登录日志（供外部服务调用）
	GetLatestLoginLog(ctx context.Context) ([]*log_grpc.LoginEntry, error)
}

type authService struct {
	store    authstore.TokenStore
	authRepo authRepo.AuthRepository
	cfg      config.AuthConfig
	notifier authnotifier.SLONotifier
	logStore log_grpc.LoginLogStore
}

// NewAuthService 创建 AuthService 实例。
func NewAuthService(
	store authstore.TokenStore,
	authRepo authRepo.AuthRepository,
	cfg config.AuthConfig,
	notifier authnotifier.SLONotifier,
	logStore log_grpc.LoginLogStore,
) AuthService {
	return &authService{
		store:    store,
		authRepo: authRepo,
		cfg:      cfg,
		notifier: notifier,
		logStore: logStore,
	}
}

func (s *authService) recordLoginLog(ctx context.Context, username, clientId, grantType string, status int16, msg, uaString, ip string) {
	ua := useragent.Parse(uaString)

	os := ua.OS
	if os == "" {
		os = "Unknown"
	}
	browser := ua.Name
	if browser == "" {
		browser = "Unknown"
	}

	log := &log_grpc.LoginEntry{
		Username:  username,
		ClientID:  clientId,
		GrantType: grantType,
		IP:        ip,
		Browser:   browser,
		OS:        os,
		Status:    status,
		Msg:       msg,
		LoginAt:   time.Now(),
	}

	go func(record *log_grpc.LoginEntry) {
		writeCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := s.logStore.WriteLoginLog(writeCtx, *record); err != nil {
			fmt.Printf("[AuthService] write login log failed: %v\n", err)
		}
	}(log)
}

// ValidateEmailCode 校验邮箱验证码，返回 userId。
func (s *authService) ValidateEmailCode(ctx context.Context, email, code string) (int64, error) {
	storedCode, err := s.store.GetEmailCode(ctx, email, "login")
	if err != nil {
		return 0, &AuthError{Code: consts.ErrInvalidCode, HTTPStatus: 400}
	}
	if storedCode != code {
		return 0, &AuthError{Code: consts.ErrInvalidCode, HTTPStatus: 400}
	}
	// 删除验证码，确保一次性使用
	if err := s.store.DeleteEmailCode(ctx, email, "login"); err != nil {
		return 0, &AuthError{Code: consts.ErrInvalidCode, HTTPStatus: 400}
	}

	user, err := s.authRepo.GetUserByEmail(email)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return 0, &AuthError{Code: consts.ErrUserNotFound, HTTPStatus: 404}
		}
		return 0, err
	}
	if user.Status != 0 {
		return 0, &AuthError{Code: consts.ErrAccountDisabled, HTTPStatus: 403}
	}
	return user.Id, nil
}

// AuthorizeLogin 处理授权页内部登录（只验证凭据，返回短期 token）。
func (s *authService) AuthorizeLogin(ctx context.Context, req *authmodel.LoginRequest) (string, error) {
	var userId int64
	var err error

	logUsername := req.Username
	if logUsername == "" {
		logUsername = req.Email
	}

	if req.GrantType == "email" {
		userId, err = s.ValidateEmailCode(ctx, req.Email, req.Code)
	} else {
		userId, err = s.ValidateUserPassword(ctx, req.Username, req.Password, req.KeyId)
	}

	if err != nil {
		s.recordLoginLog(ctx, logUsername, "", "authorize_login", 0, err.Error(), req.UserAgent, req.ClientIP)
		return "", err
	}
	// 颁发短期 token（10分钟），仅用于完成授权流程
	token, err := s.IssueShortToken(ctx, userId)
	if err != nil {
		s.recordLoginLog(ctx, logUsername, "", "authorize_login", 0, err.Error(), req.UserAgent, req.ClientIP)
		return "", err
	}
	s.recordLoginLog(ctx, logUsername, "", "authorize_login", 1, "Login successful", req.UserAgent, req.ClientIP)
	return token, nil
}

// Login 先校验客户端，再根据 grant_type 分发登录请求。
func (s *authService) Login(ctx context.Context, req *authmodel.LoginRequest) (*authmodel.TokenPair, error) {
	logUsername := req.Username
	if logUsername == "" {
		logUsername = req.Email // Email fallback specifically for email
	}
	if logUsername == "" && req.GrantType == "client_credentials" {
		logUsername = req.ClientId
	}

	client, err := s.authRepo.ValidateClient(req.ClientId, req.ClientSecret, req.GrantType)
	if err != nil {
		s.recordLoginLog(ctx, logUsername, req.ClientId, req.GrantType, 0, consts.ErrOAuthInvalidClient, req.UserAgent, req.ClientIP)
		return nil, err
	}

	if client.WebServerRedirectUri != "" && req.RedirectURI != "" {
		if req.RedirectURI != client.WebServerRedirectUri {
			s.recordLoginLog(ctx, logUsername, req.ClientId, req.GrantType, 0, consts.ErrOAuthRedirectURIMismatch, req.UserAgent, req.ClientIP)
			return nil, &AuthError{Code: consts.ErrOAuthRedirectURIMismatch, HTTPStatus: 400}
		}
	}

	var pair *authmodel.TokenPair
	switch req.GrantType {
	case "password":
		pair, err = s.loginByPassword(ctx, req, client)
	case "email":
		pair, err = s.loginByEmailCode(ctx, req, client)
	case "authorization_code":
		pair, err = s.loginByAuthCode(ctx, req, client)
		// resolve username from the issued token so the log is meaningful
		if err == nil && logUsername == "" {
			if lu, e := s.store.GetLoginUser(ctx, pair.AccessToken); e == nil {
				logUsername = lu.Username
			}
		}
	case "client_credentials":
		pair, err = s.loginByClientCredentials(ctx, req, client)
	case "refresh_token":
		pair, err = s.Refresh(ctx, req.RefreshToken, req.ClientId)
	default:
		err = &AuthError{Code: consts.ErrOAuthUnsupportedGrant, HTTPStatus: 400}
	}

	if err != nil {
		s.recordLoginLog(ctx, logUsername, req.ClientId, req.GrantType, 0, err.Error(), req.UserAgent, req.ClientIP)
		return nil, err
	}

	s.recordLoginLog(ctx, logUsername, req.ClientId, req.GrantType, 1, "Login successful", req.UserAgent, req.ClientIP)
	return pair, nil
}

func (s *authService) loginByPassword(ctx context.Context, req *authmodel.LoginRequest, client *model.OauthClient) (*authmodel.TokenPair, error) {
	user, err := s.authRepo.GetUserByUsername(req.Username)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, &AuthError{Code: consts.ErrInvalidCredentials, HTTPStatus: 401}
		}
		return nil, err
	}

	password := req.Password
	if req.KeyId != "" {
		privateKey, err := s.store.GetRSAPrivateKey(ctx, req.KeyId)
		if err == nil && privateKey != "" {
			// 一次性消费
			_ = s.store.DeleteRSAPrivateKey(ctx, req.KeyId)
			decryptedHeaderBytes, err := utils.RSADecrypt(req.Password, privateKey)
			if err == nil {
				password = decryptedHeaderBytes
			}
		}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, &AuthError{Code: consts.ErrInvalidCredentials, HTTPStatus: 401}
	}

	if user.Status != 0 {
		return nil, &AuthError{Code: consts.ErrAccountDisabled, HTTPStatus: 403}
	}

	accessTTL, refreshTTL := client.GetTTL()
	return s.issueTokenPair(ctx, user.Id, req.ClientId, accessTTL, refreshTTL)
}

func (s *authService) loginByEmailCode(ctx context.Context, req *authmodel.LoginRequest, client *model.OauthClient) (*authmodel.TokenPair, error) {
	storedCode, err := s.store.GetEmailCode(ctx, req.Email, "login")
	if err != nil {
		return nil, &AuthError{Code: consts.ErrInvalidCode, HTTPStatus: 400}
	}
	if storedCode != req.Code {
		return nil, &AuthError{Code: consts.ErrInvalidCode, HTTPStatus: 400}
	}
	// 删除验证码，确保一次性使用（Requirements 3.4）
	if err := s.store.DeleteEmailCode(ctx, req.Email, "login"); err != nil {
		return nil, &AuthError{Code: consts.ErrInvalidCode, HTTPStatus: 400}
	}

	user, err := s.authRepo.GetUserByEmail(req.Email)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, &AuthError{Code: consts.ErrUserNotFound, HTTPStatus: 404}
		}
		return nil, err
	}
	if user.Status != 0 {
		return nil, &AuthError{Code: consts.ErrAccountDisabled, HTTPStatus: 403}
	}

	accessTTL, refreshTTL := client.GetTTL()
	return s.issueTokenPair(ctx, user.Id, req.ClientId, accessTTL, refreshTTL)
}

func (s *authService) loginByAuthCode(ctx context.Context, req *authmodel.LoginRequest, client *model.OauthClient) (*authmodel.TokenPair, error) {
	if req.Code == "" {
		return nil, &AuthError{Code: consts.ErrInvalidAuthCode, HTTPStatus: 400}
	}

	// 校验授权码，获取 userId 和授权 scopes
	userId, _, scopes, err := s.store.GetAuthCodeInfo(ctx, req.Code)
	if err != nil {
		return nil, &AuthError{Code: consts.ErrInvalidAuthCode, HTTPStatus: 400}
	}
	_ = s.store.DeleteAuthCode(ctx, req.Code)

	user, err := s.authRepo.GetUserByUserId(userId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, &AuthError{Code: consts.ErrUserNotFound, HTTPStatus: 404}
		}
		return nil, err
	}
	if user.Status != 0 {
		return nil, &AuthError{Code: consts.ErrAccountDisabled, HTTPStatus: 403}
	}

	accessTTL, refreshTTL := client.GetTTL()
	return s.issueTokenPairWithScopes(ctx, user.Id, req.ClientId, scopes, accessTTL, refreshTTL)
}

func (s *authService) loginByClientCredentials(ctx context.Context, req *authmodel.LoginRequest, client *model.OauthClient) (*authmodel.TokenPair, error) {
	if req.ClientSecret == "" {
		return nil, &AuthError{Code: consts.ErrOAuthInvalidClient, HTTPStatus: 401}
	}

	accessTTL, refreshTTL := client.GetTTL()
	if accessTTL <= 0 {
		accessTTL = s.cfg.AccessTokenTTL
	}
	if refreshTTL <= 0 {
		refreshTTL = s.cfg.RefreshTokenTTL
	}

	accessToken := uuid.New().String()
	refreshToken := uuid.New().String()
	sessionID := uuid.New().String()

	// 将 scope 拆分为独立权限项（Requirements 4.3）
	var permissions []string
	if client.Scope != "" {
		for _, s := range strings.Split(client.Scope, ",") {
			if trimmed := strings.TrimSpace(s); trimmed != "" {
				permissions = append(permissions, trimmed)
			}
		}
	}
	if permissions == nil {
		permissions = []string{}
	}

	loginUser := &authmodel.LoginUser{
		UserId:      0,
		Username:    client.ClientId,
		Nickname:    client.ClientName,
		Roles:       []authmodel.Role{},
		Depts:       []authmodel.Dept{},
		Permissions: permissions,
		ClientId:    client.ClientId,
	}

	if err := s.persistSessionTokens(ctx, sessionID, accessToken, refreshToken, loginUser, accessTTL, refreshTTL); err != nil {
		return nil, err
	}

	return &authmodel.TokenPair{
		AccessToken:      accessToken,
		RefreshToken:     refreshToken,
		ExpiresIn:        accessTTL,
		RefreshExpiresIn: refreshTTL,
		TokenType:        "Bearer",
	}, nil
}

// Refresh 验证 Refresh Token，撤销旧令牌对并颁发新的 TokenPair。
// 新的 Refresh Token 继承原 Token 的剩余 TTL，防止无限滑动续期。
func (s *authService) Refresh(ctx context.Context, refreshToken, clientId string) (*authmodel.TokenPair, error) {
	sessionID, err := s.store.GetSessionIDByRefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, &AuthError{Code: consts.ErrInvalidRefreshToken, HTTPStatus: 401}
	}

	// 在删除前读取剩余 TTL，保持绝对过期时间不变
	remainingTTL, err := s.store.GetRefreshTokenTTL(ctx, refreshToken)
	if err != nil || remainingTTL <= 0 {
		return nil, &AuthError{Code: consts.ErrInvalidRefreshToken, HTTPStatus: 401}
	}

	loginUser, err := s.store.GetLoginUserBySessionID(ctx, sessionID)
	if err != nil {
		return nil, &AuthError{Code: consts.ErrInvalidRefreshToken, HTTPStatus: 401}
	}
	oldAccessToken, _ := s.store.GetAccessTokenBySessionID(ctx, sessionID)

	// 校验 client_id 一致性
	if clientId != "" && loginUser.ClientId != clientId {
		return nil, &AuthError{Code: consts.ErrInvalidRefreshToken, HTTPStatus: 401}
	}

	// 先删除旧令牌对
	_ = s.store.DeleteAccessToken(ctx, oldAccessToken)
	_ = s.store.DeleteRefreshToken(ctx, refreshToken)

	// 将剩余 TTL 转换为秒，作为新 refresh token 的过期时间
	refreshTTL := int64(remainingTTL.Seconds())

	if loginUser.UserId == 0 {
		return s.reissueClientToken(ctx, sessionID, loginUser, s.cfg.AccessTokenTTL, refreshTTL)
	}
	return s.reissueUserToken(ctx, sessionID, loginUser, s.cfg.AccessTokenTTL, refreshTTL)
}

func (s *authService) reissueClientToken(ctx context.Context, sessionID string, old *authmodel.LoginUser, accessTTL, refreshTTL int64) (*authmodel.TokenPair, error) {
	accessToken := uuid.New().String()
	refreshToken := uuid.New().String()

	loginUser := &authmodel.LoginUser{
		UserId:      0,
		Username:    old.Username,
		Nickname:    old.Nickname,
		Roles:       old.Roles,
		Permissions: old.Permissions,
		ClientId:    old.ClientId,
	}

	if err := s.persistSessionTokens(ctx, sessionID, accessToken, refreshToken, loginUser, accessTTL, refreshTTL); err != nil {
		return nil, err
	}

	return &authmodel.TokenPair{
		AccessToken:      accessToken,
		RefreshToken:     refreshToken,
		ExpiresIn:        accessTTL,
		RefreshExpiresIn: refreshTTL,
		TokenType:        "Bearer",
	}, nil
}

func (s *authService) reissueUserToken(ctx context.Context, sessionID string, old *authmodel.LoginUser, accessTTL, refreshTTL int64) (*authmodel.TokenPair, error) {
	accessToken := uuid.New().String()
	refreshToken := uuid.New().String()

	if err := s.persistSessionTokens(ctx, sessionID, accessToken, refreshToken, old, accessTTL, refreshTTL); err != nil {
		return nil, err
	}

	return &authmodel.TokenPair{
		AccessToken:      accessToken,
		RefreshToken:     refreshToken,
		ExpiresIn:        accessTTL,
		RefreshExpiresIn: refreshTTL,
		TokenType:        "Bearer",
	}, nil
}

func (s *authService) persistSessionTokens(ctx context.Context, sessionID, accessToken, refreshToken string, loginUser *authmodel.LoginUser, accessTTL, refreshTTL int64) error {
	if err := s.store.SaveSession(ctx, sessionID, loginUser, time.Duration(refreshTTL)*time.Second); err != nil {
		return err
	}
	if err := s.store.SaveAccessToken(ctx, accessToken, sessionID, time.Duration(accessTTL)*time.Second); err != nil {
		return err
	}
	if err := s.store.SaveRefreshToken(ctx, refreshToken, sessionID, time.Duration(refreshTTL)*time.Second); err != nil {
		return err
	}
	return nil
}

// Logout 撤销当前会话的 Access Token 和 Refresh Token，并清除会话记录。
func (s *authService) Logout(ctx context.Context, accessToken string) error {
	loginUser, err := s.store.GetLoginUser(ctx, accessToken)
	if err == nil {
		sessionID, _ := s.store.GetSessionIDByAccessToken(ctx, accessToken)
		if refreshToken, e := s.store.GetRefreshTokenBySessionID(ctx, sessionID); e == nil && refreshToken != "" {
			_ = s.store.DeleteRefreshToken(ctx, refreshToken)
		}
		if loginUser.UserId > 0 {
			_ = s.store.RemoveUserSession(ctx, loginUser.UserId, loginUser.ClientId)
		}
		if sessionID != "" {
			_ = s.store.DeleteSession(ctx, sessionID)
		}
	}
	_ = s.store.DeleteAccessToken(ctx, accessToken)
	return nil
}

// GlobalLogout 全局登出：清除该用户所有客户端会话，并向各客户端发送 SLO 通知。
func (s *authService) GlobalLogout(ctx context.Context, userId int64) error {
	sessions, err := s.store.GetUserSessions(ctx, userId)
	if err != nil {
		// 读取失败不阻断，继续尝试
		sessions = map[string]string{}
	}

	payload := &authnotifier.LogoutPayload{
		UId:       userId,
		Action:    "logout_all",
		Timestamp: time.Now().Unix(),
	}

	for clientId, sessionID := range sessions {
		if accessToken, err := s.store.GetAccessTokenBySessionID(ctx, sessionID); err == nil && accessToken != "" {
			_ = s.store.DeleteAccessToken(ctx, accessToken)
		}
		if refreshToken, err := s.store.GetRefreshTokenBySessionID(ctx, sessionID); err == nil && refreshToken != "" {
			_ = s.store.DeleteRefreshToken(ctx, refreshToken)
		}
		_ = s.store.DeleteSession(ctx, sessionID)
		_ = s.store.RemoveUserSession(ctx, userId, clientId)

		// 发送 Back-channel Logout 通知（失败不阻断）
		// webhookURL 需从客户端配置中读取，当前以 clientId 作为占位标识记录日志
		_ = clientId
		s.notifier.Notify(ctx, "", payload)
	}

	return nil
}

// SendEmailCode 生成验证码并存入 Redis（在生产环境中应调用发信 API，此处在只在控制台打印）。
func (s *authService) SendEmailCode(ctx context.Context, email, purpose string) error {
	// 生成真正的 6 位随机验证码
	n, _ := rand.Int(rand.Reader, big.NewInt(1000000))
	code := fmt.Sprintf("%06d", n.Int64())

	// 在控制台上打印，方便本地测试提取
	fmt.Printf("[Email Service Mock] 模拟发送邮件 -> 邮箱: %s, 业务标识: %s, 验证码: %s\n", email, purpose, code)

	ttl := time.Duration(s.cfg.EmailCodeTTL) * time.Second
	return s.store.SaveEmailCode(ctx, email, purpose, code, ttl)
}

// Register 用户注册：校验验证码，创建用户记录。
func (s *authService) Register(ctx context.Context, name, email, password, code string) (*model.User, error) {
	storedCode, err := s.store.GetEmailCode(ctx, email, "register")
	if err != nil || storedCode != code {
		return nil, &AuthError{Code: "auth.invalid_code", HTTPStatus: 400}
	}
	_ = s.store.DeleteEmailCode(ctx, email, "register")

	_, err = s.authRepo.GetUserByEmail(email)
	if err == nil {
		return nil, &AuthError{Code: "auth.emailAlreadyRegistered", HTTPStatus: 400}
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Username: email,
		Nickname: name,
		Email:    email,
		Password: string(hashed),
		Status:   0, // Active
	}

	if err := s.authRepo.RegisterUser(user); err != nil {
		return nil, err
	}

	return user, nil
}

// UpdateProfile 用户修改自己资料。同步更新缓存和数据库。
func (s *authService) UpdateProfile(ctx context.Context, userId int64, accessToken string, req *dto.UpdateProfileDTO) error {
	err := s.authRepo.UpdateProfile(userId, map[string]any{
		"nickname":  req.Nickname,
		"email":     req.Email,
		"mobile":    req.Mobile,
		"avatar":    req.Avatar,
		"sex":       req.Sex,
		"autograph": req.Autograph,
	})
	if err != nil {
		return err
	}

	// 动态更新缓存，让前端页面不用重登也能看到新资料
	if accessToken != "" {
		_ = s.store.UpdateLoginUser(ctx, accessToken, func(user *authmodel.LoginUser) {
			if req.Nickname != nil {
				user.Nickname = *req.Nickname
			}
			if req.Email != nil {
				user.Email = *req.Email
			}
			if req.Avatar != nil {
				user.Avatar = *req.Avatar
			}
		})
	}
	return nil
}

// UpdatePassword 用户修改自己密码。
func (s *authService) UpdatePassword(ctx context.Context, userId int64, req *authmodel.UpdatePasswordRequest) error {
	user, err := s.authRepo.GetUserByUserId(userId)
	if err != nil {
		return err
	}

	// 验证旧密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.CurrentPassword)); err != nil {
		return &AuthError{Code: "password.currentIncorrect", HTTPStatus: 400}
	}

	// 生成新密码哈希
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.NextPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_ = s.authRepo.UpdateProfile(userId, map[string]any{"password": string(hashed)})
	if err != nil {
		return err
	}

	// 出于安全考量，密码一旦修改成功，就应强制全局登出，让所有端的令牌失效。
	// 这里注释掉，按要求先保持会话不断（或由系统决定）。为了让本例实现 "修改完被登出" 的强安全性：
	return s.GlobalLogout(ctx, userId)
}

// ResetPassword 重置密码：校验验证码，更新用户密码。
func (s *authService) ResetPassword(ctx context.Context, email, password, code string) error {
	storedCode, err := s.store.GetEmailCode(ctx, email, "reset-password")
	if err != nil || storedCode != code {
		return &AuthError{Code: "auth.invalid_code", HTTPStatus: 400}
	}
	_ = s.store.DeleteEmailCode(ctx, email, "reset-password")

	user, err := s.authRepo.GetUserByEmail(email)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &AuthError{Code: "auth.emailNotFound", HTTPStatus: 404}
		}
		return err
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_ = s.authRepo.UpdateProfile(user.Id, map[string]any{"password": string(hashed)})
	return err
}

// GenerateAuthCode 生成 UUId 授权码，关联 userId 和 clientId，存入 Redis（SSO 入口）。
func (s *authService) GenerateAuthCode(ctx context.Context, userId int64, clientId string) (string, error) {
	code := uuid.New().String()
	ttl := 10 * time.Minute
	if err := s.store.SaveAuthCode(ctx, code, userId, clientId, nil, ttl); err != nil {
		return "", err
	}
	return code, nil
}

// issueTokenPair 加载用户信息，构建 LoginUser，将令牌对存入 Redis，并记录用户会话。
func (s *authService) issueTokenPair(ctx context.Context, userId int64, clientId string, accessTTL, refreshTTL int64) (*authmodel.TokenPair, error) {
	return s.issueTokenPairWithScopes(ctx, userId, clientId, nil, accessTTL, refreshTTL)
}

// issueTokenPairWithScopes 同 issueTokenPair，额外将授权 scopes 存入 LoginUser。
func (s *authService) issueTokenPairWithScopes(ctx context.Context, userId int64, clientId string, scopes []string, accessTTL, refreshTTL int64) (*authmodel.TokenPair, error) {
	if accessTTL <= 0 {
		accessTTL = s.cfg.AccessTokenTTL
	}
	if refreshTTL <= 0 {
		refreshTTL = s.cfg.RefreshTokenTTL
	}

	user, err := s.authRepo.GetUserByUserId(userId)
	if err != nil {
		return nil, err
	}

	dbRoles, err := s.authRepo.GetRolesByUserId(userId)
	if err != nil {
		return nil, err
	}

	roleIds := make([]int64, len(dbRoles))
	for i, role := range dbRoles {
		roleIds[i] = role.Id
	}

	userDepts, err := s.authRepo.GetUserDeptsByUserId(userId)
	if err != nil {
		return nil, err
	}
	roles := make([]authmodel.Role, len(dbRoles))
	for i, role := range dbRoles {
		roles[i] = authmodel.Role{
			Id:        role.Id,
			RoleKey:   role.RoleKey,
			DataScope: role.DataScope,
			RoleName:  role.RoleName,
		}
		// 自定义数据范围
		if role.DataScope == 9 {
			deptIds, err := s.authRepo.GetRoleDeptIdsByRoleId(role.Id)
			if err != nil {
				return nil, err
			}
			roles[i].DeptIds = deptIds
		}
	}

	depts := make([]authmodel.Dept, len(userDepts))
	for i, userDept := range userDepts {
		depts[i] = authmodel.Dept{
			DeptId: userDept.DeptId,
			PostId: userDept.PostId,
		}
	}

	var permissions []string
	if userId == consts.SuperAdminId {
		permissions = []string{consts.AllPermission}
	} else {
		permissions, err = s.authRepo.GetPermissionsByRoleIds(roleIds)
		if err != nil {
			return nil, err
		}
	}

	accessToken := uuid.New().String()
	refreshToken := uuid.New().String()
	sessionID := uuid.New().String()

	loginUser := &authmodel.LoginUser{
		UserId:      user.Id,
		Username:    user.Username,
		Nickname:    user.Nickname,
		Avatar:      user.Avatar,
		Depts:       depts,
		Roles:       roles,
		Permissions: permissions,
		ClientId:    clientId,
	}

	if err := s.persistSessionTokens(ctx, sessionID, accessToken, refreshToken, loginUser, accessTTL, refreshTTL); err != nil {
		return nil, err
	}
	// 记录用户全局会话，供 SLO 使用
	_ = s.store.AddUserSession(ctx, userId, clientId, sessionID)

	return &authmodel.TokenPair{
		AccessToken:      accessToken,
		RefreshToken:     refreshToken,
		ExpiresIn:        accessTTL,
		RefreshExpiresIn: refreshTTL,
		TokenType:        "Bearer",
		Scope:            strings.Join(scopes, " "),
	}, nil
}

// Authorize 校验客户端参数，生成 state 存入 Redis，返回授权页展示信息。
func (s *authService) Authorize(ctx context.Context, clientId, redirectURI string) (*AuthorizeInfo, error) {
	client, err := s.authRepo.GetByClientId(clientId)
	if err != nil {
		return nil, &AuthError{Code: consts.ErrOAuthInvalidClient, HTTPStatus: 401}
	}

	// 校验 redirect_uri（若客户端配置了则必须匹配）
	if client.WebServerRedirectUri != "" && redirectURI != "" && redirectURI != client.WebServerRedirectUri {
		return nil, &AuthError{Code: consts.ErrOAuthRedirectURIMismatch, HTTPStatus: 400}
	}
	effectiveRedirectURI := redirectURI
	if effectiveRedirectURI == "" {
		effectiveRedirectURI = client.WebServerRedirectUri
	}

	// 生成随机 state，将 clientId|redirectURI 存入 Redis（TTL 10 分钟）
	state := uuid.New().String()
	stateVal := clientId + "|" + effectiveRedirectURI
	if err := s.store.SaveOAuthState(ctx, state, stateVal, 10*time.Minute); err != nil {
		return nil, err
	}

	return &AuthorizeInfo{
		ClientId:    client.ClientId,
		ClientName:  client.ClientName,
		LogoURI:     client.LogoUri,
		Scope:       client.Scope,
		RedirectURI: effectiveRedirectURI,
		State:       state,
		AutoApprove: client.Autoapprove == 1,
	}, nil
}

// ApproveAuthorize 用户确认授权：校验 state，生成授权码，返回携带 code 的回调 URL。
func (s *authService) ApproveAuthorize(ctx context.Context, userId int64, state string, scopes []string) (string, error) {
	stateVal, err := s.store.GetOAuthState(ctx, state)
	if err != nil {
		return "", &AuthError{Code: consts.ErrOAuthInvalidState, HTTPStatus: 400}
	}
	_ = s.store.DeleteOAuthState(ctx, state)

	// 解析 clientId 和 redirectURI
	parts := strings.SplitN(stateVal, "|", 2)
	clientId := parts[0]
	redirectURI := ""
	if len(parts) == 2 {
		redirectURI = parts[1]
	}
	if redirectURI == "" {
		client, err := s.authRepo.GetByClientId(clientId)
		if err != nil {
			return "", &AuthError{Code: consts.ErrOAuthInvalidClient, HTTPStatus: 401}
		}
		redirectURI = client.WebServerRedirectUri
	}
	if redirectURI == "" {
		return "", &AuthError{Code: consts.ErrOAuthRedirectURIMismatch, HTTPStatus: 400}
	}

	// 生成授权码，存入 scopes
	code := uuid.New().String()
	if err := s.store.SaveAuthCode(ctx, code, userId, clientId, scopes, 10*time.Minute); err != nil {
		return "", err
	}

	// 拼接回调 URL
	sep := "?"
	if strings.Contains(redirectURI, "?") {
		sep = "&"
	}
	return redirectURI + sep + "code=" + code + "&state=" + state, nil
}

// Callback 处理授权回调：校验 state，使用 oauth2.Config.Exchange 换取令牌。
// Requirements: 1.2, 1.4, 6.1, 6.2
func (s *authService) Callback(ctx context.Context, code, state string) (*authmodel.TokenPair, error) {
	// 校验 state 参数（Requirements 1.2, 1.4）
	clientId, err := s.store.GetOAuthState(ctx, state)
	if err != nil {
		return nil, &AuthError{Code: consts.ErrOAuthInvalidState, HTTPStatus: 400}
	}
	_ = s.store.DeleteOAuthState(ctx, state)

	// 读取授权码信息以获取 userId 和 scopes
	userId, _, scopes, err := s.store.GetAuthCodeInfo(ctx, code)
	if err != nil {
		return nil, &AuthError{Code: consts.ErrInvalidAuthCode, HTTPStatus: 400}
	}
	_ = s.store.DeleteAuthCode(ctx, code)

	client, err := s.authRepo.GetByClientId(clientId)
	if err != nil {
		return nil, &AuthError{Code: consts.ErrOAuthInvalidClient, HTTPStatus: 401}
	}

	user, err := s.authRepo.GetUserByUserId(userId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, &AuthError{Code: consts.ErrUserNotFound, HTTPStatus: 404}
		}
		return nil, err
	}
	if user.Status != 0 {
		return nil, &AuthError{Code: consts.ErrAccountDisabled, HTTPStatus: 403}
	}

	accessTTL, refreshTTL := client.GetTTL()
	return s.issueTokenPairWithScopes(ctx, user.Id, clientId, scopes, accessTTL, refreshTTL)
}

// IssueShortToken 为指定用户颁发短期 access token（TTL 10 分钟，授权页登录用）。
func (s *authService) IssueShortToken(ctx context.Context, userId int64) (string, error) {
	user, err := s.authRepo.GetUserByUserId(userId)
	if err != nil {
		return "", err
	}
	token := uuid.New().String()
	sessionID := uuid.New().String()
	loginUser := &authmodel.LoginUser{
		UserId:   user.Id,
		Username: user.Username,
		Nickname: user.Nickname,
		Roles:    []authmodel.Role{},
		Depts:    []authmodel.Dept{},
		ClientId: "__authorize__",
	}
	if err := s.store.SaveSession(ctx, sessionID, loginUser, 10*time.Minute); err != nil {
		return "", err
	}
	if err := s.store.SaveAccessToken(ctx, token, sessionID, 10*time.Minute); err != nil {
		return "", err
	}
	return token, nil
}

// GetLoginUser 根据 access token 获取登录用户信息。
func (s *authService) GetLoginUser(ctx context.Context, accessToken string) (*authmodel.LoginUser, error) {
	return s.store.GetLoginUser(ctx, accessToken)
}

// ValidateUserPassword 校验用户名密码，返回 userId。
func (s *authService) ValidateUserPassword(ctx context.Context, username, password, keyId string) (int64, error) {
	user, err := s.authRepo.GetUserByUsername(username)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return 0, &AuthError{Code: consts.ErrInvalidCredentials, HTTPStatus: 401}
		}
		return 0, err
	}

	actualPassword := password
	if keyId != "" {
		privateKey, err := s.store.GetRSAPrivateKey(ctx, keyId)
		if err == nil && privateKey != "" {
			_ = s.store.DeleteRSAPrivateKey(ctx, keyId)
			if decrypted, err := utils.RSADecrypt(password, privateKey); err == nil {
				actualPassword = decrypted
			}
		}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(actualPassword)); err != nil {
		return 0, &AuthError{Code: consts.ErrInvalidCredentials, HTTPStatus: 401}
	}
	if user.Status != 0 {
		return 0, &AuthError{Code: consts.ErrAccountDisabled, HTTPStatus: 403}
	}
	return user.Id, nil
}

func (s *authService) GetRouters(ctx context.Context, langCode string) ([]authmodel.Menu, error) {
	user, err := authctx.GetCurrentUser(ctx)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("当前用户不存在")
	}
	var menus []dto.MenuDTO
	if user.UserId == consts.SuperAdminId {
		menus, err = s.authRepo.GetAllMenus(langCode)
		if err != nil {
			return nil, err
		}
	} else {
		roleIds, err := authctx.GetCurrentRoleIds(ctx)
		if err != nil {
			return nil, err
		}
		// 这里也必须捕获错误！
		menus, err = s.authRepo.GetMenusByRoleIds(roleIds, langCode)
		if err != nil {
			return nil, err
		}
	}
	return buildMenuTree(menus, 0), nil
}

func buildMenuTree(menus []dto.MenuDTO, parentId int64) []authmodel.Menu {
	var tree []authmodel.Menu
	for _, menu := range menus {
		if menu.ParentId == parentId {
			var perms []string
			if menu.Permission != "" {
				perms = []string{menu.Permission}
			}
			node := authmodel.Menu{
				Key:       fmt.Sprintf("%d", menu.Id),
				Path:      menu.Path,
				Component: menu.Component,
				Query:     menu.Query,
				Meta: &authmodel.Meta{
					Title:       menu.Title,
					Icon:        menu.Icon,
					Hidden:      menu.Visible == 0, // 直接赋值，修复逻辑
					IsFrame:     menu.IsFrame == 1,
					Permissions: perms,
				},
			}
			if menu.ActiveId != 0 {
				node.Meta.ActiveMenu = fmt.Sprintf("%d", menu.ActiveId)
			}
			// 空权限清空
			if menu.Permission == "" {
				node.Meta.Permissions = nil
			}

			children := buildMenuTree(menus, menu.Id)
			if len(children) > 0 {
				node.Children = children
			}
			tree = append(tree, node)
		}
	}
	return tree
}

// AuthError 表示带有 HTTP 状态码的认证领域错误。
type AuthError struct {
	Code       string
	HTTPStatus int
}

func (e *AuthError) Error() string {
	return e.Code
}

// GetPublicKey 生成并返回临时的 RSA 公钥和关联凭证。
func (s *authService) GetPublicKey(ctx context.Context) (*authmodel.PublicKeyResponse, error) {
	priv, pub, err := utils.GenerateRSAKey(2048)
	if err != nil {
		return nil, err
	}

	keyId := uuid.New().String()
	// 私钥存 2 分钟足以支撑登录流程
	if err := s.store.SaveRSAPrivateKey(ctx, keyId, priv, 2*time.Minute); err != nil {
		return nil, err
	}

	return &authmodel.PublicKeyResponse{
		KeyId:     keyId,
		PublicKey: pub,
	}, nil
}

// GetProfile 获取当前登录用户的详尽档案信息。
func (s *authService) GetProfile(id int64) (*dto.UserProfile, error) {
	user, err := s.authRepo.GetUserByUserId(id)
	if err != nil {
		return nil, err
	}

	roles, _ := s.authRepo.GetRolesByUserId(id)
	deptPosts, _ := s.authRepo.GetProfileDeptPosts(id)

	var profileRoles []dto.ProfileRole
	for _, role := range roles {
		profileRoles = append(profileRoles, dto.ProfileRole{
			Id:       role.Id,
			RoleKey:  role.RoleKey,
			RoleName: role.RoleName,
		})
	}

	roleIds := make([]int64, len(roles))
	for i, role := range roles {
		roleIds[i] = role.Id
	}
	permissions, _ := s.authRepo.GetPermissionsByRoleIds(roleIds)

	return &dto.UserProfile{
		User:        *user,
		Roles:       profileRoles,
		DeptPosts:   deptPosts,
		Permissions: permissions,
	}, nil
}

// LoginByUserId 直接为指定用户颁发令牌对（Passkey 等无密码登录使用）。
// TTL 优先从 oauth_client 表读取；若客户端不存在或未配置，回退到 config.yaml 中的兜底值。
func (s *authService) LoginByUserId(ctx context.Context, userId int64, clientId string) (*authmodel.TokenPair, error) {
	var accessTTL, refreshTTL int64
	if client, err := s.authRepo.GetByClientId(clientId); err == nil {
		accessTTL, refreshTTL = client.GetTTL()
	}
	if accessTTL <= 0 {
		accessTTL = s.cfg.AccessTokenTTL
	}
	if refreshTTL <= 0 {
		refreshTTL = s.cfg.RefreshTokenTTL
	}
	return s.issueTokenPair(ctx, userId, clientId, accessTTL, refreshTTL)
}

// RecordLoginLog 记录登录日志（供外部服务调用）。
func (s *authService) RecordLoginLog(ctx context.Context, username, clientId string, grantType string, status int16, msg, userAgent, ip string) {
	s.recordLoginLog(ctx, username, clientId, grantType, status, msg, userAgent, ip)
}

func (s *authService) GetLatestLoginLog(ctx context.Context) ([]*log_grpc.LoginEntry, error) {
	user, err := authctx.GetCurrentUser(ctx)
	if err != nil {
		return nil, err
	}
	entries, err := s.logStore.GetLatestLoginLog(ctx, user.Username, 10)
	if err != nil {
		return nil, err
	}

	result := make([]*log_grpc.LoginEntry, 0, len(entries))
	for i := range entries {
		entry := entries[i]
		result = append(result, &entry)
	}
	return result, nil
}
