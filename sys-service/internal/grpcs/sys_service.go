package grpcs

import (
	"context"
	"sys-service/internal/repository"
	"sys-service/internal/service"

	"internal-grpc/sys_grpc"

	"gorm.io/gorm"
)

// SysServiceServer 实现 sys-service 的系统查询 gRPC 服务。
type SysServiceServer struct {
	confSvc service.ConfService
	dictSvc service.DictDataService
}

// NewSysServiceServer 创建系统查询 gRPC 服务实例。
func NewSysServiceServer(db *gorm.DB) *SysServiceServer {
	confRepo := repository.NewConfRepository(db)
	dictRepo := repository.NewDictDataRepository(db)
	return &SysServiceServer{
		confSvc: service.NewConfService(confRepo),
		dictSvc: service.NewDictDataService(dictRepo),
	}
}

// GetConfValueByKey 查询指定 key 的系统配置值。
func (s *SysServiceServer) GetConfValueByKey(ctx context.Context, req *sys_grpc.GetConfValueByKeyRequest) (*sys_grpc.GetConfValueByKeyResponse, error) {
	value, err := s.confSvc.GetConfValueByKey(req.Key)
	if err != nil {
		return nil, err
	}
	return &sys_grpc.GetConfValueByKeyResponse{Value: value}, nil
}

// GetDictDataCache 查询指定字典类型的缓存数据。
func (s *SysServiceServer) GetDictDataCache(ctx context.Context, req *sys_grpc.GetDictDataCacheRequest) (*sys_grpc.GetDictDataCacheResponse, error) {
	items, err := s.dictSvc.GetDictDataCache(req.DictType)
	if err != nil {
		return nil, err
	}

	respItems := make([]sys_grpc.DictDataItem, len(items))
	for i, item := range items {
		translations := make(map[string]sys_grpc.DictDataTranslation, len(item.Translations))
		for langCode, trans := range item.Translations {
			translations[langCode] = sys_grpc.DictDataTranslation{
				Label: trans.Label,
				Tip:   trans.Tip,
			}
		}
		respItems[i] = sys_grpc.DictDataItem{
			ID:           item.Id,
			DictType:     item.DictType,
			DictValue:    item.DictValue,
			Status:       item.Status,
			Sort:         item.Sort,
			TagType:      item.TagType,
			TagClass:     item.TagClass,
			LangCode:     item.LangCode,
			Label:        item.Label,
			Tip:          item.Tip,
			Translations: translations,
		}
	}

	return &sys_grpc.GetDictDataCacheResponse{Items: respItems}, nil
}
