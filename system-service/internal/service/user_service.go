/**
 * 用户信息表 Service
 *
 * @author
 * @date 2026-04-19
 * @website https://www.aevons.com
 * @email 1073602@qq.com
 * @copyright 2025-2026 Aevons
 */
package service

import (
	"strconv"
	"system-service/internal/dto"
	"system-service/internal/model"
	"system-service/internal/repository"

	"github.com/calmlax/aevons-framework/core/base"

	"golang.org/x/crypto/bcrypt"
)

// 继承BaseService
type UserService interface {
	base.BaseService[model.User, *dto.UserQuery]
	// 新增用户（含角色、部门岗位，事务）
	CreateWithRelations(d dto.CreateUserDTO) (*model.User, error)
	// 修改用户（含角色、部门岗位，事务）
	UpdateWithRelations(id int64, d dto.UpdateUserDTO) error
	// 获取用户已关联角色ID
	GetRoleIds(userId int64) ([]int64, error)
	// 获取用户已关联部门岗位
	GetDeptPosts(userId int64) ([]model.UserDept, error)
	// 重置密码
	ResetPassword(userId int64, newPassword string) error
	// 更新状态
	UpdateStatus(userId int64, status int16) error
}

type userService struct {
	base.BaseService[model.User, *dto.UserQuery]
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	baseSrv := base.NewBaseService[model.User, *dto.UserQuery](repo)
	return &userService{
		BaseService: baseSrv,
		repo:        repo,
	}
}

func parseRoleIds(ids []string) []int64 {
	result := make([]int64, 0, len(ids))
	for _, s := range ids {
		if id, err := strconv.ParseInt(s, 10, 64); err == nil {
			result = append(result, id)
		}
	}
	return result
}

func (s *userService) CreateWithRelations(d dto.CreateUserDTO) (*model.User, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(d.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	var user *model.User
	err = s.repo.UserTransaction(func(repo repository.UserRepository) error {
		u := &model.User{
			Username:  d.Username,
			Nickname:  d.Nickname,
			Type:      d.Type,
			Email:     d.Email,
			Mobile:    d.Mobile,
			Sex:       d.Sex,
			Avatar:    d.Avatar,
			Autograph: d.Autograph,
			Password:  string(hashed),
			Status:    d.Status,
		}
		if err := repo.Create(u); err != nil {
			return err
		}
		if err := repo.SetUserRoles(u.Id, parseRoleIds(d.RoleIds)); err != nil {
			return err
		}
		if err := repo.SetUserDeptPosts(u.Id, d.DeptPosts); err != nil {
			return err
		}
		user = u
		return nil
	})
	return user, err
}

func (s *userService) UpdateWithRelations(id int64, d dto.UpdateUserDTO) error {
	return s.repo.UserTransaction(func(repo repository.UserRepository) error {
		// 手动构建只含用户表字段的 updates，避免 DeptPosts/RoleIds struct 污染
		updates := map[string]any{}
		if d.Username != nil {
			updates["username"] = *d.Username
		}
		if d.Nickname != nil {
			updates["nickname"] = *d.Nickname
		}
		if d.Type != nil {
			updates["type"] = *d.Type
		}
		if d.Email != nil {
			updates["email"] = *d.Email
		}
		if d.Mobile != nil {
			updates["mobile"] = *d.Mobile
		}
		if d.Sex != nil {
			updates["sex"] = *d.Sex
		}
		if d.Avatar != nil {
			updates["avatar"] = *d.Avatar
		}
		if d.Autograph != nil {
			updates["autograph"] = *d.Autograph
		}
		if d.Status != nil {
			updates["status"] = *d.Status
		}
		if d.Password != nil && *d.Password != "" {
			hashed, err := bcrypt.GenerateFromPassword([]byte(*d.Password), bcrypt.DefaultCost)
			if err != nil {
				return err
			}
			updates["password"] = string(hashed)
		}

		if len(updates) > 0 {
			if _, err := repo.Update(id, updates); err != nil {
				return err
			}
		}
		if err := repo.SetUserRoles(id, parseRoleIds(d.RoleIds)); err != nil {
			return err
		}
		return repo.SetUserDeptPosts(id, d.DeptPosts)
	})
}

func (s *userService) GetRoleIds(userId int64) ([]int64, error) {
	return s.repo.GetRoleIdsByUserId(userId)
}

func (s *userService) GetDeptPosts(userId int64) ([]model.UserDept, error) {
	return s.repo.GetUserDeptsByUserId(userId)
}

func (s *userService) ResetPassword(userId int64, newPassword string) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	return s.repo.ResetPassword(userId, string(hashed))
}

func (s *userService) UpdateStatus(userId int64, status int16) error {
	return s.repo.UpdateStatus(userId, status)
}
