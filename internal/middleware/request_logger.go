package middleware

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// RequestLogger 是一个 Gin 中间件，用于记录 HTTP 请求日志
func RequestLogger(log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		status := c.Writer.Status()
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}

		fields := []zap.Field{
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.Int("status", status),
			zap.Duration("latency", time.Since(start)),
			zap.String("client_ip", c.ClientIP()),
			zap.Int("bytes_in", int(c.Request.ContentLength)),
			zap.Int("bytes_out", c.Writer.Size()),
		}

		if rid := RequestIDFromContext(c); rid != "" {
			fields = append(fields, zap.String("request_id", rid))
		}

		span := trace.SpanFromContext(c.Request.Context())
		if span != nil {
			spanCtx := span.SpanContext()
			if spanCtx.IsValid() {
				fields = append(fields,
					zap.String("trace_id", spanCtx.TraceID().String()),
					zap.String("span_id", spanCtx.SpanID().String()),
				)
			}
		}

		if len(c.Errors) > 0 {
			fields = append(fields, zap.String("errors", strings.TrimSpace(c.Errors.String())))
		}

		switch {
		case status >= 500:
			log.Error("http request", fields...)
		case status >= 400:
			log.Warn("http request", fields...)
		default:
			log.Info("http request", fields...)
		}
	}
}
