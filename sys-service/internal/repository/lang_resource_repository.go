package repository

import (
	"sys-service/internal/dto"
	"sys-service/internal/model"

	"github.com/calmlax/aevons-framework/core/base"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type LangResourceRepository interface {
	base.BaseRepository[model.LangResource]
	ExistByCompositeKey(langCode, namespace, resourceKey string, excludeId int64) (bool, error)
	GetByNamespaceAndKey(namespace, resourceKey string) ([]model.LangResource, error)
	UpsertTranslations(namespace, resourceKey string, items []dto.TranslationItem) error
	GetKeysByNamespace(namespace string) ([]string, error)
	// 按 namespace 去重分页查询 resourceKey（支持 key/content 搜索）
	PageKeys(namespace, resourceKey, content string, offset, limit int) ([]string, int64, error)
}

type langResourceRepository struct {
	base.BaseRepository[model.LangResource]
	db *gorm.DB
}

func NewLangResourceRepository(db *gorm.DB) LangResourceRepository {
	return &langResourceRepository{
		BaseRepository: base.NewBaseRepository[model.LangResource](db),
		db:             db,
	}
}

func (r *langResourceRepository) ExistByCompositeKey(langCode, namespace, resourceKey string, excludeId int64) (bool, error) {
	var count int64
	q := r.db.Model(&model.LangResource{}).
		Where("lang_code = ? AND namespace = ? AND resource_key = ?", langCode, namespace, resourceKey)
	if excludeId > 0 {
		q = q.Where("id != ?", excludeId)
	}
	if err := q.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *langResourceRepository) GetByNamespaceAndKey(namespace, resourceKey string) ([]model.LangResource, error) {
	var list []model.LangResource
	err := r.db.Where("namespace = ? AND resource_key = ?", namespace, resourceKey).Find(&list).Error
	return list, err
}

func (r *langResourceRepository) UpsertTranslations(namespace, resourceKey string, items []dto.TranslationItem) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, item := range items {
			record := model.LangResource{
				Namespace:   namespace,
				ResourceKey: resourceKey,
				LangCode:    item.LangCode,
				Content:     item.Content,
			}
			err := tx.Where("namespace = ? AND resource_key = ? AND lang_code = ?", namespace, resourceKey, item.LangCode).
				Assign(model.LangResource{Content: item.Content}).
				FirstOrCreate(&record).Error
			if err != nil {
				return err
			}
			// 如果已存在则更新 content
			if record.Content != item.Content {
				if err := tx.Model(&record).Update("content", item.Content).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
}

func (r *langResourceRepository) GetKeysByNamespace(namespace string) ([]string, error) {
	var keys []string
	err := r.db.Model(&model.LangResource{}).
		Where("namespace = ?", namespace).
		Distinct().Pluck("resource_key", &keys).Error
	return keys, err
}

func (r *langResourceRepository) PageKeys(namespace, resourceKey, content string, offset, limit int) ([]string, int64, error) {
	q := r.db.Model(&model.LangResource{}).Where("namespace = ?", namespace)
	if resourceKey != "" {
		q = q.Where("resource_key LIKE ?", "%"+resourceKey+"%")
	}
	if content != "" {
		q = q.Where("content LIKE ?", "%"+content+"%")
	}
	// 先统计去重后的 key 数量
	var total int64
	if err := q.Distinct("resource_key").Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var keys []string
	err := q.Distinct().Order("resource_key").
		Offset(offset).Limit(limit).
		Pluck("resource_key", &keys).Error
	return keys, total, err
}

// 确保 clause 包被使用（避免 unused import）
var _ = clause.OnConflict{}
