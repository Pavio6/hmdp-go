package handler

import (
	"context"
	"errors"
	"net"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	"hmdp-backend/internal/data"
)

// HealthHandler 
type HealthHandler struct {
	db           sqlDB
	redis        *redis.Client
	kafkaBrokers []string
	log          *zap.Logger
	checkTimeout time.Duration
}
// sqlDB 定义了数据库连接需要实现的接口
type sqlDB interface {
	PingContext(ctx context.Context) error
}

// NewHealthHandler 创建一个新的 HealthHandler 实例
func NewHealthHandler(db sqlDB, redisClient *redis.Client, kafkaBrokers []string, log *zap.Logger) *HealthHandler {
	return &HealthHandler{
		db:           db,
		redis:        redisClient,
		kafkaBrokers: kafkaBrokers,
		log:          log,
		checkTimeout: 2 * time.Second,
	}
}

// Healthz 返回服务健康状态
func (h *HealthHandler) Healthz(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// Readyz 返回服务就绪状态（服务是否可以对外接收流量）
func (h *HealthHandler) Readyz(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), h.checkTimeout)
	defer cancel()

	checks := map[string]string{}
	if err := h.db.PingContext(ctx); err != nil {
		checks["mysql"] = err.Error()
	}
	if err := data.Ping(ctx, h.redis); err != nil {
		checks["redis"] = err.Error()
	}
	if err := checkKafka(ctx, h.kafkaBrokers); err != nil {
		checks["kafka"] = err.Error()
	}

	if len(checks) > 0 {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "fail",
			"checks": checks,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
// checkKafka 检查与 Kafka 的连接
func checkKafka(ctx context.Context, brokers []string) error {
	if len(brokers) == 0 {
		return errors.New("no kafka brokers configured")
	}
	// 创建 建立网络连接对象
	dialer := net.Dialer{Timeout: time.Second}
	var lastErr error
	for _, broker := range brokers {
		// 尝试连接每个 broker
		conn, err := dialer.DialContext(ctx, "tcp", broker)
		if err == nil {
			_ = conn.Close()
			return nil
		}
		lastErr = err
	}
	return lastErr
}
