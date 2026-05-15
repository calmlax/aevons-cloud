package grpcs

import (
	"context"
	"log-service/model"
	"log-service/repository"
	"log-service/service"

	"aevons-grpc/log_grpc"

	"gorm.io/gorm"
)

// LoginLogServiceServer 实现 log-service 的登录日志 gRPC 服务。
type LoginLogServiceServer struct {
	svc service.LoginLogService
}

func NewLoginLogServiceServer(db *gorm.DB) *LoginLogServiceServer {
	repo := repository.NewLoginLogRepository(db)
	svc := service.NewLoginLogService(repo)
	return &LoginLogServiceServer{svc: svc}
}

func (s *LoginLogServiceServer) WriteLoginLog(ctx context.Context, req *log_grpc.WriteLoginRequest) (*log_grpc.WriteLoginResponse, error) {
	entry := req.Entry
	record := model.LoginLog{
		Username:  entry.Username,
		ClientId:  entry.ClientID,
		GrantType: entry.GrantType,
		Os:        entry.OS,
		Browser:   entry.Browser,
		Ip:        entry.IP,
		Location:  entry.Location,
		Status:    entry.Status,
		Msg:       entry.Msg,
		LoginAt:   entry.LoginAt,
	}

	if err := s.svc.Create(&record); err != nil {
		return nil, err
	}

	return &log_grpc.WriteLoginResponse{
		Success: true,
		Message: "ok",
	}, nil
}
