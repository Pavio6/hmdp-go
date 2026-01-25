package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const requestIDKey = "request_id"

// RequestIDMiddleware 请求ID中间件：读取或生成 request_id 并写入响应头
func RequestIDMiddleware(header string) gin.HandlerFunc {
	if header == "" {
		header = "X-Request-ID"
	}
	return func(c *gin.Context) {
		rid := c.GetHeader(header)
		if rid == "" {
			rid = uuid.NewString()
		}
		// 将 request_id 写入上下文
		c.Set(requestIDKey, rid)
		// 将 request_id 写入响应头
		c.Writer.Header().Set(header, rid)
		c.Next()
	}
}

// RequestIDFromContext 从上下文中读取 request_id
func RequestIDFromContext(c *gin.Context) string {
	if value, ok := c.Get(requestIDKey); ok {
		if id, ok := value.(string); ok {
			return id
		}
	}
	return ""
}
