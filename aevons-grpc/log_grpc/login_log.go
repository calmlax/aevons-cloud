package log_grpc

import (
	"context"
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

// LoginLogWriter 定义登录日志写入器。
type LoginLogWriter interface {
	WriteLoginLog(ctx context.Context, entry LoginEntry) error
	Close() error
}

// LoginService 定义 log-service 需要实现的登录日志 gRPC 服务接口。
type LoginService interface {
	WriteLoginLog(ctx context.Context, req *WriteLoginRequest) (*WriteLoginResponse, error)
}

// NopLoginLogWriter 是登录日志空实现。
type NopLoginLogWriter struct{}

func (NopLoginLogWriter) WriteLoginLog(context.Context, LoginEntry) error { return nil }
func (NopLoginLogWriter) Close() error                                    { return nil }

// LoginLogClient 是登录日志 gRPC 客户端。
type LoginLogClient struct {
	conn *grpc.ClientConn
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

// WriteLoginLog 调用远端 log-service 写入登录日志。
func (c *LoginLogClient) WriteLoginLog(ctx context.Context, entry LoginEntry) error {
	resp := &WriteLoginResponse{}
	return invokeUnary(ctx, c.conn, WriteLoginLogMethodFullName, &WriteLoginRequest{Entry: entry}, resp)
}

// Close 关闭底层 gRPC 连接。
func (c *LoginLogClient) Close() error {
	if c == nil {
		return nil
	}
	return closeConn(c.conn)
}

// RegisterLoginService 注册登录日志 gRPC 服务。
func RegisterLoginService(registrar grpc.ServiceRegistrar, srv LoginService) {
	registerUnaryService(registrar, srv, LoginLogServiceName, (*LoginService)(nil), WriteLoginLogMethodName, writeLoginLogHandler, "login_log")
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
