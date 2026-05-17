/**
 * 字典数据表 Service
 *
 * @author
 * @date 2026-04-09 01:08:50.443674979 +0000 UTC
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
	"time"

	consts "github.com/calmlax/aevons-framework/consts"
	pkgredis "github.com/calmlax/aevons-framework/redis"
	"github.com/calmlax/aevons-framework/utils"

	"github.com/calmlax/aevons-framework/core/base"

	"golang.org/x/sync/singleflight"
)

// 继承BaseService
type DictDataService interface {
	base.BaseService[model.DictData, *dto.DictDataQuery]

	CreateDictData(ctx context.Context, dictDataDto dto.CreateDictDataDTO) (int64, error)

	UpdateDictData(ctx context.Context, dictDataDto dto.UpdateDictDataDTO) error

	ListByDictType(dictType string, langCode string) ([]dto.DictDataDTO, error)

	GetDetail(id int64) (dto.DictDataDTO, error)

	GetDictDataCache(dictType string) ([]dto.DictDataDTO, error)

	RefreshCache() error

	DeleteByIds(ctx context.Context, ids []int64) error

	UpdateSort(ctx context.Context, items []dto.SortItemDTO) error
}

type dictDataService struct {
	base.BaseService[model.DictData, *dto.DictDataQuery]
	repo  repository.DictDataRepository
	group singleflight.Group
}

func NewDictDataService(repo repository.DictDataRepository) DictDataService {
	baseSrv := base.NewBaseService[model.DictData, *dto.DictDataQuery](repo)
	return &dictDataService{
		BaseService: baseSrv,
		repo:        repo,
	}
}

func (s *dictDataService) CreateDictData(ctx context.Context, dictDataDto dto.CreateDictDataDTO) (int64, error) {
	count, err := s.repo.Count(&dto.DictDataQuery{
		DictType: &dictDataDto.DictType,
	})
	if err != nil {
		return 0, err
	}
	dictDataDto.Sort = int(count) + 1
	return s.repo.CreateDictData(ctx, dictDataDto)
}

func (s *dictDataService) UpdateDictData(ctx context.Context, dictDataDto dto.UpdateDictDataDTO) error {

	return s.repo.UpdateDictData(ctx, dictDataDto)
}
func (s *dictDataService) ListByDictType(dictType string, langCode string) ([]dto.DictDataDTO, error) {
	if utils.IsEmpty(dictType) {
		return []dto.DictDataDTO{}, nil
	}
	return s.repo.ListByDictTypeAndLangCode(dictType, langCode)
}

// GetDictData 获取字典数据（带缓存、防击穿、防穿透）
func (s *dictDataService) GetDictDataCache(dictType string) ([]dto.DictDataDTO, error) {
	// 空参数直接返回空数组
	if utils.IsEmpty(dictType) {
		return []dto.DictDataDTO{}, nil
	}

	ctx := context.Background()
	cacheKey := consts.DictCacheKeyPrefix + dictType

	// 1. 先查缓存
	var dictDataList []dto.DictDataDTO
	if err := pkgredis.GetJSON(ctx, cacheKey, &dictDataList); err == nil {
		return dictDataList, nil
	}

	// 2. 【防击穿核心】同一 key 只允许一个请求访问DB
	result, err, _ := s.group.Do(cacheKey, func() (interface{}, error) {
		// 双重检查缓存，避免高并发重复请求
		if err := pkgredis.GetJSON(ctx, cacheKey, &dictDataList); err == nil {
			return dictDataList, nil
		}

		// 3. 查询数据库
		list, err := s.repo.ListByDictType(dictType)
		if err != nil {
			return []dto.DictDataDTO{}, err
		}

		// 4. 空数据缓存（防穿透）
		if len(list) == 0 {
			_ = pkgredis.SetJSON(ctx, cacheKey, []dto.DictDataDTO{}, 5*time.Minute)
			return list, nil
		}

		// 5. 正常数据缓存24小时
		_ = pkgredis.SetJSON(ctx, cacheKey, list, 24*time.Hour)
		return list, nil
	})

	if err != nil {
		return []dto.DictDataDTO{}, err
	}

	return result.([]dto.DictDataDTO), nil
}

// RefreshCache 扫描并删除所有的系统配置 Redis 缓存（基于前缀）。
func (s *dictDataService) RefreshCache() error {
	ctx := context.Background()
	keys, err := pkgredis.Client.Keys(ctx, consts.DictCacheKeyPrefix+"*").Result()
	if err != nil {
		return err
	}
	if len(keys) > 0 {
		return pkgredis.Del(ctx, keys...)
	}
	return nil
}

// GetDetail 获取字典数据详情
func (s *dictDataService) GetDetail(id int64) (dto.DictDataDTO, error) {

	dd, err := s.repo.GetById(id)
	if err != nil {
		return dto.DictDataDTO{}, err
	}
	tl := make(map[string]dto.DictDataTlDTO)

	tls, err := s.repo.GetTranslationsById(id)
	if err != nil {
		return dto.DictDataDTO{}, err
	}

	for _, item := range tls {
		tl[item.LangCode] = dto.DictDataTlDTO{
			Label: item.Label,
			Tip:   item.Tip,
		}
	}
	return dto.DictDataDTO{
		Id:           dd.Id,
		DictType:     dd.DictType,
		DictValue:    dd.DictValue,
		Status:       dd.Status,
		Sort:         dd.Sort,
		TagType:      dd.TagType,
		TagClass:     dd.TagClass,
		Translations: tl,
	}, nil
}

func (s *dictDataService) DeleteByIds(ctx context.Context, ids []int64) error {
	return s.repo.DeleteByIds(ctx, ids)
}

func (s *dictDataService) UpdateSort(ctx context.Context, items []dto.SortItemDTO) error {
	return s.repo.UpdateSort(ctx, items)
}
