package observability

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type HTTPMetrics struct {
	registry    *prometheus.Registry
	inFlight    prometheus.Gauge
	reqTotal    *prometheus.CounterVec
	reqDuration *prometheus.HistogramVec
}

// NewHTTPMetrics 创建 HTTP 指标收集器，并注册到给定的 Registry
func NewHTTPMetrics(registry *prometheus.Registry, serviceName string) *HTTPMetrics {
	if registry == nil {
		// 未传入 registry 时，创建默认的并挂载 Go/runtime 与进程指标。
		registry = prometheus.NewRegistry()
		registry.MustRegister(collectors.NewGoCollector())
		registry.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
	}

	constLabels := prometheus.Labels{}
	if serviceName != "" {
		constLabels["service"] = serviceName
	}
	// 当前正在处理的请求数
	inFlight := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace:   "http",
		Subsystem:   "server",
		Name:        "in_flight_requests",
		Help:        "Number of in-flight HTTP requests.",
		ConstLabels: constLabels,
	})
	// HTTP请求总数
	reqTotal := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace:   "http",
		Subsystem:   "server",
		Name:        "requests_total",
		Help:        "Total number of HTTP requests.",
		ConstLabels: constLabels,
	}, []string{"method", "path", "status"})
	// 请求耗时分布
	reqDuration := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace:   "http",
		Subsystem:   "server",
		Name:        "request_duration_seconds",
		Help:        "HTTP request duration in seconds.",
		Buckets:     prometheus.DefBuckets,
		ConstLabels: constLabels,
	}, []string{"method", "path", "status"})

	registry.MustRegister(inFlight, reqTotal, reqDuration)

	return &HTTPMetrics{
		registry:    registry,
		inFlight:    inFlight,
		reqTotal:    reqTotal,
		reqDuration: reqDuration,
	}
}

// Handler 返回 Prometheus 指标导出 handler
func (m *HTTPMetrics) Handler() http.Handler {
	return promhttp.HandlerFor(m.registry, promhttp.HandlerOpts{})
}

// Middleware 返回 Gin 中间件，用于采集请求统计与耗时
func (m *HTTPMetrics) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		m.inFlight.Inc()
		start := time.Now()
		c.Next()
		m.inFlight.Dec()

		status := c.Writer.Status()
		path := c.FullPath()
		if path == "" {
			// 未命中路由时回退到原始路径，避免丢失路径标签。
			path = c.Request.URL.Path
		}
		method := c.Request.Method
		statusLabel := strconv.Itoa(status)
		m.reqTotal.WithLabelValues(method, path, statusLabel).Inc()
		m.reqDuration.WithLabelValues(method, path, statusLabel).Observe(time.Since(start).Seconds())
	}
}
