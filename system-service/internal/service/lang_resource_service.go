package service

import (
	"errors"
	"system-service/internal/dto"
	"system-service/internal/model"
	"system-service/internal/repository"

	"github.com/calmlax/aevons-framework/core/base"
)

type LangResourceService interface {
	base.BaseService[model.LangResource, *dto.LangResourceQuery]
	CreateResource(d dto.CreateLangResourceDTO) (*model.LangResource, error)
	UpdateResource(id int64, d dto.UpdateLangResourceDTO) error
	// 按 namespace+resourceKey 获取所有语言翻译
	GetTranslations(namespace, resourceKey string) ([]model.LangResource, error)
	// 批量保存翻译（upsert：有则更新，无则新增）
	SaveTranslations(namespace, resourceKey string, items []dto.TranslationItem) error
	// 按 namespace 获取去重的 resourceKey 列表
	GetKeysByNamespace(namespace string) ([]string, error)
	// 去重分页查询 resourceKey
	PageKeys(namespace, resourceKey, content string, offset, limit int) ([]string, int64, error)
}

type langResourceService struct {
	base.BaseService[model.LangResource, *dto.LangResourceQuery]
	repo repository.LangResourceRepository
}

func NewLangResourceService(repo repository.LangResourceRepository) LangResourceService {
	return &langResourceService{
		BaseService: base.NewBaseService[model.LangResource, *dto.LangResourceQuery](repo),
		repo:        repo,
	}
}

func (s *langResourceService) CreateResource(d dto.CreateLangResourceDTO) (*model.LangResource, error) {
	exists, err := s.repo.ExistByCompositeKey(d.LangCode, d.Namespace, d.ResourceKey, 0)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("语言编码+命名空间+资源Key已存在")
	}
	r := &model.LangResource{
		ResourceKey: d.ResourceKey,
		Namespace:   d.Namespace,
		LangCode:    d.LangCode,
		Content:     d.Content,
		Status:      d.Status,
	}
	if err := s.repo.Create(r); err != nil {
		return nil, err
	}
	return r, nil
}

func (s *langResourceService) UpdateResource(id int64, d dto.UpdateLangResourceDTO) error {
	langCode := ""
	namespace := ""
	resourceKey := ""
	if d.LangCode != nil {
		langCode = *d.LangCode
	}
	if d.Namespace != nil {
		namespace = *d.Namespace
	}
	if d.ResourceKey != nil {
		resourceKey = *d.ResourceKey
	}
	if langCode != "" && namespace != "" && resourceKey != "" {
		exists, err := s.repo.ExistByCompositeKey(langCode, namespace, resourceKey, id)
		if err != nil {
			return err
		}
		if exists {
			return errors.New("语言编码+命名空间+资源Key已存在")
		}
	}
	updates := map[string]any{}
	if d.ResourceKey != nil {
		updates["resource_key"] = *d.ResourceKey
	}
	if d.Namespace != nil {
		updates["namespace"] = *d.Namespace
	}
	if d.LangCode != nil {
		updates["lang_code"] = *d.LangCode
	}
	if d.Content != nil {
		updates["content"] = *d.Content
	}
	if d.Status != nil {
		updates["status"] = *d.Status
	}
	_, err := s.repo.Update(id, updates)
	return err
}

func (s *langResourceService) GetTranslations(namespace, resourceKey string) ([]model.LangResource, error) {
	return s.repo.GetByNamespaceAndKey(namespace, resourceKey)
}

func (s *langResourceService) SaveTranslations(namespace, resourceKey string, items []dto.TranslationItem) error {
	return s.repo.UpsertTranslations(namespace, resourceKey, items)
}

func (s *langResourceService) GetKeysByNamespace(namespace string) ([]string, error) {
	return s.repo.GetKeysByNamespace(namespace)
}

func (s *langResourceService) PageKeys(namespace, resourceKey, content string, offset, limit int) ([]string, int64, error) {
	return s.repo.PageKeys(namespace, resourceKey, content, offset, limit)
}
