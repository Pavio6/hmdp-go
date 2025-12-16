package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"hmdp-backend/internal/dto"
)

// ErrorHandler mimics WebExceptionAdvice by converting panics to JSON responses.
func ErrorHandler(log *zap.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if rec := recover(); rec != nil {
				log.Error("panic recovered", zap.Any("error", rec))
				ctx.JSON(http.StatusInternalServerError, dto.Fail("服务器异常"))
				ctx.Abort()
			}
		}()
		ctx.Next()
	}
}
