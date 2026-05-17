package handler

import (
	"net/http"
	"strconv"

	"gen-service/dto"
	"gen-service/model"
	"gen-service/service"

	"github.com/calmlax/aevons-framework/core/base"
	apperr "github.com/calmlax/aevons-framework/errors"
	"github.com/calmlax/aevons-framework/response"
	"github.com/calmlax/aevons-framework/utils"

	"github.com/gin-gonic/gin"
)

type GenTableHandler struct {
	crud *base.BaseHandler[model.GenTable, *dto.GenTableQuery, dto.CreateGenTableDTO, dto.UpdateGenTableDTO]
	svc  service.GenTableService
}

func NewGenTableHandler(svc service.GenTableService) *GenTableHandler {
	return &GenTableHandler{
		crud: base.NewBaseHandler[model.GenTable, *dto.GenTableQuery, dto.CreateGenTableDTO, dto.UpdateGenTableDTO](svc),
		svc:  svc,
	}
}

// List 查询代码生成表列表。
//
// @Summary      查询代码生成表列表
// @Tags         代码生成
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  response.Response
// @Router       /table/list [get]
func (h *GenTableHandler) List(c *gin.Context) {
	h.crud.HandleList(c)
}

// Page 分页查询代码生成表。
//
// @Summary      分页查询代码生成表
// @Tags         代码生成
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  response.Response
// @Router       /table/page [get]
func (h *GenTableHandler) Page(c *gin.Context) {
	h.crud.HandlePage(c)
}

// Get 获取代码生成表详情。
//
// @Summary      获取代码生成表详情
// @Tags         代码生成
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "代码生成表ID"
// @Success      200  {object}  response.Response
// @Router       /table/{id} [get]
func (h *GenTableHandler) Get(c *gin.Context) {
	h.crud.HandleGet(c)
}

// Create 新增代码生成表。
//
// @Summary      新增代码生成表
// @Tags         代码生成
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      dto.CreateGenTableDTO  true  "新增代码生成表"
// @Success      200      {object}  response.Response
// @Router       /table [post]
func (h *GenTableHandler) Create(c *gin.Context) {
	h.crud.HandleCreate(c)
}

// Update 更新代码生成表。
//
// @Summary      修改代码生成表
// @Tags         代码生成
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      int                   true  "代码生成表ID"
// @Param        request  body      dto.UpdateGenTableDTO true  "修改代码生成表"
// @Success      200      {object}  response.Response
// @Router       /table/{id} [put]
func (h *GenTableHandler) Update(c *gin.Context) {
	h.crud.HandleUpdate(c)
}

// BatchDelete 批量删除代码生成表。
//
// @Summary      批量删除代码生成表
// @Tags         代码生成
// @Produce      json
// @Security     BearerAuth
// @Param        ids   path      string  true  "逗号分隔的代码生成表ID"
// @Success      200   {object}  response.Response
// @Router       /table/{ids} [delete]
func (h *GenTableHandler) BatchDelete(c *gin.Context) {
	h.crud.HandleBatchDelete(c)
}

// DBTables 列出可导入的数据库表。
//
// @Summary      查询可导入数据库表
// @Tags         代码生成
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  response.Response
// @Router       /table/db [get]
func (h *GenTableHandler) DBTables(c *gin.Context) {
	list, err := h.svc.DBTables()
	if err != nil {
		response.FailBy(c, apperr.ErrQueryFailed)
		return
	}
	response.Success(c, list)
}

// ImportTables 导入选中的数据库表。
//
// @Summary      导入数据库表
// @Tags         代码生成
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      []string  true  "待导入的数据库表名列表"
// @Success      200      {object}  response.Response
// @Router       /table/import [post]
func (h *GenTableHandler) ImportTables(c *gin.Context) {
	var tableNames []string
	if err := c.ShouldBindJSON(&tableNames); err != nil {
		response.FailBy(c, apperr.ErrInvalidBody)
		return
	}
	if err := h.svc.ImportDbTables(c, tableNames); err != nil {
		response.FailBy(c, apperr.ErrCreateFailed, map[string]any{"error": err.Error()})
		return
	}
	response.Success(c, nil)
}

// SynchDbTable 同步数据库表结构。
//
// @Summary      同步数据库表结构
// @Tags         代码生成
// @Produce      json
// @Security     BearerAuth
// @Param        tableId  query     int  true  "代码生成表ID"
// @Success      200      {object}  response.Response
// @Router       /table/synch [get]
func (h *GenTableHandler) SynchDbTable(c *gin.Context) {
	tableID, err := strconv.ParseInt(c.Query("tableId"), 10, 64)
	if err != nil {
		response.FailBy(c, apperr.ErrInvalidId)
		return
	}
	if err := h.svc.SynchDbTable(c, tableID); err != nil {
		response.FailBy(c, apperr.ErrUpdateFailed, map[string]any{"error": err.Error()})
		return
	}
	response.Success(c, nil)
}

// CodePreview 代码预览。
//
// @Summary      代码预览
// @Tags         代码生成
// @Produce      json
// @Security     BearerAuth
// @Param        tableId  query     int  true  "代码生成表ID"
// @Success      200      {object}  response.Response
// @Router       /table/preview [get]
func (h *GenTableHandler) CodePreview(c *gin.Context) {
	tableID, err := strconv.ParseInt(c.Query("tableId"), 10, 64)
	if err != nil {
		response.FailBy(c, apperr.ErrInvalidId)
		return
	}

	data, err := h.svc.CodePreview(tableID)
	if err != nil {
		response.FailBy(c, apperr.ErrDataNotFound, map[string]any{"error": err.Error()})
		return
	}
	response.Success(c, data)
}

// Download 下载生成代码压缩包。
//
// @Summary      下载生成代码
// @Tags         代码生成
// @Produce      application/octet-stream
// @Security     BearerAuth
// @Param        tableIds  query     string  true  "逗号分隔的代码生成表ID"
// @Success      200       {file}    file
// @Router       /table/download [get]
func (h *GenTableHandler) Download(c *gin.Context) {
	tableIDs, err := utils.StrToNumberArray[int64](c.Query("tableIds"), ",")
	if err != nil {
		response.FailBy(c, apperr.ErrInvalidId)
		return
	}

	data, err := h.svc.DownloadCodeZip(tableIDs)
	if err != nil || data == nil {
		response.FailBy(c, apperr.ErrDataNotFound, map[string]any{"error": err.Error()})
		return
	}

	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Expose-Headers", "Content-Disposition")
	c.Header("Content-Disposition", `attachment; filename="aevo-gen-code.zip"`)
	c.Header("Content-Length", strconv.Itoa(len(data)))
	c.Data(http.StatusOK, "application/octet-stream", data)
}
