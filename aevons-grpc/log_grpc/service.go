package log_grpc

import (
	"context"
	"errors"
	"time"

	"github.com/calmlax/aevons-framework/grpcx"

	"google.golang.org/grpc"
)

const (
	ServiceName              = "aevons.log.v1.OperLogService"
	WriteMethodName          = "WriteOperLog"
	WriteMethodFullName      = "/" + ServiceName + "/" + WriteMethodName
	LoginServiceName         = "aevons.log.v1.LoginLogService"
	WriteLoginMethodName     = "WriteLoginLog"
	WriteLoginMethodFullName = "/" + LoginServiceName + "/" + WriteLoginMethodName
)

var ErrEmptyTarget = errors.New("operlog grpc target is empty")

// Entry 定义通用操作日志载荷，供各业务服务通过 gRPC 同步写入 log-server。
type Entry struct {
	Module      string    `json:"module"`
	Type        string    `json:"type"`
	Description string    `json:"description"`
	Method      string    `json:"method"`
	URL         string    `json:"url"`
	IP          string    `json:"ip"`
	Location    string    `json:"location"`
	Payload     string    `json:"payload"`
	Result      string    `json:"result"`
	Device      string    `json:"device"`
	OS          string    `json:"os"`
	Browser     string    `json:"browser"`
	Status      int16     `json:"status"`
	Error       string    `json:"error"`
	TimeMS      int64     `json:"timeMs"`
	UserID      int64     `json:"userId"`
	Username    string    `json:"username"`
	OperAt      time.Time `json:"operAt"`
}

type WriteRequest struct {
	Entry Entry `json:"entry"`
}

type WriteResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// LoginEntry 定义通用登录日志载荷，供各业务服务通过 gRPC 同步写入 log-server。
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

type WriteLoginRequest struct {
	Entry LoginEntry `json:"entry"`
}

type WriteLoginResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// Writer 定义操作日志写入器，供中间件统一调用。
type Writer interface {
	Write(ctx context.Context, entry Entry) error
	Close() error
}

// LoginWriter 定义登录日志写入器。
type LoginWriter interface {
	WriteLogin(ctx context.Context, entry LoginEntry) error
	Close() error
}

// Service 定义 log-server 需要实现的 gRPC 服务接口。
type Service interface {
	WriteOperLog(ctx context.Context, req *WriteRequest) (*WriteResponse, error)
}

// LoginService 定义 log-server 需要实现的登录日志 gRPC 服务接口。
type LoginService interface {
	WriteLoginLog(ctx context.Context, req *WriteLoginRequest) (*WriteLoginResponse, error)
}

type NopWriter struct{}

func (NopWriter) Write(context.Context, Entry) error { return nil }
func (NopWriter) Close() error                       { return nil }

type NopLoginWriter struct{}

func (NopLoginWriter) WriteLogin(context.Context, LoginEntry) error { return nil }
func (NopLoginWriter) Close() error                                 { return nil }

type Client struct {
	conn *grpc.ClientConn
}

func NewClient(target string, opts ...grpc.DialOption) (*Client, error) {
	if target == "" {
		return nil, ErrEmptyTarget
	}

	conn, err := grpcx.NewClientConn(target, opts...)
	if err != nil {
		return nil, err
	}

	return &Client{conn: conn}, nil
}

func (c *Client) Write(ctx context.Context, entry Entry) error {
	resp := &WriteResponse{}
	return c.conn.Invoke(ctx, WriteMethodFullName, &WriteRequest{Entry: entry}, resp)
}

func (c *Client) Close() error {
	if c == nil || c.conn == nil {
		return nil
	}
	return c.conn.Close()
}

type LoginClient struct {
	conn *grpc.ClientConn
}

func NewLoginClient(target string, opts ...grpc.DialOption) (*LoginClient, error) {
	if target == "" {
		return nil, ErrEmptyTarget
	}

	conn, err := grpcx.NewClientConn(target, opts...)
	if err != nil {
		return nil, err
	}

	return &LoginClient{conn: conn}, nil
}

func (c *LoginClient) WriteLogin(ctx context.Context, entry LoginEntry) error {
	resp := &WriteLoginResponse{}
	return c.conn.Invoke(ctx, WriteLoginMethodFullName, &WriteLoginRequest{Entry: entry}, resp)
}

func (c *LoginClient) Close() error {
	if c == nil || c.conn == nil {
		return nil
	}
	return c.conn.Close()
}

func RegisterService(registrar grpc.ServiceRegistrar, srv Service) {
	registrar.RegisterService(&grpc.ServiceDesc{
		ServiceName: ServiceName,
		HandlerType: (*Service)(nil),
		Methods: []grpc.MethodDesc{
			{
				MethodName: WriteMethodName,
				Handler:    writeOperLogHandler,
			},
		},
		Streams:  []grpc.StreamDesc{},
		Metadata: "operlog",
	}, srv)
}

func RegisterLoginService(registrar grpc.ServiceRegistrar, srv LoginService) {
	registrar.RegisterService(&grpc.ServiceDesc{
		ServiceName: LoginServiceName,
		HandlerType: (*LoginService)(nil),
		Methods: []grpc.MethodDesc{
			{
				MethodName: WriteLoginMethodName,
				Handler:    writeLoginLogHandler,
			},
		},
		Streams:  []grpc.StreamDesc{},
		Metadata: "login_log",
	}, srv)
}

func writeOperLogHandler(
	srv any,
	ctx context.Context,
	dec func(any) error,
	interceptor grpc.UnaryServerInterceptor,
) (any, error) {
	in := new(WriteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(Service).WriteOperLog(ctx, in)
	}

	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: WriteMethodFullName,
	}
	handler := func(ctx context.Context, req any) (any, error) {
		return srv.(Service).WriteOperLog(ctx, req.(*WriteRequest))
	}
	return interceptor(ctx, in, info, handler)
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
		FullMethod: WriteLoginMethodFullName,
	}
	handler := func(ctx context.Context, req any) (any, error) {
		return srv.(LoginService).WriteLoginLog(ctx, req.(*WriteLoginRequest))
	}
	return interceptor(ctx, in, info, handler)
}
