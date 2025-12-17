package middleware

import (
	"hmdp-backend/internal/dto/result"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ErrorHandler 通过将 panic 转换为 JSON 响应来模仿 WebExceptionAdvice。
func ErrorHandler(log *zap.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if rec := recover(); rec != nil {
				log.Error("panic recovered", zap.Any("error", rec))
				ctx.JSON(http.StatusInternalServerError, result.Fail("服务器异常"))
				ctx.Abort()
			}
		}()
		ctx.Next()
	}
}
