package sys_grpc

import (
	"context"

	"github.com/calmlax/aevons-framework/config"
	"github.com/calmlax/aevons-framework/grpcx"
	"google.golang.org/grpc"
)

const (
	// RegistryServiceName 是 sys-service 在 Consul 中的服务注册名。
	RegistryServiceName = "sys-service"
)

// newServiceConn 使用默认系统服务注册名，通过 Consul 发现并建立 gRPC 连接。
func newServiceConn(consulCfg config.ConsulConfig, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	return grpcx.NewClientConnFromConsul(consulCfg, RegistryServiceName, opts...)
}

// invokeUnary 封装一元 gRPC 调用，减少契约文件里的重复 Invoke 样板。
func invokeUnary(ctx context.Context, conn *grpc.ClientConn, fullMethod string, req, resp any) error {
	return conn.Invoke(ctx, fullMethod, req, resp)
}

// closeConn 安全关闭底层 gRPC 连接。
func closeConn(conn *grpc.ClientConn) error {
	return grpcx.CloseClientConn(conn)
}
