package repository

import (
	"auth-service/dto"
	"auth-service/model"
	"strings"
	"time"

	"github.com/calmlax/aevons-framework/auth"
	"github.com/calmlax/aevons-framework/consts"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthRepository interface {
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
	// GetByClientId 按 Client ID 查找 OAuth 客户端
	GetByClientId(clientId string) (*model.OauthClient, error)
	// ValidateClient 校验客户端参数，返回客户端信息或认证错误
	ValidateClient(clientId, clientSecret, grantType string) (*model.OauthClient, error)
	// GetUserByEmail 根据邮箱获取用户信息
	GetUserByEmail(email string) (*model.User, error)
	// GetUserByUserId 根据用户名查询用户
	GetUserByUsername(username string) (*model.User, error)
	// GetUserByUserId 根据用户ID查询用户
	GetUserByUserId(userId int64) (*model.User, error)
	// RegisterUser 注册新用户
	RegisterUser(user *model.User) error
	// UpdateProfile 更新用户资料
	UpdateProfile(userId int64, profile map[string]any) error
	// GetRolesByUserId 获取用户角色列表
	GetRolesByUserId(userId int64) ([]model.Role, error)
	// GetUserDeptsByUserId 获取用户部门列表
	GetUserDeptsByUserId(userId int64) ([]model.UserDept, error)
	// GetRoleDeptIdsByRoleId 获取角色关联的部门ID列表
	GetRoleDeptIdsByRoleId(roleId int64) ([]int64, error)
	// GetPermissionsByRoleIds 获取角色关联的权限字符串列表
	GetPermissionsByRoleIds(roleIds []int64) ([]string, error)
	// GetAllMenus 获取所有菜单列表
	GetAllMenus(langCode string) ([]dto.MenuDTO, error)
	// GetMenusByRoleIds 获取角色关联的菜单列表
	GetMenusByRoleIds(roleIds []int64, langCode string) ([]dto.MenuDTO, error)
	// GetProfileDeptPosts 获取用户资料页的部门岗位信息
	GetProfileDeptPosts(userId int64) ([]dto.ProfileDeptPost, error)
}

type authRepository struct{ db *gorm.DB }

func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &authRepository{db: db}
}

func (r *authRepository) GetCredentialByUserId(userId int64) ([]*model.UserCredential, error) {
	var list []*model.UserCredential
	err := r.db.Where("user_id = ? AND is_revoked = 0", userId).Find(&list).Error
	return list, err
}

func (r *authRepository) GetCredentialByCredentialId(credentialId []byte) (*model.UserCredential, error) {
	var c model.UserCredential
	err := r.db.Where("credential_id = ? AND is_revoked = 0", credentialId).First(&c).Error
	return &c, err
}

func (r *authRepository) CreateCredential(c *model.UserCredential) error {
	return r.db.Create(c).Error
}

func (r *authRepository) UpdateCredentialSignatureCount(id int64, count uint64) error {
	now := time.Now()
	return r.db.Model(&model.UserCredential{}).Where("id = ?", id).
		Updates(map[string]any{"signature_count": count, "last_used_at": now}).Error
}

func (r *authRepository) RevokeCredential(id int64) error {
	return r.db.Model(&model.UserCredential{}).Where("id = ?", id).
		Update("is_revoked", true).Error
}

func (r *authRepository) ListCredentialByUserId(userId int64) ([]*model.UserCredential, error) {
	var list []*model.UserCredential
	err := r.db.Where("user_id = ?", userId).Order("created_at DESC").Find(&list).Error
	return list, err
}

func (r *authRepository) GetByClientId(clientId string) (*model.OauthClient, error) {
	var c model.OauthClient
	err := r.db.Where("client_id = ?", clientId).First(&c).Error
	return &c, err
}

func (r *authRepository) ValidateClient(clientId, clientSecret, grantType string) (*model.OauthClient, error) {
	client, err := r.GetByClientId(clientId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, &auth.AuthError{Code: consts.ErrOAuthInvalidClient, HTTPStatus: 401}
		}
		return nil, err
	}
	if clientSecret != "" {
		if err := bcrypt.CompareHashAndPassword([]byte(client.ClientSecret), []byte(clientSecret)); err != nil {
			return nil, &auth.AuthError{Code: consts.ErrOAuthInvalidClient, HTTPStatus: 401}
		}
	}
	if !r.isGrantTypeAllowed(client.AuthorizedGrantTypes, grantType) {
		return nil, &auth.AuthError{Code: consts.ErrOAuthUnsupportedGrant, HTTPStatus: 400}
	}
	return client, nil
}

