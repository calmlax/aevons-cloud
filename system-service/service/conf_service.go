/**
 * 参数配置表 Service
 *
 * @author
 * @date 2026-04-09 00:38:25.504785055 +0000 UTC
 * @website https://www.aevons.com
 * @email 1073602@qq.com
 * @copyright 2025-2026 Aevons
 */
package service

import (
	"context"
	"errors"
	"fmt"
	"system-service/dto"
	"system-service/model"
	"system-service/repository"
	"time"

	"github.com/calmlax/aevons-framework/consts"
	"github.com/calmlax/aevons-framework/redis"
	"github.com/calmlax/aevons-framework/utils"

	"github.com/calmlax/aevons-framework/core/base"
	errdef "github.com/calmlax/aevons-framework/err"

	"gorm.io/gorm"
)

// 继承BaseService
type ConfService interface {
	base.BaseService[model.Conf, *dto.ConfQuery]
	GetConfByKey(key string) (*model.Conf, error)
	GetConfValueByKey(key string) (string, error)
	RefreshCache() error
	DecryptIfNeeded(value string, isSecret bool) string
	EncryptIfNeeded(value string, isSecret bool) string
}

type confService struct {
	base.BaseService[model.Conf, *dto.ConfQuery]
	repo repository.ConfRepository
}

func NewConfService(repo repository.ConfRepository) ConfService {
	baseSrv := base.NewBaseService[model.Conf, *dto.ConfQuery](repo)
	return &confService{
		BaseService: baseSrv,
		repo:        repo,
	}
}

// DecryptIfNeeded 在配置数据需要保密 (IsSecret) 的前提下自动对其尝试 AES 解密。
func (s *confService) DecryptIfNeeded(value string, isSecret bool) string {
	if isSecret {
		if decryptedVal, err := utils.DecryptAES(value, utils.DefaultSymmetricKey); err == nil {
			value = decryptedVal
		} else {
			fmt.Print(err)
		}
	}
	return value
}

// EncryptIfNeeded 在配置数据需要保密 (IsSecret) 的前提下自动对其尝试 AES 加密。
func (s *confService) EncryptIfNeeded(value string, isSecret bool) string {
	if isSecret {
		if encryptedVal, err := utils.EncryptAES(value, utils.DefaultSymmetricKey); err == nil {
			value = encryptedVal
		}
	}
	return value
}

// GetConfByKey 获取指定键的系统配置记录，包含 Redis 级缓存并自动拦截解析密文。
func (s *confService) GetConfByKey(key string) (*model.Conf, error) {
	cacheKey := consts.ConfCacheKeyPrefix + key
	var conf model.Conf

	ctx := context.Background()
	// 1. 优先尝试从 Redis 缓存获取 JSON 反序列化数据
	if err := redis.GetJSON(ctx, cacheKey, &conf); err == nil {
		// 缓存命中则按需解密并返回
		conf.ConfValue = s.DecryptIfNeeded(conf.ConfValue, conf.IsSecret == 1)
		return &conf, nil
	}

	// 2. 缓存未命中时查底层 DB
	dbConf, err := s.repo.GetByField("conf_key", key)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errdef.ErrNotFound
		}
		return nil, err
	}
	conf = *dbConf
	// 3. 将密文态数据库查出的结果直接存入 Redis 以供复用，过期时间置为 24 小时
	_ = redis.SetJSON(ctx, cacheKey, conf, 24*time.Hour)

	conf.ConfValue = s.DecryptIfNeeded(conf.ConfValue, conf.IsSecret == 1)
	return &conf, nil
}

// GetConfValueByKey 获取指定键的系统配置记录
func (s *confService) GetConfValueByKey(key string) (string, error) {
	conf, err := s.GetConfByKey(key)
	if err != nil {
		return "", err
	}
	return conf.ConfValue, nil
}

// RefreshCache 扫描并删除所有的系统配置 Redis 缓存（基于前缀）。
func (s *confService) RefreshCache() error {
	ctx := context.Background()
	keys, err := redis.Client.Keys(ctx, consts.ConfCacheKeyPrefix+"*").Result()
	if err != nil {
		return err
	}
	if len(keys) > 0 {
		return redis.Del(ctx, keys...)
	}
	return nil
}
