package grpcs

import (
	"context"
	"log-service/model"
	"log-service/repository"
	"log-service/service"

	"aevons-grpc/log_grpc"

	"gorm.io/gorm"
)

// OperLogServiceServer 实现 log-service 的操作日志 gRPC 服务。
type OperLogServiceServer struct {
	svc service.OperLogService
}

func NewOperLogServiceServer(db *gorm.DB) *OperLogServiceServer {
	repo := repository.NewOperLogRepository(db)
	svc := service.NewOperLogService(repo)
	return &OperLogServiceServer{svc: svc}
}

func (s *OperLogServiceServer) WriteOperLog(ctx context.Context, req *log_grpc.WriteRequest) (*log_grpc.WriteResponse, error) {
	entry := req.Entry
	record := model.OperLog{
		Module:      entry.Module,
		Type:        entry.Type,
		Description: entry.Description,
		Method:      entry.Method,
		Url:         entry.URL,
		Ip:          entry.IP,
		Location:    entry.Location,
		Payload:     entry.Payload,
		Result:      entry.Result,
		Device:      entry.Device,
		Os:          entry.OS,
		Browser:     entry.Browser,
		Status:      entry.Status,
		Error:       entry.Error,
		Time:        entry.TimeMS,
		UserId:      entry.UserID,
		Username:    entry.Username,
		OperAt:      entry.OperAt,
	}

	if err := s.svc.Create(&record); err != nil {
		return nil, err
	}

	return &log_grpc.WriteResponse{
		Success: true,
		Message: "ok",
	}, nil
}
