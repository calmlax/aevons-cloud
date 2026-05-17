package service

import (
	"errors"
	"sys-service/internal/dto"
	"sys-service/internal/model"
	"sys-service/internal/repository"

	"github.com/calmlax/aevons-framework/core/base"
)

type LangService interface {
	base.BaseService[model.Lang, *dto.LangQuery]
	// 新增语言：若设为默认则清除其他默认
	CreateLang(d dto.CreateLangDTO) (*model.Lang, error)
	// 修改语言：若设为默认则清除其他默认
	UpdateLang(id int64, d dto.UpdateLangDTO) error
}

type langService struct {
	base.BaseService[model.Lang, *dto.LangQuery]
	repo repository.LangRepository
}

func NewLangService(repo repository.LangRepository) LangService {
	return &langService{
		BaseService: base.NewBaseService[model.Lang, *dto.LangQuery](repo),
		repo:        repo,
	}
}

func (s *langService) CreateLang(d dto.CreateLangDTO) (*model.Lang, error) {
	// 语言编码唯一校验
	exists, err := s.repo.ExistByField("lang_code", d.LangCode)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("语言编码已存在")
	}
	lang := &model.Lang{
		LangCode:  d.LangCode,
		LangName:  d.LangName,
		IsDefault: d.IsDefault,
		Sort:      d.Sort,
		Status:    d.Status,
		Remark:    d.Remark,
	}
	if err := s.repo.Create(lang); err != nil {
		return nil, err
	}
	// 设为默认时清除其他默认
	if d.IsDefault == 1 {
		_ = s.repo.ClearOtherDefault(lang.Id)
	}
	return lang, nil
}

func (s *langService) UpdateLang(id int64, d dto.UpdateLangDTO) error {
	// 语言编码唯一校验（排除自身）
	if d.LangCode != nil {
		exists, err := s.repo.ExistByFieldExcludeId("lang_code", *d.LangCode, id)
		if err != nil {
			return err
		}
		if exists {
			return errors.New("语言编码已存在")
		}
	}
	updates := map[string]any{}
	if d.LangCode != nil {
		updates["lang_code"] = *d.LangCode
	}
	if d.LangName != nil {
		updates["lang_name"] = *d.LangName
	}
	if d.IsDefault != nil {
		updates["is_default"] = *d.IsDefault
	}
	if d.Sort != nil {
		updates["sort"] = *d.Sort
	}
	if d.Status != nil {
		updates["status"] = *d.Status
	}
	if d.Remark != nil {
		updates["remark"] = *d.Remark
	}
	if _, err := s.repo.Update(id, updates); err != nil {
		return err
	}
	// 设为默认时清除其他默认
	if d.IsDefault != nil && *d.IsDefault == 1 {
		_ = s.repo.ClearOtherDefault(id)
	}
	return nil
}
