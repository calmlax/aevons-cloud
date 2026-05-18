package grpcs

import (
	"context"
	"log-service/internal/model"
	"log-service/internal/repository"
	"log-service/internal/service"

	"internal-grpc/log_grpc"

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

func (s *LoginLogServiceServer) GetLatestLoginLog(ctx context.Context, req *log_grpc.GetLatestLoginLogRequest) (*log_grpc.GetLatestLoginLogResponse, error) {
	records, err := s.svc.GetLatestLoginLog(req.Username, req.Limit)
	if err != nil {
		return nil, err
	}

	entries := make([]log_grpc.LoginEntry, len(records))
	for i, record := range records {
		entries[i] = log_grpc.LoginEntry{
			ID:        record.Id,
			Username:  record.Username,
			ClientID:  record.ClientId,
			GrantType: record.GrantType,
			OS:        record.Os,
			Browser:   record.Browser,
			IP:        record.Ip,
			Location:  record.Location,
			Status:    record.Status,
			Msg:       record.Msg,
			LoginAt:   record.LoginAt,
		}
	}

	return &log_grpc.GetLatestLoginLogResponse{
		Entries: entries,
	}, nil
}
