package repository

import (
	"auth-service/model"
	"time"

	"gorm.io/gorm"
)

type CredentialRepository interface {
	// GetByUserId 获取用户所有未吊销的凭据
	GetByUserId(userId int64) ([]*model.UserCredential, error)
	// GetByCredentialId 按凭据 ID 查找
	GetByCredentialId(credentialId []byte) (*model.UserCredential, error)
	// Create 保存新凭据
	Create(c *model.UserCredential) error
	// UpdateSignatureCount 更新签名计数器和最后使用时间
	UpdateSignatureCount(id int64, count uint64) error
	// Revoke 吊销凭据
	Revoke(id int64) error
	// ListByUserId 列出用户所有凭据（含已吊销）
	ListByUserId(userId int64) ([]*model.UserCredential, error)
}

type credentialRepository struct{ db *gorm.DB }

func NewCredentialRepository(db *gorm.DB) CredentialRepository {
	return &credentialRepository{db: db}
}

func (r *credentialRepository) GetByUserId(userId int64) ([]*model.UserCredential, error) {
	var list []*model.UserCredential
	err := r.db.Where("user_id = ? AND is_revoked = 0", userId).Find(&list).Error
	return list, err
}

func (r *credentialRepository) GetByCredentialId(credentialId []byte) (*model.UserCredential, error) {
	var c model.UserCredential
	err := r.db.Where("credential_id = ? AND is_revoked = 0", credentialId).First(&c).Error
	return &c, err
}

func (r *credentialRepository) Create(c *model.UserCredential) error {
	return r.db.Create(c).Error
}

func (r *credentialRepository) UpdateSignatureCount(id int64, count uint64) error {
	now := time.Now()
	return r.db.Model(&model.UserCredential{}).Where("id = ?", id).
		Updates(map[string]any{"signature_count": count, "last_used_at": now}).Error
}

func (r *credentialRepository) Revoke(id int64) error {
	return r.db.Model(&model.UserCredential{}).Where("id = ?", id).
		Update("is_revoked", true).Error
}

func (r *credentialRepository) ListByUserId(userId int64) ([]*model.UserCredential, error) {
	var list []*model.UserCredential
	err := r.db.Where("user_id = ?", userId).Order("created_at DESC").Find(&list).Error
	return list, err
}
