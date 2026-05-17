/**
 * 终端应用 Service
 *
 * @author
 * @date 2026-04-09 01:26:40.390618977 +0000 UTC
 * @website https://www.aevons.com
 * @email 1073602@qq.com
 * @copyright 2025-2026 Aevons
 */
package service

import (
	"context"
	"strings"
	"system-service/internal/dto"
	"system-service/internal/model"
	"system-service/internal/repository"

	"github.com/calmlax/aevons-framework/consts"

	"github.com/calmlax/aevons-framework/core/base"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// 继承BaseService
type OauthClientService interface {
	base.BaseService[model.OauthClient, *dto.OauthClientQuery]
	// ValidateClient 验证终端应用
	ValidateClient(ctx context.Context, clientId, clientSecret, grantType string) (*model.OauthClient, error)
	// GetByClientId 根据 clientId 查询客户端配置（不校验 secret 和 grant_type）。
	GetByClientId(ctx context.Context, clientId string) (*model.OauthClient, error)
}

type oauthClientService struct {
	base.BaseService[model.OauthClient, *dto.OauthClientQuery]
	repo repository.OauthClientRepository
}

func NewOauthClientService(repo repository.OauthClientRepository) OauthClientService {
	baseSrv := base.NewBaseService[model.OauthClient, *dto.OauthClientQuery](repo)
	return &oauthClientService{
		BaseService: baseSrv,
		repo:        repo,
	}
}

func (s *oauthClientService) GetByClientId(_ context.Context, clientId string) (*model.OauthClient, error) {
	client, err := s.repo.GetByField("client_id", clientId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, &AuthError{Code: consts.ErrOAuthInvalidClient, HTTPStatus: 401}
		}
		return nil, err
	}
	return client, nil
}

func (s *oauthClientService) ValidateClient(_ context.Context, clientId, clientSecret, grantType string) (*model.OauthClient, error) {
	client, err := s.repo.GetByField("client_id", clientId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, &AuthError{Code: consts.ErrOAuthInvalidClient, HTTPStatus: 401}
		}
		return nil, err
	}
	if clientSecret != "" {
		if err := bcrypt.CompareHashAndPassword([]byte(client.ClientSecret), []byte(clientSecret)); err != nil {
			return nil, &AuthError{Code: consts.ErrOAuthInvalidClient, HTTPStatus: 401}
		}
	}
	if !s.isGrantTypeAllowed(client.AuthorizedGrantTypes, grantType) {
		return nil, &AuthError{Code: consts.ErrOAuthUnsupportedGrant, HTTPStatus: 400}
	}
	return client, nil
}

func (s *oauthClientService) isGrantTypeAllowed(authorizedGrantTypes, grantType string) bool {
	if authorizedGrantTypes == "" {
		return false
	}
	for _, gt := range strings.Split(authorizedGrantTypes, ",") {
		if strings.TrimSpace(gt) == grantType {
			return true
		}
	}
	return false
}

// AuthError 表示带有 HTTP 状态码的认证领域错误。
type AuthError struct {
	Code       string
	HTTPStatus int
}

func (e *AuthError) Error() string {
	return e.Code
}
