/**
 * 菜单权限表 Service
 *
 * @author
 * @date 2026-04-09 01:17:32.070579992 +0000 UTC
 * @website https://www.aevons.com
 * @email 1073602@qq.com
 * @copyright 2025-2026 Aevons
 */
package service

import (
	"context"
	"system-service/internal/dto"
	"system-service/internal/model"
	"system-service/internal/repository"

	"github.com/calmlax/aevons-framework/utils"

	"github.com/calmlax/aevons-framework/core/base"
)

// 继承BaseService
type MenuService interface {
	base.BaseService[model.Menu, *dto.MenuQuery]

	CreateMenu(ctx context.Context, dto dto.CreateMenuDTO) (int64, error)

	UpdateMenu(ctx context.Context, dto dto.UpdateMenuDTO) error

	ListByLangCode(langCode string) ([]dto.MenuDTO, error)

	GetDetail(id int64) (dto.MenuDTO, error)

	DeleteByIds(ctx context.Context, ids []int64) error
}

type menuService struct {
	base.BaseService[model.Menu, *dto.MenuQuery]
	repo repository.MenuRepository
}

func NewMenuService(repo repository.MenuRepository) MenuService {
	baseSrv := base.NewBaseService[model.Menu, *dto.MenuQuery](repo)
	return &menuService{
		BaseService: baseSrv,
		repo:        repo,
	}
}

func (s *menuService) CreateMenu(ctx context.Context, dto dto.CreateMenuDTO) (int64, error) {

	return s.repo.CreateMenu(ctx, dto)
}

func (s *menuService) UpdateMenu(ctx context.Context, dto dto.UpdateMenuDTO) error {

	return s.repo.UpdateMenu(ctx, dto)
}

func (s *menuService) ListByLangCode(langCode string) ([]dto.MenuDTO, error) {

	return s.repo.ListByLangCode(langCode)
}

func (s *menuService) GetDetail(id int64) (dto.MenuDTO, error) {

	menu, err := s.repo.GetById(id)
	if err != nil {
		return dto.MenuDTO{}, err
	}
	tl := make(map[string]string)

	tls, err := s.repo.GetTranslationsById(id)
	if err != nil {
		return dto.MenuDTO{}, err
	}

	for _, item := range tls {
		tl[item.LangCode] = item.Title
	}
	var menuDto dto.MenuDTO
	utils.Copy(&menuDto, menu)
	menuDto.Translations = tl
	return menuDto, nil
}

func (s *menuService) DeleteByIds(ctx context.Context, ids []int64) error {
	return s.repo.DeleteByIds(ctx, ids)
}
