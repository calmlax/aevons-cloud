package log_grpc

import (
	"context"
	"sync"
	"time"

	"github.com/calmlax/aevons-framework/config"
	"google.golang.org/grpc"
)

const (
	// LoginLogServiceName 是登录日志服务的 gRPC 服务名。
	LoginLogServiceName = "aevons.log.v1.LoginLogService"
	// WriteLoginLogMethodName 是写入登录日志的方法名。
	WriteLoginLogMethodName = "WriteLoginLog"
	// WriteLoginLogMethodFullName 是写入登录日志的完整 gRPC 方法名。
	WriteLoginLogMethodFullName = "/" + LoginLogServiceName + "/" + WriteLoginLogMethodName
	// GetLatestLoginLogMethodName 是查询最近登录日志的方法名。
	GetLatestLoginLogMethodName = "GetLatestLoginLog"
	// GetLatestLoginLogMethodFullName 是查询最近登录日志的完整 gRPC 方法名。
	GetLatestLoginLogMethodFullName = "/" + LoginLogServiceName + "/" + GetLatestLoginLogMethodName
)

// LoginEntry 定义通用登录日志载荷。
type LoginEntry struct {
	ID        int64     `json:"id,string,omitempty"`
	Username  string    `json:"username"`
	ClientID  string    `json:"clientId"`
	GrantType string    `json:"grantType"`
	OS        string    `json:"os"`
	Browser   string    `json:"browser"`
	IP        string    `json:"ip"`
	Location  string    `json:"location"`
	Status    int16     `json:"status"`
	Msg       string    `json:"msg"`
	LoginAt   time.Time `json:"loginAt"`
}

// WriteLoginRequest 是登录日志写入请求。
type WriteLoginRequest struct {
	Entry LoginEntry `json:"entry"`
}

// WriteLoginResponse 是登录日志写入响应。
type WriteLoginResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// GetLatestLoginLogRequest 是查询最近登录日志请求。
type GetLatestLoginLogRequest struct {
	Username string `json:"username"`
	Limit    int    `json:"limit"`
}

// GetLatestLoginLogResponse 是查询最近登录日志响应。
type GetLatestLoginLogResponse struct {
	Entries []LoginEntry `json:"entries"`
}

// LoginLogStore 定义登录日志完整访问接口。
type LoginLogStore interface {
	WriteLoginLog(ctx context.Context, entry LoginEntry) error
	GetLatestLoginLog(ctx context.Context, username string, limit int) ([]LoginEntry, error)
	Close() error
}

// LoginService 定义 log-service 需要实现的登录日志 gRPC 服务接口。
type LoginService interface {
	WriteLoginLog(ctx context.Context, req *WriteLoginRequest) (*WriteLoginResponse, error)
	GetLatestLoginLog(ctx context.Context, req *GetLatestLoginLogRequest) (*GetLatestLoginLogResponse, error)
}

// NopLoginLogStore 是登录日志空实现。
type NopLoginLogStore struct{}

func (NopLoginLogStore) WriteLoginLog(context.Context, LoginEntry) error { return nil }
func (NopLoginLogStore) GetLatestLoginLog(context.Context, string, int) ([]LoginEntry, error) {
	return nil, nil
}
func (NopLoginLogStore) Close() error { return nil }

// LoginLogClient 是登录日志 gRPC 客户端。
type LoginLogClient struct {
	conn *grpc.ClientConn
}

// ConsulLoginLogStore 是带 Consul 自动重连能力的登录日志写入器。
// 当 log-service 启动晚于业务服务时，不需要重启业务服务也能恢复写日志。
type ConsulLoginLogStore struct {
	consulCfg config.ConsulConfig
	mu        sync.Mutex
	client    *LoginLogClient
}

// NewLoginLogClient 通过 Consul 服务发现创建登录日志 gRPC 客户端。
// 默认使用 RegistryServiceName 作为日志中心服务注册名，业务侧无需再配置。
func NewLoginLogClient(consulCfg config.ConsulConfig, opts ...grpc.DialOption) (*LoginLogClient, error) {
	conn, err := newServiceConn(consulCfg, opts...)
	if err != nil {
		return nil, err
	}
	return &LoginLogClient{conn: conn}, nil
}

// NewConsulLoginLogStore 创建一个可按需连接、失败后自动重试的登录日志存储实现。
func NewConsulLoginLogStore(consulCfg config.ConsulConfig) LoginLogStore {
	return &ConsulLoginLogStore{consulCfg: consulCfg}
}

