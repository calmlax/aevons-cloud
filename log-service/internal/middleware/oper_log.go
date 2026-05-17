package middleware

import (
	"bytes"
	"context"
	"net/http"
	"time"

	frameworkauth "github.com/calmlax/aevons-framework/auth"
	"github.com/calmlax/aevons-framework/consts"
	"github.com/calmlax/aevons-framework/utils"
	"github.com/calmlax/aevons-framework/xlog"
	"github.com/gin-gonic/gin"
)

type OperLogEntry struct {
	Module      string
	Type        string
	Description string
	Method      string
	URL         string
	IP          string
	Location    string
	Payload     string
	Result      string
	Device      string
	OS          string
	Browser     string
	Status      int16
	Error       string
	TimeMS      int64
	UserID      int64
	Username    string
	OperAt      time.Time
}

type OperLogWriter interface {
	Write(ctx context.Context, entry OperLogEntry) error
}

// OperLog 采集请求上下文，并通过统一 writer 写入操作日志。
func OperLog(writer OperLogWriter, module string, bizType consts.BizType) gin.HandlerFunc {
	if writer == nil {
		return func(c *gin.Context) {
			c.Next()
		}
	}

	return func(c *gin.Context) {
		start := time.Now()
		req := c.Request
		payload := utils.CaptureRequestPayload(req, 2000)

		w := &utils.ResponseBodyWriter{
			ResponseWriter: c.Writer,
			Body:           &bytes.Buffer{},
		}
		c.Writer = w

		c.Next()

		statusCode := c.Writer.Status()
		status := int16(1)
		if statusCode >= http.StatusBadRequest {
			status = 0
		}

		agent := utils.ParseClientAgent(req.UserAgent())
		userID, username, _ := frameworkauth.GetCurrentIdentity(c.Request.Context())

		ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
		defer cancel()

		if err := writer.Write(ctx, OperLogEntry{
			Module:      module,
			Type:        string(bizType),
			Description: module,
			Method:      req.Method,
			URL:         utils.BuildRequestURL(req),
			IP:          utils.ClientIPFromRequest(req),
			Location:    "",
			Payload:     payload,
			Result:      utils.TruncateText(w.Body.String(), 2000),
			Device:      agent.Device,
			OS:          agent.OS,
			Browser:     agent.Browser,
			Status:      status,
			Error:       utils.ExtractResponseError(statusCode, c, w.Body.Bytes(), 2000),
			TimeMS:      time.Since(start).Milliseconds(),
			UserID:      userID,
			Username:    username,
			OperAt:      start,
		}); err != nil {
			xlog.Warn("write oper log failed: %v", err)
		}
	}
}
