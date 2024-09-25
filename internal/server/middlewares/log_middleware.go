package middlewares

import (
	"github.com/labstack/echo/v4"
	"github.com/shestooy/go-musthave-metrics-tpl.git/internal/logger"
	"go.uber.org/zap"
	"time"
)

type logEntry struct {
	Time     string `json:"time"`
	RemoteIP string `json:"remote_ip"`
	Host     string `json:"host"`
	Method   string `json:"method"`
	URI      string `json:"uri"`
	Status   int    `json:"status"`
	Error    string `json:"error,omitempty"`
	Latency  string `json:"latency_human"`
	BytesIn  int64  `json:"bytes_in"`
	BytesOut int64  `json:"bytes_out"`
}

func Logging(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) (err error) {
		req := c.Request()
		res := c.Response()
		start := time.Now()

		if err = next(c); err != nil {
			c.Error(err)
		}

		end := time.Since(start)

		log := logEntry{
			Time:     start.Format(time.RFC3339Nano),
			RemoteIP: c.RealIP(),
			Host:     req.Host,
			Method:   req.Method,
			URI:      req.RequestURI,
			Status:   res.Status,
			Latency:  end.String(),
			BytesIn:  req.ContentLength,
			BytesOut: res.Size,
		}

		if err != nil {
			log.Error = err.Error()
		}

		logger.Log.Info("logging", zap.Any("fields", log))

		return err
	}
}
