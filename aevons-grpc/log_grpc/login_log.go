package log_grpc

import (
	"context"
	"time"

	"github.com/calmlax/aevons-framework/grpcx"
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
	WriteLogin(ctx context.Context, entry LoginEntry) error
	Close() error
}

// LoginService 定义 log-service 需要实现的登录日志 gRPC 服务接口。
type LoginService interface {
	WriteLoginLog(ctx context.Context, req *WriteLoginRequest) (*WriteLoginResponse, error)
}

// NopLoginLogWriter 是登录日志空实现。
type NopLoginLogWriter struct{}

func (NopLoginLogWriter) WriteLogin(context.Context, LoginEntry) error { return nil }
func (NopLoginLogWriter) Close() error                                 { return nil }

// LoginLogClient 是登录日志 gRPC 客户端。
type LoginLogClient struct {
	conn *grpc.ClientConn
}

// NewLoginLogClient 创建登录日志 gRPC 客户端。
func NewLoginLogClient(target string, opts ...grpc.DialOption) (*LoginLogClient, error) {
	conn, err := grpcx.NewClientConn(target, opts...)
	if err != nil {
		return nil, err
	}
	return &LoginLogClient{conn: conn}, nil
}

// WriteLogin 调用远端 log-service 写入登录日志。
func (c *LoginLogClient) WriteLogin(ctx context.Context, entry LoginEntry) error {
	resp := &WriteLoginResponse{}
	return c.conn.Invoke(ctx, WriteLoginLogMethodFullName, &WriteLoginRequest{Entry: entry}, resp)
}

// Close 关闭底层 gRPC 连接。
func (c *LoginLogClient) Close() error {
	if c == nil {
		return nil
	}
	return grpcx.CloseClientConn(c.conn)
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
