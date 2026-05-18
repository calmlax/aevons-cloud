/**
 * 终端应用 DTO
 *
 * @author
 * @date 2026-04-09 01:26:40.390618977 +0000 UTC
 * @website https://www.aevons.com
 * @email 1073602@qq.com
 * @copyright 2025-2026 Aevons
 */
package dto

import (
	"github.com/calmlax/aevons-framework/core/base"
)

// q 支持查询操作: eq, ne, gt, gte, lt, lte, like, like_l, like_r, in, not_in, between, is_null, not_null
// gorm.column 必须对应真实数据库列名
type OauthClientQuery struct {
	// 客户端ID
	ClientId *string `form:"clientId" gorm:"column:client_id" q:"eq"`
	// 客户端秘钥
	ClientSecret *string `form:"clientSecret" gorm:"column:client_secret" q:"eq"`
	// 客户端名称
	ClientName *string `form:"clientName" gorm:"column:client_name" q:"like"`
	// 客户端LOGO
	LogoUri *string `form:"logoUri" gorm:"column:logo_uri" q:"eq"`
	// 授权范围
	Scope *string `form:"scope" gorm:"column:scope" q:"eq"`
	// 授权类型
	AuthorizedGrantTypes *[]string `form:"authorizedGrantTypes" gorm:"column:authorized_grant_types" q:"in"`
	// 回调地址
	WebServerRedirectUri *string `form:"webServerRedirectUri" gorm:"column:web_server_redirect_uri" q:"eq"`
	// 自动授权（0否，1是）
	Autoapprove *int16 `form:"autoapprove" gorm:"column:autoapprove" q:"eq"`
	// 资源权限，按服务名称逗号分割
	Resources *string `form:"resources" gorm:"column:resources" q:"like"`
	base.BaseQuery
}

type CreateOauthClientDTO struct {

	//客户端ID
	ClientId string `json:"clientId" binding:"required,max=32"`
	//客户端秘钥
	ClientSecret string `json:"clientSecret" binding:"max=256"`
	//客户端名称
	ClientName string `json:"clientName" binding:"max=255"`
	//客户端LOGO
	LogoUri string `json:"logoUri" binding:"max=255"`
	//授权范围
	Scope string `json:"scope" binding:"max=256"`
	//授权类型
	AuthorizedGrantTypes string `json:"authorizedGrantTypes" binding:"max=256"`
	//回调地址
	WebServerRedirectUri string `json:"webServerRedirectUri" binding:"max=256"`
	//访问令牌有效期（秒）
	AccessTokenValidity int `json:"accessTokenValidity"`
	//刷新令牌有效期（秒）
	RefreshTokenValidity int `json:"refreshTokenValidity"`
	//自动授权（0否，1是）
	Autoapprove int16 `json:"autoapprove"`
	//资源权限，按服务名称逗号分割
	Resources string `json:"resources" binding:"max=2048"`
}

type UpdateOauthClientDTO struct {
	//ID
	Id *int64 `json:"id" binding:"required"`
	//客户端ID
	ClientId *string `json:"clientId" binding:"required,max=32"`
	//客户端秘钥
	ClientSecret *string `json:"clientSecret" binding:"max=256"`
	//客户端名称
	ClientName *string `json:"clientName" binding:"max=255"`
	//客户端LOGO
	LogoUri *string `json:"logoUri" binding:"max=255"`
	//授权范围
	Scope *string `json:"scope" binding:"max=256"`
	//授权类型
	AuthorizedGrantTypes *string `json:"authorizedGrantTypes" binding:"max=256"`
	//回调地址
	WebServerRedirectUri *string `json:"webServerRedirectUri" binding:"max=256"`
	//访问令牌有效期（秒）
	AccessTokenValidity *int `json:"accessTokenValidity"`
	//刷新令牌有效期（秒）
	RefreshTokenValidity *int `json:"refreshTokenValidity"`
	//自动授权（0否，1是）
	Autoapprove *int16 `json:"autoapprove"`
	//资源权限，按服务名称逗号分割
	Resources *string `json:"resources" binding:"max=2048"`
}
