package log_grpc

import (
	"context"

	"github.com/calmlax/aevons-framework/config"
	"github.com/calmlax/aevons-framework/core/grpcx"
	"google.golang.org/grpc"
)

const (
	// RegistryServiceName 是 log-service 在 Consul 中的服务注册名。
	RegistryServiceName = "log-service"
)

// newServiceConn 使用默认日志服务注册名，通过 Consul 发现并建立 gRPC 连接。
func newServiceConn(consulCfg config.ConsulConfig, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	return grpcx.NewClientConnFromConsul(consulCfg, RegistryServiceName, opts...)
}

// invokeUnary 封装一元 gRPC 调用，减少各契约文件里的重复 Invoke 样板。
func invokeUnary(ctx context.Context, conn *grpc.ClientConn, fullMethod string, req, resp any) error {
	return conn.Invoke(ctx, fullMethod, req, resp)
}

// closeConn 安全关闭底层 gRPC 连接。
func closeConn(conn *grpc.ClientConn) error {
	return grpcx.CloseClientConn(conn)
}

// registerUnaryService 封装一元服务注册，减少各契约文件里的 ServiceDesc 样板。
func registerUnaryService(
	registrar grpc.ServiceRegistrar,
	srv any,
	serviceName string,
	handlerType any,
	methodName string,
	handler grpc.MethodHandler,
	metadata string,
) {
	registrar.RegisterService(&grpc.ServiceDesc{
		ServiceName: serviceName,
		HandlerType: handlerType,
		Methods: []grpc.MethodDesc{
			{
				MethodName: methodName,
				Handler:    handler,
			},
		},
		Streams:  []grpc.StreamDesc{},
		Metadata: metadata,
	}, srv)
}
