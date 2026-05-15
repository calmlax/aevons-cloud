/**
 * 用户信息表 Handler
 *
 * @author
 * @date 2026-04-19
 * @website https://www.aevons.com
 * @email 1073602@qq.com
 * @copyright 2025-2026 Aevons
 */
package handler

import (
	"strconv"
	"system-service/dto"
	"system-service/model"
	"system-service/service"

	apperr "github.com/calmlax/aevons-framework/errors"
	"github.com/calmlax/aevons-framework/response"

	"github.com/calmlax/aevons-framework/core/base"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	*base.BaseHandler[
		model.User,
		*dto.UserQuery,
		dto.CreateUserDTO,
		dto.UpdateUserDTO,
	]
	svc service.UserService
}

func NewUserHandler(svc service.UserService) *UserHandler {
	return &UserHandler{
		BaseHandler: base.NewBaseHandler[
			model.User,
			*dto.UserQuery,
			dto.CreateUserDTO,
			dto.UpdateUserDTO,
		](svc),
		svc: svc,
	}
}

// Create 新增用户（含角色、部门岗位，事务）
func (h *UserHandler) Create(c *gin.Context) {
	var d dto.CreateUserDTO
	if err := c.ShouldBindJSON(&d); err != nil {
		response.FailBy(c, apperr.ErrInvalidBody)
		return
	}
	user, err := h.svc.CreateWithRelations(d)
	if err != nil {
		response.FailBy(c, apperr.ErrCreateFailed, map[string]any{"error": err.Error()})
		return
	}
	response.Success(c, user)
}

// Update 修改用户（含角色、部门岗位，事务）
func (h *UserHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.FailBy(c, apperr.ErrInvalidId)
		return
	}
	var d dto.UpdateUserDTO
	if err := c.ShouldBindJSON(&d); err != nil {
		response.FailBy(c, apperr.ErrInvalidBody)
		return
	}
	if err := h.svc.UpdateWithRelations(id, d); err != nil {
		response.FailBy(c, apperr.ErrUpdateFailed, map[string]any{"error": err.Error()})
		return
	}
	response.Success(c, id)
}

// GetRelations 获取用户已关联的角色ID和部门岗位
func (h *UserHandler) GetRelations(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.FailBy(c, apperr.ErrInvalidId)
		return
	}
	roleIds, err := h.svc.GetRoleIds(id)
	if err != nil {
		response.FailBy(c, apperr.ErrQueryFailed)
		return
	}
	deptPosts, err := h.svc.GetDeptPosts(id)
	if err != nil {
		response.FailBy(c, apperr.ErrQueryFailed)
		return
	}
	roleIdStrs := make([]string, len(roleIds))
	for i, rid := range roleIds {
		roleIdStrs[i] = strconv.FormatInt(rid, 10)
	}

	response.Success(c, gin.H{
		"roleIds":   roleIdStrs,
		"deptPosts": deptPosts,
	})
}

// ResetPassword 重置用户密码
func (h *UserHandler) ResetPassword(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.FailBy(c, apperr.ErrInvalidId)
		return
	}
	var d dto.ResetPasswordDTO
	if err := c.ShouldBindJSON(&d); err != nil {
		response.FailBy(c, apperr.ErrInvalidBody)
		return
	}
	if err := h.svc.ResetPassword(id, d.Password); err != nil {
		response.FailBy(c, apperr.ErrUpdateFailed, map[string]any{"error": err.Error()})
		return
	}
	response.Success(c, id)
}

// UpdateStatus 更新用户状态
func (h *UserHandler) UpdateStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.FailBy(c, apperr.ErrInvalidId)
		return
	}
	var body struct {
		Status int16 `json:"status"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		response.FailBy(c, apperr.ErrInvalidBody)
		return
	}
	if err := h.svc.UpdateStatus(id, body.Status); err != nil {
		response.FailBy(c, apperr.ErrUpdateFailed, map[string]any{"error": err.Error()})
		return
	}
	response.Success(c, id)
}
