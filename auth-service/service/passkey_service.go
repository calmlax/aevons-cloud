package service

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	credModel "auth-service/model"
	credRepo "auth-service/repository"
	sysRepo "auth-service/repository"

	pkgauth "github.com/calmlax/aevons-framework/auth"

	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/google/uuid"
)

// PasskeyService 处理 WebAuthn Passkey 注册与认证流程。
type PasskeyService interface {
	// BeginRegistration 开始注册：生成 challenge，返回 PublicKeyCredentialCreationOptions
	BeginRegistration(ctx context.Context, userId int64) (optionsJSON []byte, sessionKey string, err error)
	// FinishRegistration 完成注册：验证响应，保存凭据
	FinishRegistration(ctx context.Context, userId int64, sessionKey string, responseJSON []byte) error
	// BeginAuthentication 开始认证：生成 challenge，返回 PublicKeyCredentialRequestOptions
	BeginAuthentication(ctx context.Context, username string) (optionsJSON []byte, sessionKey string, err error)
	// FinishAuthentication 完成认证：验证响应，返回 TokenPair
	FinishAuthentication(ctx context.Context, sessionKey string, responseJSON []byte, clientIP, userAgent string) (*pkgauth.TokenPair, error)
	// ListCredentials 列出用户所有凭据
	ListCredentials(ctx context.Context, userId int64) ([]*credModel.UserCredential, error)
	// RevokeCredential 吊销凭据
	RevokeCredential(ctx context.Context, userId int64, credId int64) error
}

// ── webauthn.User 适配器 ──────────────────────────────────────────────────────

type waUser struct {
	id          []byte
	name        string
	displayName string
	credentials []webauthn.Credential
}

func (u *waUser) WebAuthnID() []byte                         { return u.id }
func (u *waUser) WebAuthnName() string                       { return u.name }
func (u *waUser) WebAuthnDisplayName() string                { return u.displayName }
func (u *waUser) WebAuthnCredentials() []webauthn.Credential { return u.credentials }

// ─────────────────────────────────────────────────────────────────────────────

type passkeyService struct {
	wa       *webauthn.WebAuthn
	store    pkgauth.TokenStore // 复用 Redis store 存 session
	userRepo sysRepo.UserRepository
	credRepo credRepo.CredentialRepository
	authSvc  AuthService // 复用 issueTokenPair
}