func (r *authRepository) isGrantTypeAllowed(authorizedGrantTypes, grantType string) bool {
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

// GetUserByEmail 根据邮箱查询用户。
func (r *authRepository) GetUserByEmail(email string) (*model.User, error) {
	var user model.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByUsername 根据用户名查询用户。
func (r *authRepository) GetUserByUsername(username string) (*model.User, error) {
	var user model.User
	if err := r.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByUserId 根据用户ID查询用户。
func (r *authRepository) GetUserByUserId(userId int64) (*model.User, error) {
	var user model.User
	if err := r.db.Where("id = ?", userId).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// RegisterUser 注册新用户。
func (r *authRepository) RegisterUser(user *model.User) error {
	return r.db.Create(user).Error
}

// UpdateProfile 更新用户资料。
func (r *authRepository) UpdateProfile(userId int64, profile map[string]any) error {
	return r.db.Model(&model.User{}).Where("id = ?", userId).Updates(profile).Error
}

// GetRolesByUserId 查询用户的角色列表。
// 路径：sys_user_role → sys_role
func (r *authRepository) GetRolesByUserId(userId int64) ([]model.Role, error) {
	var roles []model.Role
	err := r.db.Table("sys_role").
		Select("sys_role.*").
		Joins("INNER JOIN sys_user_role ON sys_user_role.role_id = sys_role.id").
		Where("sys_user_role.user_id = ? AND sys_role.status = 0", userId).Find(&roles).Error
	return roles, err
}

// GetUserDeptsByUserId 查询用户的用户部门关联列表。
func (r *authRepository) GetUserDeptsByUserId(userId int64) ([]model.UserDept, error) {
	var userDepts []model.UserDept
	if err := r.db.Where("user_id = ?", userId).Find(&userDepts).Error; err != nil {
		return nil, err
	}
	return userDepts, nil
}

// GetRoleDeptIdsByRoleId 查询角色的部门ID列表。
// 路径：sys_role_dept → sys_dept
func (r *authRepository) GetRoleDeptIdsByRoleId(roleId int64) ([]int64, error) {
	var roleDeptIds []int64
	if err := r.db.Table("sys_role_dept").Where("role_id = ?", roleId).Pluck("dept_id", &roleDeptIds).Error; err != nil {
		return nil, err
	}
	return roleDeptIds, nil
}

// GetPermissionsByRoleIds 查询用户的权限标识列表（菜单 permission）。
// 路径：sys_role_menu → sys_menu.permission
func (r *authRepository) GetPermissionsByRoleIds(roleIds []int64) ([]string, error) {
	if len(roleIds) == 0 {
		return []string{}, nil
	}
	var permissions []string
	err := r.db.Table("sys_menu m").
		Select("DISTINCT m.permission").
		Joins("INNER JOIN sys_role_menu ON sys_role_menu.menu_id = m.id").
		Where("sys_role_menu.role_id IN ? AND m.permission != ''", roleIds).
		Pluck("m.permission", &permissions).Error
	return permissions, err
}

// GetAllMenus 查询所有菜单（管理员）。
func (r *authRepository) GetAllMenus(langCode string) ([]dto.MenuDTO, error) {
	var menus []dto.MenuDTO
	err := r.db.Table("sys_menu m").
		Select("m.id,m.parent_id,tl.lang_code,tl.title,m.type,m.sort,m.path,m.component,m.`query`,m.visible,m.`status`,m.is_frame,m.permission,m.icon,m.active_id").
		Joins("LEFT JOIN sys_menu_tl tl ON m.id = tl.menu_id AND tl.lang_code = ?", langCode).
		Distinct().
		Where("m.type IN (1, 2) AND m.status = 0").
		Order("m.parent_id, m.sort").
		Find(&menus).Error
	return menus, err
}

// GetMenusByRoleIds 查询角色关联的菜单（父节点已在保存时补全，直接 IN 查询）。
func (r *authRepository) GetMenusByRoleIds(roleIds []int64, langCode string) ([]dto.MenuDTO, error) {
	if len(roleIds) == 0 {
		return []dto.MenuDTO{}, nil
	}
	var menus []dto.MenuDTO
	err := r.db.Table("sys_menu m").
		Select("m.id, m.parent_id, tl.lang_code, tl.title, m.type, m.sort, m.path, m.component, m.`query`, m.visible, m.`status`, m.is_frame, m.permission, m.icon, m.active_id").
		Joins("INNER JOIN sys_role_menu rm ON rm.menu_id = m.id").
		Joins("LEFT JOIN sys_menu_tl tl ON m.id = tl.menu_id AND tl.lang_code = ?", langCode).
		Where("rm.role_id IN ? AND m.type IN (1,2) AND m.status = 0", roleIds).
		Distinct().
		Order("m.parent_id, m.sort").
		Find(&menus).Error
	return menus, err
}

// GetProfileDeptPosts 获取用户的部门岗位详细信息映射。
func (r *authRepository) GetProfileDeptPosts(userId int64) ([]dto.ProfileDeptPost, error) {
	var results []dto.ProfileDeptPost
	err := r.db.Table("sys_user_dept ud").
		Select("d.id as dept_id, d.dept_name as dept_name, p.id as post_id, p.post_name as post_name").
		Joins("INNER JOIN sys_dept d ON ud.dept_id = d.id").
		Joins("INNER JOIN sys_post p ON ud.post_id = p.id").
		Where("ud.user_id = ?", userId).
		Find(&results).Error
	return results, err
}
