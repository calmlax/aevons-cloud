package repository

import (
	"auth-service/model"
	"time"

	"gorm.io/gorm"
)

type CredentialRepository interface {
	// GetCredentialByUserId 获取用户所有未吊销的凭据
	GetCredentialByUserId(userId int64) ([]*model.UserCredential, error)
	// GetCredentialByCredentialId 按凭据 ID 查找
	GetCredentialByCredentialId(credentialId []byte) (*model.UserCredential, error)
	// CreateCredential 保存新凭据
	CreateCredential(c *model.UserCredential) error
	// UpdateCredentialSignatureCount 更新签名计数器和最后使用时间
	UpdateCredentialSignatureCount(id int64, count uint64) error
	// RevokeCredential 吊销凭据
	RevokeCredential(id int64) error
	// ListCredentialByUserId 列出用户所有凭据（含已吊销）
	ListCredentialByUserId(userId int64) ([]*model.UserCredential, error)
}

type credentialRepository struct{ db *gorm.DB }

func NewCredentialRepository(db *gorm.DB) CredentialRepository {
	return &credentialRepository{db: db}
}

func (r *credentialRepository) GetCredentialByUserId(userId int64) ([]*model.UserCredential, error) {
	var list []*model.UserCredential
	err := r.db.Where("user_id = ? AND is_revoked = 0", userId).Find(&list).Error
	return list, err
}

func (r *credentialRepository) GetCredentialByCredentialId(credentialId []byte) (*model.UserCredential, error) {
	var c model.UserCredential
	err := r.db.Where("credential_id = ? AND is_revoked = 0", credentialId).First(&c).Error
	return &c, err
}

func (r *credentialRepository) CreateCredential(c *model.UserCredential) error {
	return r.db.Create(c).Error
}

func (r *credentialRepository) UpdateCredentialSignatureCount(id int64, count uint64) error {
	now := time.Now()
	return r.db.Model(&model.UserCredential{}).Where("id = ?", id).
		Updates(map[string]any{"signature_count": count, "last_used_at": now}).Error
}

func (r *credentialRepository) RevokeCredential(id int64) error {
	return r.db.Model(&model.UserCredential{}).Where("id = ?", id).
		Update("is_revoked", true).Error
}

func (r *credentialRepository) ListCredentialByUserId(userId int64) ([]*model.UserCredential, error) {
	var list []*model.UserCredential
	err := r.db.Where("user_id = ?", userId).Order("created_at DESC").Find(&list).Error
	return list, err
}
