/**
 * 岗位信息表 Service
 *
 * @author
 * @date 2026-04-09 01:37:54.290392924 +0000 UTC
 * @website https://www.aevons.com
 * @email 1073602@qq.com
 * @copyright 2025-2026 Aevons
 */
package service

import (
	"sys-service/internal/dto"
	"sys-service/internal/model"
	"sys-service/internal/repository"

	"github.com/calmlax/aevons-framework/core/base"
)

type PostService interface {
	base.BaseService[model.Post, *dto.PostQuery]
	ExistByPostKey(postKey string) (bool, error)
	ExistByPostKeyExcludeId(postKey string, excludeId int64) (bool, error)
}

type postService struct {
	base.BaseService[model.Post, *dto.PostQuery]
	repo repository.PostRepository
}

func NewPostService(repo repository.PostRepository) PostService {
	return &postService{
		BaseService: base.NewBaseService[model.Post, *dto.PostQuery](repo),
		repo:        repo,
	}
}

func (s *postService) ExistByPostKey(postKey string) (bool, error) {
	return s.repo.ExistByField("post_key", postKey)
}

func (s *postService) ExistByPostKeyExcludeId(postKey string, excludeId int64) (bool, error) {
	return s.repo.ExistByFieldExcludeId("post_key", postKey, excludeId)
}
