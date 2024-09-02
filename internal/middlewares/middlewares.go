package middlewares

import (
	"compress/gzip"
	l "github.com/shestooy/go-musthave-metrics-tpl.git/internal/logger"
	"go.uber.org/zap"
	"io"
	"net/http"
	"strings"
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

type compressWriter struct {
	w  http.ResponseWriter
	zw *gzip.Writer
}

func newCompressWriter(w http.ResponseWriter) *compressWriter {
	return &compressWriter{
		w:  w,
		zw: gzip.NewWriter(w),
	}
}

func (c *compressWriter) Header() http.Header {
	return c.w.Header()
}

func (c *compressWriter) Write(p []byte) (int, error) {
	return c.zw.Write(p)
}

func (c *compressWriter) WriteHeader(statusCode int) {
	if statusCode < 300 {
		c.w.Header().Set("Content-Encoding", "gzip")
	}
	c.w.WriteHeader(statusCode)
}

func (c *compressWriter) Close() error {
	return c.zw.Close()
}

type compressReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

func newCompressReader(r io.ReadCloser) (*compressReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &compressReader{
		r:  r,
		zr: zr,
	}, nil
}

func (c compressReader) Read(p []byte) (n int, err error) {
	return c.zr.Read(p)
}

func (c *compressReader) Close() error {
	if err := c.r.Close(); err != nil {
		return err
	}
	return c.zr.Close()
}

func GzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		w := res

		acceptEnc := req.Header.Get("Accept-Encoding")
		if ok := strings.Contains(acceptEnc, "gzip"); ok {
			cw := newCompressWriter(res)
			w = cw
			w.Header().Set("Content-Encoding", "gzip")
			defer cw.Close()
		}
		if ok := strings.Contains(req.Header.Get("Content-Encoding"), "gzip"); ok {
			cr, err := newCompressReader(req.Body)
			if err != nil {
				res.WriteHeader(http.StatusInternalServerError)
				return
			}
			req.Body = cr
			defer cr.Close()
		}

		next.ServeHTTP(w, req)
	})
}
