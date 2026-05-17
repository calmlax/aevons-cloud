/**
 * 终端应用 Model
 *
 * @author
 * @date 2026-04-09 01:26:40.390618977 +0000 UTC
 * @website https://www.aevons.com
 * @email 1073602@qq.com
 * @copyright 2025-2026 Aevons
 */
package model

import (
	"github.com/calmlax/aevons-framework/core/base"
)

type OauthClient struct {

	//ID
	Id int64 `gorm:"column:id;primaryKey;autoIncrement;type:bigint;comment:ID" json:"id,string"`

	//客户端ID
	ClientId string `gorm:"column:client_id;type:varchar(32);comment:客户端ID" json:"clientId"`

	//客户端秘钥
	ClientSecret string `gorm:"column:client_secret;type:varchar(256);comment:客户端秘钥" json:"-"`

	//客户端名称
	ClientName string `gorm:"column:client_name;type:varchar(255);comment:客户端名称" json:"clientName"`

	//客户端LOGO
	LogoUri string `gorm:"column:logo_uri;type:varchar(255);comment:客户端LOGO" json:"logoUri"`

	//授权范围
	Scope string `gorm:"column:scope;type:varchar(256);comment:授权范围" json:"scope"`

	//授权类型
	AuthorizedGrantTypes string `gorm:"column:authorized_grant_types;type:varchar(256);comment:授权类型" json:"authorizedGrantTypes"`

	//回调地址
	WebServerRedirectUri string `gorm:"column:web_server_redirect_uri;type:varchar(256);comment:回调地址" json:"webServerRedirectUri"`

	//访问令牌有效期（秒）
	AccessTokenValidity int `gorm:"column:access_token_validity;type:int;comment:访问令牌有效期（秒）" json:"accessTokenValidity"`

	//刷新令牌有效期（秒）
	RefreshTokenValidity int `gorm:"column:refresh_token_validity;type:int;comment:刷新令牌有效期（秒）" json:"refreshTokenValidity"`

	//自动授权（0否，1是）
	Autoapprove int16 `gorm:"column:autoapprove;type:tinyint(1);comment:自动授权（0否，1是）" json:"autoapprove"`

	base.DefaultModel
}

// TableName 指定表名
func (OauthClient) TableName() string {
	return "sys_oauth_client"
}

// GetTTL 返回客户端配置的令牌有效期（秒）。
// 若未配置（值为 0），调用方应回退到系统默认值。
func (c *OauthClient) GetTTL() (accessTTL int64, refreshTTL int64) {
	return int64(c.AccessTokenValidity), int64(c.RefreshTokenValidity)
}
