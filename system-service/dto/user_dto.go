/**
 * 用户信息表 DTO
 *
 * @author
 * @date 2026-04-19
 * @website https://www.aevons.com
 * @email 1073602@qq.com
 * @copyright 2025-2026 Aevons
 */
package dto

import (
	"time"

	"github.com/calmlax/aevons-framework/core/base"
)

// q 支持查询操作: eq, ne, gt, gte, lt, lte, like, like_l, like_r, in, not_in, between, is_null, not_null
// gorm.column 必须对应真实数据库列名
type UserQuery struct {

	// 用户名
	Username *string `form:"username" gorm:"column:username" q:"like"`
	// 呢称
	Nickname *string `form:"nickname" gorm:"column:nickname" q:"like"`
	// 用户类型
	Type *int16 `form:"type" gorm:"column:type" q:"eq"`
	// 邮箱
	Email *string `form:"email" gorm:"column:email" q:"like_l"`
	// 手机号
	Mobile *string `form:"mobile" gorm:"column:mobile" q:"like_l"`
	// 性别（0未知 1男 2女）
	Sex *int16 `form:"sex" gorm:"column:sex" q:"eq"`
	// 状态（0正常 1停用）
	Status *int16 `form:"status" gorm:"column:status" q:"eq"`
	base.BaseQuery
}

type CreateUserDTO struct {
	//用户名
	Username string `json:"username" binding:"required,max=32"`
	//呢称
	Nickname string `json:"nickname" binding:"required,max=50"`
	//用户类型
	Type int16 `json:"type"`
	//邮箱
	Email string `json:"email" binding:"max=64"`
	//手机号
	Mobile string `json:"mobile" binding:"max=11"`
	//性别（0未知 1男 2女）
	Sex int16 `json:"sex"`
	//头像
	Avatar string `json:"avatar" binding:"max=255"`
	//电子签名
	Autograph string `json:"autograph" binding:"max=255"`
	//密码
	Password string `json:"password" binding:"required,max=128"`
	//状态（0正常 1停用）
	Status int16 `json:"status"`
	//角色ID列表
	RoleIds []string `json:"roleIds"`
	//部门岗位列表
	DeptPosts []UserDeptPostDTO `json:"deptPosts"`
}

type UpdateUserDTO struct {
	//用户编号
	Id *int64 `json:"id,string" binding:"required"`
	//用户名
	Username *string `json:"username" binding:"required,max=32"`
	//呢称
	Nickname *string `json:"nickname" binding:"required,max=50"`
	//用户类型
	Type *int16 `json:"type"`
	//邮箱
	Email *string `json:"email" binding:"max=64"`
	//手机号
	Mobile *string `json:"mobile" binding:"max=11"`
	//性别（0未知 1男 2女）
	Sex *int16 `json:"sex"`
	//头像
	Avatar *string `json:"avatar" binding:"max=255"`
	//电子签名
	Autograph *string `json:"autograph" binding:"max=255"`
	//密码（为空则不修改）
	Password *string `json:"password" binding:"omitempty,max=128"`
	//状态（0正常 1停用）
	Status *int16 `json:"status"`
	//角色ID列表
	RoleIds []string `json:"roleIds"`
	//部门岗位列表
	DeptPosts []UserDeptPostDTO `json:"deptPosts"`
}

// UserDeptPostDTO 部门岗位关联
type UserDeptPostDTO struct {
	DeptId int64 `json:"deptId,string"`
	PostId int64 `json:"postId,string"`
}

// ResetPasswordDTO 重置密码
type ResetPasswordDTO struct {
	Password string `json:"password" binding:"required,max=128"`
}

type UserDTO struct {
	//用户编号
	Id int64 `excel:"column:用户编号;index:0;dict:" json:"id,string" binding:"required"`
	//用户名
	Username string `excel:"column:用户名;index:1;dict:" json:"username" binding:"required,max=32"`
	//呢称
	Nickname string `excel:"column:呢称;index:2;dict:" json:"nickname" binding:"required,max=50"`
	//用户类型
	Type int16 `excel:"column:用户类型;index:3;dict:sys_user_type" json:"type" binding:"required"`
	//邮箱
	Email string `excel:"column:邮箱;index:4;dict:" json:"email" binding:"max=64"`
	//手机号
	Mobile string `excel:"column:手机号;index:5;dict:" json:"mobile" binding:"max=11"`
	//性别（0未知 1男 2女）
	Sex int16 `excel:"column:性别（0未知 1男 2女）;index:6;dict:sys_user_sex" json:"sex" binding:"required"`
	//头像
	Avatar string `excel:"column:头像;index:7;dict:" json:"avatar" binding:"max=255"`
	//电子签名
	Autograph string `excel:"column:电子签名;index:8;dict:" json:"autograph" binding:"max=255"`
	//密码
	Password string `excel:"column:密码;index:9;dict:" json:"password" binding:"required,max=128"`
	//状态（0正常 1停用）
	Status int16 `excel:"column:状态（0正常 1停用）;index:10;dict:sys_common_status" json:"status" binding:"required"`
	//
	CreatedBy int64 `excel:"column:;index:11;dict:" json:"createdBy,string"`
	//
	CreatedAt time.Time `excel:"column:;index:12;dict:" json:"createdAt"`
	//
	UpdatedBy int64 `excel:"column:;index:13;dict:" json:"updatedBy,string"`
	//
	UpdatedAt time.Time `excel:"column:;index:14;dict:" json:"updatedAt"`
}
