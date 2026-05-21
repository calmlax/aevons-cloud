package excel

import (
	"context"
	"sys-service/internal/service"

	"github.com/calmlax/aevons-framework/core/base"
	fwexcel "github.com/calmlax/aevons-framework/excel"
	"github.com/gin-gonic/gin"
)

// DictProviderBuilder 定义 sys-service 的 Excel 字典提供器构建能力。
// Handler 只依赖这个小接口，不直接感知 DictDataService 等更重的业务服务。
type DictProviderBuilder interface {
	Build(c *gin.Context) fwexcel.DictProvider
}

type dictProviderBuilder struct {
	dictSvc service.DictDataService
}

// NewDictProviderBuilder 创建一个基于字典服务的 Excel 字典提供器构建器。
func NewDictProviderBuilder(dictSvc service.DictDataService) DictProviderBuilder {
	return &dictProviderBuilder{dictSvc: dictSvc}
}

// Build 基于当前请求语言构造 sys-service 通用 Excel 字典提供器。
func (b *dictProviderBuilder) Build(c *gin.Context) fwexcel.DictProvider {
	langCode := base.GetLanguage(c)
	return func(ctx context.Context, dictKey string) ([]fwexcel.DictItem, error) {
		list, err := b.dictSvc.GetDictDataCache(dictKey)
		if err != nil {
			return nil, err
		}

		items := make([]fwexcel.DictItem, 0, len(list))
		for _, item := range list {
			if item.LangCode == langCode {
				items = append(items, fwexcel.DictItem{
					Value: item.DictValue,
					Label: item.Label,
				})
			}
		}
		return items, nil
	}
}
