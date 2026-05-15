package repository

import (
	"system-service/model"

	"github.com/calmlax/aevons-framework/core/base"

	"gorm.io/gorm"
)

type LangRepository interface {
	base.BaseRepository[model.Lang]
	// 清除其他语言的默认标记
	ClearOtherDefault(excludeId int64) error
}

type langRepository struct {
	base.BaseRepository[model.Lang]
	db *gorm.DB
}

func NewLangRepository(db *gorm.DB) LangRepository {
	return &langRepository{
		BaseRepository: base.NewBaseRepository[model.Lang](db),
		db:             db,
	}
}

func (r *langRepository) ClearOtherDefault(excludeId int64) error {
	return r.db.Model(&model.Lang{}).
		Where("id != ? AND is_default = 1", excludeId).
		Update("is_default", 0).Error
}
