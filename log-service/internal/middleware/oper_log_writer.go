package middleware

import (
	"context"
	"log-service/internal/model"
	"log-service/internal/repository"
	"log-service/internal/service"

	"gorm.io/gorm"
)

// DBOperLogWriter 供 log-service 自身使用，直接落库写入操作日志。
type DBOperLogWriter struct {
	svc service.OperLogService
}

func NewDBOperLogWriter(db *gorm.DB) OperLogWriter {
	repo := repository.NewOperLogRepository(db)
	return &DBOperLogWriter{
		svc: service.NewOperLogService(repo),
	}
}

func (w *DBOperLogWriter) Write(_ context.Context, entry OperLogEntry) error {
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
	return w.svc.Create(&record)
}

func (w *DBOperLogWriter) Close() error {
	return nil
}
