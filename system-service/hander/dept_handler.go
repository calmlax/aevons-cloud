package handler

import (
	"system-service/dto"
	"system-service/model"
	"system-service/service"

	apperr "github.com/calmlax/aevons-framework/errors"
	"github.com/calmlax/aevons-framework/response"

	"github.com/calmlax/aevons-framework/core/base"

	"github.com/gin-gonic/gin"
)

type DeptHandler struct {
	*base.BaseHandler[
		model.Dept,
		*dto.DeptQuery,
		dto.CreateDeptDTO,
		dto.UpdateDeptDTO,
	]
	svc service.DeptService
}

func NewDeptHandler(svc service.DeptService) *DeptHandler {
	return &DeptHandler{
		BaseHandler: base.NewBaseHandler[
			model.Dept,
			*dto.DeptQuery,
			dto.CreateDeptDTO,
			dto.UpdateDeptDTO,
		](svc),
		svc: svc,
	}
}

// ListTree 返回树形部门列表
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
func (h *DeptHandler) Update(c *gin.Context) {
	id, ok := h.GetId(c)
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
func (h *DeptHandler) Delete(c *gin.Context) {
	id, ok := h.GetId(c)
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
