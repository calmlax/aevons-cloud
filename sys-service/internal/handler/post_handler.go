package handler

import (
	"sys-service/internal/dto"
	"sys-service/internal/model"
	"sys-service/internal/service"

	apperr "github.com/calmlax/aevons-framework/errors"
	"github.com/calmlax/aevons-framework/response"
	"github.com/calmlax/aevons-framework/utils"

	"github.com/calmlax/aevons-framework/core/base"

	"github.com/gin-gonic/gin"
)

type PostHandler struct {
	crud *base.BaseHandler[model.Post, *dto.PostQuery, dto.CreatePostDTO, dto.UpdatePostDTO]
	svc  service.PostService
}

func NewPostHandler(svc service.PostService) *PostHandler {
	return &PostHandler{
		crud: base.NewBaseHandler[model.Post, *dto.PostQuery, dto.CreatePostDTO, dto.UpdatePostDTO](svc),
		svc:  svc,
	}
}

// List 查询岗位列表。
//
// @Summary      查询岗位列表
// @Tags         岗位管理
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  response.Response
// @Router       /post/list [get]
func (h *PostHandler) List(c *gin.Context) {
	h.crud.HandleList(c)
}

// Page 分页查询岗位。
//
// @Summary      分页查询岗位
// @Tags         岗位管理
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  response.Response
// @Router       /post/page [get]
func (h *PostHandler) Page(c *gin.Context) {
	h.crud.HandlePage(c)
}

// Get 获取岗位详情。
//
// @Summary      获取岗位详情
// @Tags         岗位管理
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "岗位ID"
// @Success      200  {object}  response.Response
// @Router       /post/{id} [get]
func (h *PostHandler) Get(c *gin.Context) {
	h.crud.HandleGet(c)
}

// BatchDelete 批量删除岗位。
//
// @Summary      批量删除岗位
// @Tags         岗位管理
// @Produce      json
// @Security     BearerAuth
// @Param        ids   path      string  true  "逗号分隔的岗位ID"
// @Success      200   {object}  response.Response
// @Router       /post/{ids} [delete]
func (h *PostHandler) BatchDelete(c *gin.Context) {
	h.crud.HandleBatchDelete(c)
}

// Create 创建岗位，校验 postKey 唯一性
//
// @Summary      新增岗位
// @Tags         岗位管理
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      dto.CreatePostDTO  true  "新增岗位"
// @Success      200      {object}  response.Response
// @Router       /post [post]
func (h *PostHandler) Create(c *gin.Context) {
	var d dto.CreatePostDTO
	if err := c.ShouldBindJSON(&d); err != nil {
		response.FailBy(c, apperr.ErrInvalidBody)
		return
	}
	exists, err := h.svc.ExistByPostKey(d.PostKey)
	if err != nil {
		response.FailBy(c, apperr.ErrQueryFailed)
		return
	}
	if exists {
		response.Fail(c, 400, 4002, "err.sys.post.key_exists")
		return
	}
	var m model.Post
	utils.Copy(&m, d)
	if err := h.svc.Create(&m); err != nil {
		response.FailBy(c, apperr.ErrCreateFailed)
		return
	}
	response.Success(c, m)
}

// Update 更新岗位，校验 postKey 唯一性（排除自身）
//
// @Summary      修改岗位
// @Tags         岗位管理
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      int               true  "岗位ID"
// @Param        request  body      dto.UpdatePostDTO true  "修改岗位"
// @Success      200      {object}  response.Response
// @Router       /post/{id} [put]
func (h *PostHandler) Update(c *gin.Context) {
	id, ok := base.GetId(c)
	if !ok {
		return
	}
	var d dto.UpdatePostDTO
	if err := c.ShouldBindJSON(&d); err != nil {
		response.FailBy(c, apperr.ErrInvalidBody)
		return
	}
	if d.PostKey != nil {
		exists, err := h.svc.ExistByPostKeyExcludeId(*d.PostKey, id)
		if err != nil {
			response.FailBy(c, apperr.ErrQueryFailed)
			return
		}
		if exists {
			response.Fail(c, 400, 4002, "err.sys.post.key_exists")
			return
		}
	}
	mp := utils.StructToMapIgnoreNil(d)
	delete(mp, "id")
	if _, err := h.svc.Update(id, mp); err != nil {
		response.FailBy(c, apperr.ErrUpdateFailed)
		return
	}
	response.Success(c, id)
}
