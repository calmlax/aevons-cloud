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
	"sys-service/internal/dto"
	sysexcel "sys-service/internal/excel"
	"sys-service/internal/model"
	"sys-service/internal/service"

	"github.com/calmlax/aevons-framework/core/scope"
	apperr "github.com/calmlax/aevons-framework/errors"
	"github.com/calmlax/aevons-framework/excel"
	"github.com/calmlax/aevons-framework/response"

	"github.com/calmlax/aevons-framework/core/base"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	crud         *base.BaseHandler[model.User, *dto.UserQuery, dto.CreateUserDTO, dto.UpdateUserDTO]
	svc          service.UserService
	dictProvider sysexcel.DictProviderBuilder
}

func NewUserHandler(svc service.UserService, dictProvider sysexcel.DictProviderBuilder) *UserHandler {
	return &UserHandler{
		crud:         base.NewBaseHandler[model.User, *dto.UserQuery, dto.CreateUserDTO, dto.UpdateUserDTO](svc),
		svc:          svc,
		dictProvider: dictProvider,
	}
}

// List 查询用户列表。
//
// @Summary      查询用户列表
// @Tags         用户管理
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  response.Response
// @Router       /api/sys/v1/user/list [get]
func (h *UserHandler) List(c *gin.Context) {
	h.crud.HandleList(c)
}

// Page 分页查询用户。
//
// @Summary      分页查询用户
// @Tags         用户管理
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  response.Response
// @Router       /api/sys/v1/user/page [get]
func (h *UserHandler) Page(c *gin.Context) {
	h.crud.HandlePage(c)
}

// Get 获取用户详情。
//
// @Summary      获取用户详情
// @Tags         用户管理
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "用户ID"
// @Success      200  {object}  response.Response
// @Router       /api/sys/v1/user/{id} [get]
func (h *UserHandler) Get(c *gin.Context) {
	h.crud.HandleGet(c)
}

// BatchDelete 批量删除用户。
//
// @Summary      批量删除用户
// @Tags         用户管理
// @Produce      json
// @Security     BearerAuth
// @Param        ids   path      string  true  "逗号分隔的用户ID"
// @Success      200   {object}  response.Response
// @Router       /api/sys/v1/user/{ids} [delete]
func (h *UserHandler) BatchDelete(c *gin.Context) {
	h.crud.HandleBatchDelete(c)
}

// Create 新增用户（含角色、部门岗位，事务）
//
// @Summary      新增用户
// @Tags         用户管理
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      dto.CreateUserDTO  true  "新增用户"
// @Success      200      {object}  response.Response
// @Router       /api/sys/v1/user [post]
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
//
// @Summary      修改用户
// @Tags         用户管理
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      int                true  "用户ID"
// @Param        request  body      dto.UpdateUserDTO  true  "修改用户"
// @Success      200      {object}  response.Response
// @Router       /api/sys/v1/user/{id} [put]
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
//
// @Summary      获取用户关联信息
// @Tags         用户管理
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "用户ID"
// @Success      200  {object}  response.Response
// @Router       /api/sys/v1/user/{id}/relations [get]
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
//
// @Summary      重置用户密码
// @Tags         用户管理
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      int                   true  "用户ID"
// @Param        request  body      dto.ResetPasswordDTO  true  "重置密码"
// @Success      200      {object}  response.Response
// @Router       /api/sys/v1/user/{id}/reset-password [put]
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
//
// @Summary      更新用户状态
// @Tags         用户管理
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      int   true  "用户ID"
// @Param        request  body      object  true  "状态信息"
// @Success      200      {object}  response.Response
// @Router       /api/sys/v1/user/{id}/status [put]
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

// Export 导出用户 Excel。
//
// @Summary      导出用户
// @Tags         用户管理
// @Produce      application/octet-stream
// @Security     BearerAuth
// @Router       /api/sys/v1/user/export [get]
func (h *UserHandler) Export(c *gin.Context) {
	q, ok := base.BindQuery[*dto.UserQuery](c)
	if !ok {
		return
	}
	list, err := h.svc.ListExcel(q, scope.GetDBScopes(c)...)
	if err != nil {
		response.FailBy(c, apperr.ErrQueryFailed)
		return
	}
	if err := excel.Export(c, excel.ExportParam{
		DataList:     list,
		StructPtr:    &dto.UserDTO{},
		FileName:     "用户列表",
		SheetName:    "用户",
		DictProvider: h.dictProvider.Build(c),
	}); err != nil {
		response.FailServerError(c, "excel.export.failed", map[string]any{"error": err.Error()})
	}
}

// ImportTemplate 下载用户导入模板。
//
// @Summary      下载用户导入模板
// @Tags         用户管理
// @Produce      application/octet-stream
// @Security     BearerAuth
// @Router       /api/sys/v1/user/import/template [get]
func (h *UserHandler) ImportTemplate(c *gin.Context) {
	if err := excel.Export(c, excel.ExportParam{
		DataList:     []dto.UserDTO{},
		StructPtr:    &dto.UserDTO{},
		FileName:     "用户导入模板",
		SheetName:    "用户导入模板",
		DictProvider: h.dictProvider.Build(c),
	}); err != nil {
		response.FailServerError(c, "excel.export.failed", map[string]any{"error": err.Error()})
	}
}

// Import 解析用户导入文件。
//
// @Summary      导入用户 Excel
// @Tags         用户管理
// @Accept       mpfd
// @Produce      json
// @Security     BearerAuth
// @Param        file  formData  file  true  "Excel 文件"
// @Success      200   {object}  response.Response
// @Router       /api/sys/v1/user/import [post]
func (h *UserHandler) Import(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		response.FailBy(c, apperr.ErrInvalidBody)
		return
	}
	result, err := excel.Import(c.Request.Context(), excel.ImportParam{
		File:         file,
		StructPtr:    &dto.UserDTO{},
		DictProvider: h.dictProvider.Build(c),
	})
	if err != nil {
		response.FailServerError(c, "excel.import.failed", map[string]any{"error": err.Error()})
		return
	}
	response.Success(c, result)
}
