package dto

import "system-service/model"

// UpdateProfileDTO 个人资料更新请求（仅允许修改自己的基本信息）
type UpdateProfileDTO struct {
	Nickname  *string `json:"nickname"  binding:"omitempty,max=50"`
	Email     *string `json:"email"     binding:"omitempty,email,max=64"`
	Mobile    *string `json:"mobile"    binding:"omitempty,max=11"`
	Sex       *int16  `json:"sex"`
	Avatar    *string `json:"avatar"    binding:"omitempty,max=255"`
	Autograph *string `json:"autograph" binding:"omitempty,max=255"`
}

// UserProfile 返回给个人的完整用户信息。
type UserProfile struct {
	User        model.User        `json:"user"`
	Roles       []ProfileRole     `json:"roles"`
	DeptPosts   []ProfileDeptPost `json:"dept_posts"`
	Permissions []string          `json:"permissions"`
}

type ProfileRole struct {
	Id       int64  `json:"id,string"`
	RoleKey  string `json:"role_key"`
	RoleName string `json:"role_name"`
}

type ProfileDeptPost struct {
	DeptId   int64  `json:"dept_id,string"`
	DeptName string `json:"dept_name"`
	PostId   int64  `json:"post_id,string"`
	PostName string `json:"post_name"`
}
