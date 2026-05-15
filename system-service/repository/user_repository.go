/**
 * 用户信息表 Repository
 *
 * @author
 * @date 2026-04-09 02:01:11.845628976 +0000 UTC
 * @website https://www.aevons.com
 * @email 1073602@qq.com
 * @copyright 2025-2026 Aevons
 */
package repository

import (
	"system-service/dto"
	"system-service/model"

	"github.com/calmlax/aevons-framework/core/base"

	"gorm.io/gorm"
)

// 继承BaseRepository
type UserRepository interface {
	base.BaseRepository[model.User]
	GetByUsername(username string) (*model.User, error)
	GetByEmail(email string) (*model.User, error)
	GetRolesByUserId(userId int64) ([]model.Role, error)
	GetUserDeptsByUserId(userId int64) ([]model.UserDept, error)
	GetRoleDeptIdsByRoleId(roleId int64) ([]int64, error)
	GetPermissionsByRoleIds(roleIds []int64) ([]string, error)
	GetMenusByRoleIds(roleIds []int64, langCode string) ([]dto.MenuDTO, error)
	GetAllMenus(langCode string) ([]dto.MenuDTO, error)
	GetProfileDeptPosts(userId int64) ([]dto.ProfileDeptPost, error)
	UserTransaction(fn func(repo UserRepository) error) error
	// 用户角色关联
	GetRoleIdsByUserId(userId int64) ([]int64, error)
	SetUserRoles(userId int64, roleIds []int64) error
	// 用户部门岗位关联
	SetUserDeptPosts(userId int64, deptPosts []dto.UserDeptPostDTO) error
	// 重置密码
	ResetPassword(userId int64, hashedPassword string) error
	// 更新状态
	UpdateStatus(userId int64, status int16) error
}

// userRepository 结构体
type userRepository struct {
	base.BaseRepository[model.User]
	db *gorm.DB
}

// 创建 UserRepository 实例
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		BaseRepository: base.NewBaseRepository[model.User](db),
		db:             db,
	}
}

// Transaction 在事务中执行用户操作。
func (r *userRepository) UserTransaction(fn func(repo UserRepository) error) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		txRepo := &userRepository{
			BaseRepository: base.NewBaseRepository[model.User](tx),
			db:             tx,
		}
		return fn(txRepo)
	})
}

// GetByUsername 根据用户名查询用户。
func (r *userRepository) GetByUsername(username string) (*model.User, error) {
	var user model.User
	if err := r.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByEmail 根据邮箱查询用户。
func (r *userRepository) GetByEmail(email string) (*model.User, error) {
	var user model.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetRolesByUserId 查询用户的角色列表。
// 路径：sys_user_role → sys_role
func (r *userRepository) GetRolesByUserId(userId int64) ([]model.Role, error) {
	var roles []model.Role
	err := r.db.Table("sys_role").
		Select("sys_role.*").
		Joins("INNER JOIN sys_user_role ON sys_user_role.role_id = sys_role.id").
		Where("sys_user_role.user_id = ? AND sys_role.status = 0", userId).Find(&roles).Error
	return roles, err
}

// GetUserDeptsByUserId 查询用户的用户部门关联列表。
func (r *userRepository) GetUserDeptsByUserId(userId int64) ([]model.UserDept, error) {
	var userDepts []model.UserDept
	if err := r.db.Where("user_id = ?", userId).Find(&userDepts).Error; err != nil {
		return nil, err
	}
	return userDepts, nil
}

// GetRoleDeptIdsByRoleId 查询角色的部门ID列表。
// 路径：sys_role_dept → sys_dept
func (r *userRepository) GetRoleDeptIdsByRoleId(roleId int64) ([]int64, error) {
	var roleDeptIds []int64
	if err := r.db.Table("sys_role_dept").Where("role_id = ?", roleId).Pluck("dept_id", &roleDeptIds).Error; err != nil {
		return nil, err
	}
	return roleDeptIds, nil
}

// GetMenusByRoleIds 查询角色关联的菜单（父节点已在保存时补全，直接 IN 查询）。
func (r *userRepository) GetMenusByRoleIds(roleIds []int64, langCode string) ([]dto.MenuDTO, error) {
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

// GetAllMenus 查询所有菜单（管理员）。
func (r *userRepository) GetAllMenus(langCode string) ([]dto.MenuDTO, error) {
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

// GetPermissionsByRoleIds 查询用户的权限标识列表（菜单 permission）。
// 路径：sys_role_menu → sys_menu.permission
func (r *userRepository) GetPermissionsByRoleIds(roleIds []int64) ([]string, error) {
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

// GetProfileDeptPosts 获取用户的部门岗位详细信息映射。
func (r *userRepository) GetProfileDeptPosts(userId int64) ([]dto.ProfileDeptPost, error) {
	var results []dto.ProfileDeptPost
	err := r.db.Table("sys_user_dept ud").
		Select("d.id as dept_id, d.dept_name as dept_name, p.id as post_id, p.post_name as post_name").
		Joins("INNER JOIN sys_dept d ON ud.dept_id = d.id").
		Joins("INNER JOIN sys_post p ON ud.post_id = p.id").
		Where("ud.user_id = ?", userId).
		Find(&results).Error
	return results, err
}

// GetRoleIdsByUserId 查询用户已关联的角色ID列表。
func (r *userRepository) GetRoleIdsByUserId(userId int64) ([]int64, error) {
	var roleIds []int64
	err := r.db.Table("sys_user_role").Where("user_id = ?", userId).Pluck("role_id", &roleIds).Error
	return roleIds, err
}

// SetUserRoles 替换用户角色关联（先删后插）。
func (r *userRepository) SetUserRoles(userId int64, roleIds []int64) error {
	if err := r.db.Where("user_id = ?", userId).Delete(&model.UserRole{}).Error; err != nil {
		return err
	}
	if len(roleIds) == 0 {
		return nil
	}
	records := make([]model.UserRole, len(roleIds))
	for i, rid := range roleIds {
		records[i] = model.UserRole{UserId: userId, RoleId: rid}
	}
	return r.db.Create(&records).Error
}

// SetUserDeptPosts 替换用户部门岗位关联（先删后插）。
func (r *userRepository) SetUserDeptPosts(userId int64, deptPosts []dto.UserDeptPostDTO) error {
	if err := r.db.Where("user_id = ?", userId).Delete(&model.UserDept{}).Error; err != nil {
		return err
	}
	if len(deptPosts) == 0 {
		return nil
	}
	records := make([]model.UserDept, len(deptPosts))
	for i, dp := range deptPosts {
		records[i] = model.UserDept{UserId: userId, DeptId: dp.DeptId, PostId: dp.PostId}
	}
	return r.db.Create(&records).Error
}

// ResetPassword 重置用户密码。
func (r *userRepository) ResetPassword(userId int64, hashedPassword string) error {
	return r.db.Model(&model.User{}).Where("id = ?", userId).Update("password", hashedPassword).Error
}

// UpdateStatus 更新用户状态。
func (r *userRepository) UpdateStatus(userId int64, status int16) error {
	return r.db.Model(&model.User{}).Where("id = ?", userId).Update("status", status).Error
}
