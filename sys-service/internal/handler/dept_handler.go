package handler

import (
	"sys-service/internal/dto"
	"sys-service/internal/model"
	"sys-service/internal/service"

	apperr "github.com/calmlax/aevons-framework/errors"
	"github.com/calmlax/aevons-framework/response"

	"github.com/calmlax/aevons-framework/core/base"

	"github.com/gin-gonic/gin"
)

type DeptHandler struct {
	crud *base.BaseHandler[model.Dept, *dto.DeptQuery, dto.CreateDeptDTO, dto.UpdateDeptDTO]
	svc  service.DeptService
}

func NewDeptHandler(svc service.DeptService) *DeptHandler {
	return &DeptHandler{
		crud: base.NewBaseHandler[model.Dept, *dto.DeptQuery, dto.CreateDeptDTO, dto.UpdateDeptDTO](svc),
		svc:  svc,
	}
}

// Get 获取部门详情。
//
// @Summary      获取部门详情
// @Tags         部门管理
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "部门ID"
// @Success      200  {object}  response.Response
// @Router       /dept/{id} [get]
func (h *DeptHandler) Get(c *gin.Context) {
	h.crud.HandleGet(c)
}

// ListTree 返回树形部门列表
//
// @Summary      查询部门树
// @Tags         部门管理
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  response.Response
// @Router       /dept/list [get]
func (h *DeptHandler) ListTree(c *gin.Context) {
	var q dto.DeptQuery
	_ = c.ShouldBindQuery(&q)
	list, err := h.svc.ListTree(q)
	if err != nil {
		response.FailBy(c, apperr.ErrQueryFailed)
		return
	}
	response.Success(c, list)
}

// Create 创建部门（自动计算 ancestors）
//
// @Summary      新增部门
// @Tags         部门管理
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      dto.CreateDeptDTO  true  "新增部门"
// @Success      200      {object}  response.Response
// @Router       /dept [post]
func (h *DeptHandler) Create(c *gin.Context) {
	var d dto.CreateDeptDTO
	if err := c.ShouldBindJSON(&d); err != nil {
		response.FailBy(c, apperr.ErrInvalidBody)
		return
	}
	dept, err := h.svc.CreateDept(c, d)
	if err != nil {
		response.FailBy(c, apperr.ErrCreateFailed, map[string]any{"error": err.Error()})
		return
	}
	response.Success(c, dept)
}

// Update 更新部门（自动重算 ancestors）
//
// @Summary      修改部门
// @Tags         部门管理
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      int               true  "部门ID"
// @Param        request  body      dto.UpdateDeptDTO true  "修改部门"
// @Success      200      {object}  response.Response
// @Router       /dept/{id} [put]
func (h *DeptHandler) Update(c *gin.Context) {
	id, ok := base.GetId(c)
	if !ok {
		return
	}
	var d dto.UpdateDeptDTO
	if err := c.ShouldBindJSON(&d); err != nil {
		response.FailBy(c, apperr.ErrInvalidBody)
		return
	}
	if err := h.svc.UpdateDept(c, id, d); err != nil {
		response.FailBy(c, apperr.ErrUpdateFailed, map[string]any{"error": err.Error()})
		return
	}
	response.Success(c, id)
}

// Delete 删除部门（检查是否有子部门）
//
// @Summary      删除部门
// @Tags         部门管理
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "部门ID"
// @Success      200  {object}  response.Response
// @Router       /dept/{id} [delete]
func (h *DeptHandler) Delete(c *gin.Context) {
	id, ok := base.GetId(c)
	if !ok {
		return
	}
	has, err := h.svc.HasChildren(id)
	if err != nil || has {
		response.Fail(c, 400, 4001, "err.sys.dept.has_children")
		return
	}
	if err := h.svc.Delete(id); err != nil {
		response.FailBy(c, apperr.ErrDeleteFailed)
		return
	}
	response.Success(c, id)
}
