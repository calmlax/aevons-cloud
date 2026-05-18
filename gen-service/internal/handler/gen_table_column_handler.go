package handler

import (
	"gen-service/internal/dto"
	"gen-service/internal/model"
	"gen-service/internal/service"

	"github.com/calmlax/aevons-framework/core/base"
	apperr "github.com/calmlax/aevons-framework/errors"
	"github.com/calmlax/aevons-framework/response"
	"github.com/calmlax/aevons-framework/utils"

	"github.com/gin-gonic/gin"
)

type GenTableColumnHandler struct {
	crud *base.BaseHandler[model.GenTableColumn, *dto.GenTableColumnQuery, dto.CreateGenTableColumnDTO, dto.UpdateGenTableColumnDTO]
	svc  service.GenTableColumnService
}

func NewGenTableColumnHandler(svc service.GenTableColumnService) *GenTableColumnHandler {
	return &GenTableColumnHandler{
		crud: base.NewBaseHandler[model.GenTableColumn, *dto.GenTableColumnQuery, dto.CreateGenTableColumnDTO, dto.UpdateGenTableColumnDTO](svc),
		svc:  svc,
	}
}

// List 查询代码生成字段列表。
//
// @Summary      查询代码生成字段列表
// @Tags         代码生成字段
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  response.Response
// @Router       /table/column/list [get]
func (h *GenTableColumnHandler) List(c *gin.Context) {
	h.crud.HandleList(c)
}

// BatchUpdate 批量更新数据表字段信息。
//
// @Summary      批量更新代码生成字段
// @Tags         代码生成字段
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      []dto.UpdateGenTableColumnDTO  true  "批量更新字段信息"
// @Success      200      {object}  response.Response
// @Router       /table/column/batch-update [put]
func (h *GenTableColumnHandler) BatchUpdate(c *gin.Context) {
	var reqList []dto.UpdateGenTableColumnDTO
	if err := c.ShouldBindJSON(&reqList); err != nil {
		errDef, arg := utils.GetValidateErrorKey(err)
		response.FailBy(c, *errDef, map[string]any{"field": arg})
		return
	}

	ok, err := h.svc.BatchUpdate(reqList)
	if err != nil {
		response.FailBy(c, apperr.ErrUpdateFailed, map[string]any{"error": err.Error()})
		return
	}
	response.Success(c, ok)
}
