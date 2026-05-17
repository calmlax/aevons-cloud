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
	"sys-service/internal/dto"
	"sys-service/internal/model"

	"github.com/calmlax/aevons-framework/core/base"

	"gorm.io/gorm"
)

// 继承BaseRepository
type UserRepository interface {
	base.BaseRepository[model.User]
	// 获取用户部门岗位关联列表
	GetUserDeptsByUserId(userId int64) ([]model.UserDept, error)
	// 事务执行用户相关操作
	UserTransaction(fn func(repo UserRepository) error) error
	// 用户角色关联
	GetRoleIdsByUserId(userId int64) ([]int64, error)
	// 设置用户角色关联（先删后插）
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

// GetUserDeptsByUserId 查询用户的用户部门关联列表。
func (r *userRepository) GetUserDeptsByUserId(userId int64) ([]model.UserDept, error) {
	var userDepts []model.UserDept
	if err := r.db.Where("user_id = ?", userId).Find(&userDepts).Error; err != nil {
		return nil, err
	}
	return userDepts, nil
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
