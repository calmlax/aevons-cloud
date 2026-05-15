package log_grpc

import (
	"context"
	"time"

	"github.com/calmlax/aevons-framework/grpcx"
	"google.golang.org/grpc"
)

const (
	// OperLogServiceName 是操作日志服务的 gRPC 服务名。
	OperLogServiceName = "aevons.log.v1.OperLogService"
	// WriteOperLogMethodName 是写入操作日志的方法名。
	WriteOperLogMethodName = "WriteOperLog"
	// WriteOperLogMethodFullName 是写入操作日志的完整 gRPC 方法名。
	WriteOperLogMethodFullName = "/" + OperLogServiceName + "/" + WriteOperLogMethodName
)

// Entry 定义通用操作日志载荷，供各业务服务同步写入 log-service。
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

// WriteRequest 是操作日志写入请求。
type WriteRequest struct {
	Entry Entry `json:"entry"`
}

// WriteResponse 是操作日志写入响应。
type WriteResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// OperLogWriter 定义操作日志写入器，供中间件统一调用。
type OperLogWriter interface {
	Write(ctx context.Context, entry Entry) error
	Close() error
}

// Service 定义 log-service 需要实现的操作日志 gRPC 服务接口。
type Service interface {
	WriteOperLog(ctx context.Context, req *WriteRequest) (*WriteResponse, error)
}

// NopOperLogWriter 是操作日志空实现，适合在未配置日志中心时兜底。
type NopOperLogWriter struct{}

func (NopOperLogWriter) Write(context.Context, Entry) error { return nil }
func (NopOperLogWriter) Close() error                       { return nil }

// OperLogClient 是操作日志 gRPC 客户端。
type OperLogClient struct {
	conn *grpc.ClientConn
}

// NewOperLogClient 创建操作日志 gRPC 客户端。
func NewOperLogClient(target string, opts ...grpc.DialOption) (*OperLogClient, error) {
	conn, err := grpcx.NewClientConn(target, opts...)
	if err != nil {
		return nil, err
	}
	return &OperLogClient{conn: conn}, nil
}

// Write 调用远端 log-service 写入操作日志。
func (c *OperLogClient) Write(ctx context.Context, entry Entry) error {
	resp := &WriteResponse{}
	return c.conn.Invoke(ctx, WriteOperLogMethodFullName, &WriteRequest{Entry: entry}, resp)
}

// Close 关闭底层 gRPC 连接。
func (c *OperLogClient) Close() error {
	if c == nil {
		return nil
	}
	return grpcx.CloseClientConn(c.conn)
}

// RegisterService 注册操作日志 gRPC 服务。
func RegisterService(registrar grpc.ServiceRegistrar, srv Service) {
	registrar.RegisterService(&grpc.ServiceDesc{
		ServiceName: OperLogServiceName,
		HandlerType: (*Service)(nil),
		Methods: []grpc.MethodDesc{
			{
				MethodName: WriteOperLogMethodName,
				Handler:    writeOperLogHandler,
			},
		},
		Streams:  []grpc.StreamDesc{},
		Metadata: "oper_log",
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
		FullMethod: WriteOperLogMethodFullName,
	}
	handler := func(ctx context.Context, req any) (any, error) {
		return srv.(Service).WriteOperLog(ctx, req.(*WriteRequest))
	}
	return interceptor(ctx, in, info, handler)
}