// WriteLoginLog 调用远端 log-service 写入登录日志。
func (c *LoginLogClient) WriteLoginLog(ctx context.Context, entry LoginEntry) error {
	resp := &WriteLoginResponse{}
	return invokeUnary(ctx, c.conn, WriteLoginLogMethodFullName, &WriteLoginRequest{Entry: entry}, resp)
}

// GetLatestLoginLog 调用远端 log-service 查询指定用户最近的登录日志。
func (c *LoginLogClient) GetLatestLoginLog(ctx context.Context, username string, limit int) ([]LoginEntry, error) {
	resp := &GetLatestLoginLogResponse{}
	err := invokeUnary(ctx, c.conn, GetLatestLoginLogMethodFullName, &GetLatestLoginLogRequest{
		Username: username,
		Limit:    limit,
	}, resp)
	if err != nil {
		return nil, err
	}
	return resp.Entries, nil
}

// Close 关闭底层 gRPC 连接。
func (c *LoginLogClient) Close() error {
	if c == nil {
		return nil
	}
	return closeConn(c.conn)
}

func (s *ConsulLoginLogStore) getClient() (*LoginLogClient, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.client != nil {
		return s.client, nil
	}

	client, err := NewLoginLogClient(s.consulCfg)
	if err != nil {
		return nil, err
	}
	s.client = client
	return client, nil
}

func (s *ConsulLoginLogStore) resetClient() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.client != nil {
		_ = s.client.Close()
	}
	s.client = nil
}

func (s *ConsulLoginLogStore) WriteLoginLog(ctx context.Context, entry LoginEntry) error {
	client, err := s.getClient()
	if err != nil {
		return err
	}

	if err := client.WriteLoginLog(ctx, entry); err == nil {
		return nil
	}

	s.resetClient()

	client, err = s.getClient()
	if err != nil {
		return err
	}
	return client.WriteLoginLog(ctx, entry)
}

func (s *ConsulLoginLogStore) GetLatestLoginLog(ctx context.Context, username string, limit int) ([]LoginEntry, error) {
	client, err := s.getClient()
	if err != nil {
		return nil, err
	}

	entries, err := client.GetLatestLoginLog(ctx, username, limit)
	if err == nil {
		return entries, nil
	}

	s.resetClient()

	client, err = s.getClient()
	if err != nil {
		return nil, err
	}
	return client.GetLatestLoginLog(ctx, username, limit)
}

func (s *ConsulLoginLogStore) Close() error {
	s.resetClient()
	return nil
}

// RegisterLoginService 注册登录日志 gRPC 服务。
func RegisterLoginService(registrar grpc.ServiceRegistrar, srv LoginService) {
	registrar.RegisterService(&grpc.ServiceDesc{
		ServiceName: LoginLogServiceName,
		HandlerType: (*LoginService)(nil),
		Methods: []grpc.MethodDesc{
			{
				MethodName: WriteLoginLogMethodName,
				Handler:    writeLoginLogHandler,
			},
			{
				MethodName: GetLatestLoginLogMethodName,
				Handler:    getLatestLoginLogHandler,
			},
		},
		Streams:  []grpc.StreamDesc{},
		Metadata: "login_log",
	}, srv)
}

func writeLoginLogHandler(
	srv any,
	ctx context.Context,
	dec func(any) error,
	interceptor grpc.UnaryServerInterceptor,
) (any, error) {
	in := new(WriteLoginRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LoginService).WriteLoginLog(ctx, in)
	}

	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: WriteLoginLogMethodFullName,
	}
	handler := func(ctx context.Context, req any) (any, error) {
		return srv.(LoginService).WriteLoginLog(ctx, req.(*WriteLoginRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func getLatestLoginLogHandler(
	srv any,
	ctx context.Context,
	dec func(any) error,
	interceptor grpc.UnaryServerInterceptor,
) (any, error) {
	in := new(GetLatestLoginLogRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LoginService).GetLatestLoginLog(ctx, in)
	}

	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GetLatestLoginLogMethodFullName,
	}
	handler := func(ctx context.Context, req any) (any, error) {
		return srv.(LoginService).GetLatestLoginLog(ctx, req.(*GetLatestLoginLogRequest))
	}
	return interceptor(ctx, in, info, handler)
}