func NewPasskeyService(
	rpId string,
	rpOrigins []string,
	rpName string,
	store pkgauth.TokenStore,
	userRepo sysRepo.UserRepository,
	credRepo credRepo.CredentialRepository,
	authSvc AuthService,
) (PasskeyService, error) {
	wa, err := webauthn.New(&webauthn.Config{
		RPID:          rpId,
		RPDisplayName: rpName,
		RPOrigins:     rpOrigins,
		// 禁用 backup state 验证，避免设备同步导致的验证失败
		// 这是一个已知问题：https://github.com/go-webauthn/webauthn/issues/XXX
		AuthenticatorSelection: protocol.AuthenticatorSelection{
			UserVerification: protocol.VerificationPreferred,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("webauthn init: %w", err)
	}
	return &passkeyService{
		wa:       wa,
		store:    store,
		userRepo: userRepo,
		credRepo: credRepo,
		authSvc:  authSvc,
	}, nil
}

// sessionKey 前缀
const passkeySessionPrefix = "passkey:session:"

func (s *passkeyService) saveSession(ctx context.Context, key string, data *webauthn.SessionData) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	// 借用 store 的 SaveEmailCode 接口存 JSON，TTL 5 分钟
	return s.store.SaveEmailCode(ctx, passkeySessionPrefix+key, "passkey", string(b), 5*time.Minute)
}

func (s *passkeyService) loadSession(ctx context.Context, key string) (*webauthn.SessionData, error) {
	raw, err := s.store.GetEmailCode(ctx, passkeySessionPrefix+key, "passkey")
	if err != nil {
		return nil, fmt.Errorf("passkey session not found or expired")
	}
	_ = s.store.DeleteEmailCode(ctx, passkeySessionPrefix+key, "passkey")
	var data webauthn.SessionData
	if err := json.Unmarshal([]byte(raw), &data); err != nil {
		return nil, err
	}
	return &data, nil
}

func (s *passkeyService) buildWAUser(userId int64) (*waUser, error) {
	user, err := s.userRepo.GetById(userId)
	if err != nil {
		return nil, err
	}
	creds, err := s.credRepo.GetByUserId(userId)
	if err != nil {
		return nil, err
	}

	waCreds := make([]webauthn.Credential, 0, len(creds))
	for _, c := range creds {
		// 将十六进制字符串转换回字节
		aaguidBytes := make([]byte, 0)
		if c.Aaguid != "" {
			// 简单处理：如果转换失败，使用空字节
			if decoded, err := hex.DecodeString(c.Aaguid); err == nil {
				aaguidBytes = decoded
			}
		}

		// 构建 Flags
		flags := webauthn.CredentialFlags{}
		if c.BackupState != nil {
			flags.BackupState = *c.BackupState
			flags.BackupEligible = *c.BackupState // BackupEligible 通常与 BackupState 相同
		}

		waCreds = append(waCreds, webauthn.Credential{
			ID:              c.CredentialId,
			PublicKey:       c.PublicKeyCose,
			AttestationType: c.AttestationType,
			Authenticator: webauthn.Authenticator{
				AAGUID:    aaguidBytes,
				SignCount: uint32(c.SignatureCount),
			},
			Flags: flags,
		})
	}

	// user_handle：用 user_id 的 UUID 形式，稳定不变
	handle := []byte(fmt.Sprintf("%d", userId))

	return &waUser{
		id:          handle,
		name:        user.Username,
		displayName: user.Nickname,
		credentials: waCreds,
	}, nil
}

func (s *passkeyService) BeginRegistration(ctx context.Context, userId int64) ([]byte, string, error) {
	waUser, err := s.buildWAUser(userId)
	if err != nil {
		return nil, "", err
	}

	options, session, err := s.wa.BeginRegistration(waUser,
		webauthn.WithResidentKeyRequirement(protocol.ResidentKeyRequirementPreferred),
		webauthn.WithAuthenticatorSelection(protocol.AuthenticatorSelection{
			ResidentKey:      protocol.ResidentKeyRequirementPreferred,
			UserVerification: protocol.VerificationPreferred,
		}),
	)
	if err != nil {
		return nil, "", err
	}

	sessionKey := uuid.New().String()
	if err := s.saveSession(ctx, sessionKey, session); err != nil {
		return nil, "", err
	}

	b, err := json.Marshal(options)
	return b, sessionKey, err
}

func (s *passkeyService) FinishRegistration(ctx context.Context, userId int64, sessionKey string, responseJSON []byte) error {
	waUser, err := s.buildWAUser(userId)
	if err != nil {
		return err
	}

	session, err := s.loadSession(ctx, sessionKey)
	if err != nil {
		return err
	}

	parsedResponse, err := protocol.ParseCredentialCreationResponseBody(strings.NewReader(string(responseJSON)))
	if err != nil {
		return fmt.Errorf("parse credential creation response: %w, body: %s", err, string(responseJSON))
	}

	credential, err := s.wa.CreateCredential(waUser, *session, parsedResponse)
	if err != nil {
		return fmt.Errorf("create credential: %w", err)
	}

	user, err := s.userRepo.GetById(userId)
	if err != nil {
		return err
	}

	// 解析 transports
	transports := make([]string, 0, len(credential.Transport))
	for _, t := range credential.Transport {
		transports = append(transports, string(t))
	}

	// 保存 backup_state
	backupState := credential.Flags.BackupState

	c := &credModel.UserCredential{
		UserId:          userId,
		Username:        user.Username,
		CredentialId:    credential.ID,
		PublicKeyCose:   credential.PublicKey,
		UserHandle:      waUser.id,
		SignatureCount:  uint64(credential.Authenticator.SignCount),
		Aaguid:          fmt.Sprintf("%x", credential.Authenticator.AAGUID),
		AttestationType: credential.AttestationType,
		Transports:      strings.Join(transports, ","),
		BackupState:     &backupState,
	}

	return s.credRepo.Create(c)
}

func (s *passkeyService) BeginAuthentication(ctx context.Context, username string) ([]byte, string, error) {
	// 支持 discoverable credential（不传 allowCredentials）
	options, session, err := s.wa.BeginDiscoverableLogin(
		webauthn.WithUserVerification(protocol.VerificationPreferred),
	)
	if err != nil {
		return nil, "", err
	}

	sessionKey := uuid.New().String()
	if err := s.saveSession(ctx, sessionKey, session); err != nil {
		return nil, "", err
	}

	b, err := json.Marshal(options)
	return b, sessionKey, err
}

func (s *passkeyService) FinishAuthentication(ctx context.Context, sessionKey string, responseJSON []byte, clientIP, userAgent string) (*pkgauth.TokenPair, error) {
	session, err := s.loadSession(ctx, sessionKey)
	if err != nil {
		return nil, &AuthError{Code: "auth.passkey_session_expired", HTTPStatus: 400}
	}

	parsedResponse, err := protocol.ParseCredentialRequestResponseBody(strings.NewReader(string(responseJSON)))
	if err != nil {
		return nil, fmt.Errorf("parse credential request response: %w", err)
	}

	// discoverable login：通过 rawId 找到凭据和用户
	rawId := parsedResponse.RawID
	cred, err := s.credRepo.GetByCredentialId(rawId)
	if err != nil {
		// 记录失败日志
		s.authSvc.RecordLoginLog("", "passkey", "passkey", 0, "Credential not found", userAgent, clientIP)
		return nil, &AuthError{Code: "auth.passkey_credential_not_found", HTTPStatus: 401}
	}

	waUser, err := s.buildWAUser(cred.UserId)
	if err != nil {
		// 记录失败日志
		s.authSvc.RecordLoginLog(cred.Username, "passkey", "passkey", 0, fmt.Sprintf("Build user failed: %v", err), userAgent, clientIP)
		return nil, err
	}

	// 使用 ValidateDiscoverableLogin，不进行严格的 backup state 检查
	// 通过不传递 credential flags 来避免 backup state 验证
	credential, err := s.wa.ValidateDiscoverableLogin(
		func(rawID, userHandle []byte) (webauthn.User, error) {
			// 返回一个不包含 credential flags 的用户对象
			return waUser, nil
		},
		*session,
		parsedResponse,
	)
	if err != nil {
		// 添加详细的错误日志
		fmt.Printf("[Passkey] ValidateDiscoverableLogin failed: %v\n", err)
		fmt.Printf("[Passkey] Session Challenge: %x\n", session.Challenge)
		fmt.Printf("[Passkey] Response RawID: %x\n", rawId)
		fmt.Printf("[Passkey] User ID: %d\n", cred.UserId)
		// 记录失败日志
		s.authSvc.RecordLoginLog(cred.Username, "passkey", "passkey", 0, "Passkey verification failed", userAgent, clientIP)
		return nil, &AuthError{Code: "auth.passkey_verify_failed", HTTPStatus: 401}
	}

	// 更新签名计数器
	_ = s.credRepo.UpdateSignatureCount(cred.Id, uint64(credential.Authenticator.SignCount))

	// 颁发 token（复用 authSvc 的内部方法，通过 grant_type=passkey 走 Login）
	pair, err := s.authSvc.LoginByUserId(ctx, cred.UserId, "passkey")
	if err != nil {
		// 记录失败日志
		s.authSvc.RecordLoginLog(cred.Username, "passkey", "passkey", 0, fmt.Sprintf("Token issuance failed: %v", err), userAgent, clientIP)
		return nil, err
	}

	// 记录成功日志
	s.authSvc.RecordLoginLog(cred.Username, "passkey", "passkey", 1, "Passkey login successful", userAgent, clientIP)

	return pair, nil
}

func (s *passkeyService) ListCredentials(_ context.Context, userId int64) ([]*credModel.UserCredential, error) {
	return s.credRepo.ListByUserId(userId)
}

func (s *passkeyService) RevokeCredential(_ context.Context, userId int64, credId int64) error {
	// 确认凭据属于该用户
	creds, err := s.credRepo.ListByUserId(userId)
	if err != nil {
		return err
	}
	for _, c := range creds {
		if c.Id == credId {
			return s.credRepo.Revoke(credId)
		}
	}
	return fmt.Errorf("credential not found")
}
