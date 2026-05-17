package model

import "time"

// UserCredential 用户 Passkey 凭据（对应 sys_user_credential 表）
type UserCredential struct {
	Id              int64      `gorm:"column:id;primaryKey;autoIncrement"                json:"id,string"`
	UserId          int64      `gorm:"column:user_id;not null;index:idx_user_id"         json:"user_id,string"`
	Username        string     `gorm:"column:username;type:varchar(64);not null"         json:"username"`
	CredentialId    []byte     `gorm:"column:credential_id;type:varbinary(255);not null;uniqueIndex:uk_credential_id" json:"-"`
	PublicKeyCose   []byte     `gorm:"column:public_key_cose;type:varbinary(512);not null" json:"-"`
	UserHandle      []byte     `gorm:"column:user_handle;type:varbinary(64);not null"    json:"-"`
	SignatureCount  uint64     `gorm:"column:signature_count;not null;default:0"         json:"signature_count"`
	Aaguid          string     `gorm:"column:aaguid;type:char(36)"                       json:"aaguid"`
	AttestationType string     `gorm:"column:attestation_type;type:varchar(32)"          json:"attestation_type"`
	Attachment      string     `gorm:"column:attachment;type:varchar(32)"                json:"attachment"`
	Transports      string     `gorm:"column:transports;type:varchar(255)"               json:"transports"`
	DeviceType      string     `gorm:"column:device_type;type:varchar(32)"               json:"device_type"`
	BackupState     *bool      `gorm:"column:backup_state;type:tinyint(1)"               json:"backup_state"`
	DeviceName      string     `gorm:"column:device_name;type:varchar(255)"              json:"device_name"`
	IsRevoked       bool       `gorm:"column:is_revoked;not null;default:0"              json:"is_revoked"`
	LastUsedAt      *time.Time `gorm:"column:last_used_at"                               json:"last_used_at"`
	CreatedAt       time.Time  `gorm:"column:created_at;autoCreateTime"                  json:"created_at"`
	UpdatedAt       *time.Time `gorm:"column:updated_at;autoUpdateTime"                  json:"updated_at"`
}

func (UserCredential) TableName() string { return "sys_user_credential" }
