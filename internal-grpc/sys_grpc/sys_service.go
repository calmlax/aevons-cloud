package sys_grpc

import (
	"context"

	"github.com/calmlax/aevons-framework/config"
	"google.golang.org/grpc"
)

const (
	// SysServiceName 是系统服务的 gRPC 服务名。
	SysServiceName = "aevons.sys.v1.SysService"
	// GetConfValueByKeyMethodName 是按 key 查询系统配置值的方法名。
	GetConfValueByKeyMethodName = "GetConfValueByKey"
	// GetConfValueByKeyMethodFullName 是按 key 查询系统配置值的完整 gRPC 方法名。
	GetConfValueByKeyMethodFullName = "/" + SysServiceName + "/" + GetConfValueByKeyMethodName
	// GetDictDataCacheMethodName 是按字典类型查询字典缓存的方法名。
	GetDictDataCacheMethodName = "GetDictDataCache"
	// GetDictDataCacheMethodFullName 是按字典类型查询字典缓存的完整 gRPC 方法名。
	GetDictDataCacheMethodFullName = "/" + SysServiceName + "/" + GetDictDataCacheMethodName
)

// DictDataTranslation 定义字典翻译信息。
type DictDataTranslation struct {
	Label string `json:"label"`
	Tip   string `json:"tip"`
}

// DictDataItem 定义字典数据载荷。
type DictDataItem struct {
	ID           int64                          `json:"id,string,omitempty"`
	DictType     string                         `json:"dictType"`
	DictValue    string                         `json:"dictValue"`
	Status       int16                          `json:"status"`
	Sort         int                            `json:"sort"`
	TagType      string                         `json:"tagType"`
	TagClass     string                         `json:"tagClass"`
	LangCode     string                         `json:"langCode"`
	Label        string                         `json:"label"`
	Tip          string                         `json:"tip"`
	Translations map[string]DictDataTranslation `json:"translations,omitempty"`
}

// GetConfValueByKeyRequest 是按 key 查询系统配置值请求。
type GetConfValueByKeyRequest struct {
	Key string `json:"key"`
}

// GetConfValueByKeyResponse 是按 key 查询系统配置值响应。
type GetConfValueByKeyResponse struct {
	Value string `json:"value"`
}

// GetDictDataCacheRequest 是按字典类型查询字典缓存请求。
type GetDictDataCacheRequest struct {
	DictType string `json:"dictType"`
}

// GetDictDataCacheResponse 是按字典类型查询字典缓存响应。
type GetDictDataCacheResponse struct {
	Items []DictDataItem `json:"items"`
}

// Service 定义 sys-service 需要实现的系统查询 gRPC 服务接口。
type Service interface {
	GetConfValueByKey(ctx context.Context, req *GetConfValueByKeyRequest) (*GetConfValueByKeyResponse, error)
	GetDictDataCache(ctx context.Context, req *GetDictDataCacheRequest) (*GetDictDataCacheResponse, error)
}

// Client 是系统查询 gRPC 客户端。
type Client struct {
	conn *grpc.ClientConn
}

// NewClient 通过 Consul 服务发现创建系统查询 gRPC 客户端。
func NewClient(consulCfg config.ConsulConfig, opts ...grpc.DialOption) (*Client, error) {
	conn, err := newServiceConn(consulCfg, opts...)
	if err != nil {
		return nil, err
	}
	return &Client{conn: conn}, nil
}

// GetConfValueByKey 调用远端 sys-service 查询指定 key 的配置值。
func (c *Client) GetConfValueByKey(ctx context.Context, key string) (string, error) {
	resp := &GetConfValueByKeyResponse{}
	if err := invokeUnary(ctx, c.conn, GetConfValueByKeyMethodFullName, &GetConfValueByKeyRequest{Key: key}, resp); err != nil {
		return "", err
	}
	return resp.Value, nil
}

// GetDictDataCache 调用远端 sys-service 查询指定字典类型的缓存数据。
func (c *Client) GetDictDataCache(ctx context.Context, dictType string) ([]DictDataItem, error) {
	resp := &GetDictDataCacheResponse{}
	if err := invokeUnary(ctx, c.conn, GetDictDataCacheMethodFullName, &GetDictDataCacheRequest{DictType: dictType}, resp); err != nil {
		return nil, err
	}
	return resp.Items, nil
}

// Close 关闭底层 gRPC 连接。
func (c *Client) Close() error {
	if c == nil {
		return nil
	}
	return closeConn(c.conn)
}

// RegisterService 注册系统查询 gRPC 服务。
func RegisterService(registrar grpc.ServiceRegistrar, srv Service) {
	registrar.RegisterService(&grpc.ServiceDesc{
		ServiceName: SysServiceName,
		HandlerType: (*Service)(nil),
		Methods: []grpc.MethodDesc{
			{
				MethodName: GetConfValueByKeyMethodName,
				Handler:    getConfValueByKeyHandler,
			},
			{
				MethodName: GetDictDataCacheMethodName,
				Handler:    getDictDataCacheHandler,
			},
		},
		Streams:  []grpc.StreamDesc{},
		Metadata: "sys_service",
	}, srv)
}

func getConfValueByKeyHandler(
	srv any,
	ctx context.Context,
	dec func(any) error,
	interceptor grpc.UnaryServerInterceptor,
) (any, error) {
	in := new(GetConfValueByKeyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(Service).GetConfValueByKey(ctx, in)
	}

	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GetConfValueByKeyMethodFullName,
	}
	handler := func(ctx context.Context, req any) (any, error) {
		return srv.(Service).GetConfValueByKey(ctx, req.(*GetConfValueByKeyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func getDictDataCacheHandler(
	srv any,
	ctx context.Context,
	dec func(any) error,
	interceptor grpc.UnaryServerInterceptor,
) (any, error) {
	in := new(GetDictDataCacheRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(Service).GetDictDataCache(ctx, in)
	}

	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GetDictDataCacheMethodFullName,
	}
	handler := func(ctx context.Context, req any) (any, error) {
		return srv.(Service).GetDictDataCache(ctx, req.(*GetDictDataCacheRequest))
	}
	return interceptor(ctx, in, info, handler)
}
