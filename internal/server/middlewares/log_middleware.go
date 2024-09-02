package middlewares

import (
	l "github.com/shestooy/go-musthave-metrics-tpl.git/internal/logger"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
	size       int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	size, err := rw.ResponseWriter.Write(b)
	rw.size += size
	return size, err
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		start := time.Now()

		l.Log.Info("Request received",
			zap.String("URI", req.RequestURI),
			zap.String("Method", req.Method))

		rW := &responseWriter{
			ResponseWriter: res,
			statusCode:     http.StatusOK}
		next.ServeHTTP(rW, req)
		d := time.Since(start)

		l.Log.Info("Response sent",
			zap.Int("Status", rW.statusCode),
			zap.Int("Content Length", rW.size),
			zap.Duration("Duration", d))
	})
}
