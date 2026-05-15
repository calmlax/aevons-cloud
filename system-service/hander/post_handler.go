package handler

import (
	"system-service/dto"
	"system-service/model"
	"system-service/service"

	apperr "github.com/calmlax/aevons-framework/errors"
	"github.com/calmlax/aevons-framework/response"
	"github.com/calmlax/aevons-framework/utils"

	"github.com/calmlax/aevons-framework/core/base"

	"github.com/gin-gonic/gin"
)

type PostHandler struct {
	*base.BaseHandler[
		model.Post,
		*dto.PostQuery,
		dto.CreatePostDTO,
		dto.UpdatePostDTO,
	]
	svc service.PostService
}

func NewPostHandler(svc service.PostService) *PostHandler {
	return &PostHandler{
		BaseHandler: base.NewBaseHandler[
			model.Post,
			*dto.PostQuery,
			dto.CreatePostDTO,
			dto.UpdatePostDTO,
		](svc),
		svc: svc,
	}
}

// Create 创建岗位，校验 postKey 唯一性
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
func (h *PostHandler) Update(c *gin.Context) {
	id, ok := h.GetId(c)
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
