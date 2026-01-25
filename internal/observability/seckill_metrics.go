package observability

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// SeckillMetrics 定义秒杀相关的指标
type SeckillMetrics struct {
	seckillTotal        *prometheus.CounterVec
	seckillLatency      *prometheus.HistogramVec // 秒杀请求耗时分布
	kafkaPublishTotal   *prometheus.CounterVec 
	kafkaConsumeTotal   *prometheus.CounterVec
	kafkaConsumeLatency *prometheus.HistogramVec // Kafka消费处理耗时分布
	retryTotal          *prometheus.CounterVec
}

func NewSeckillMetrics(registry *prometheus.Registry, serviceName string) *SeckillMetrics {
	if registry == nil {
		registry = NewMetricsRegistry()
	}

	constLabels := prometheus.Labels{}
	if serviceName != "" {
		constLabels["service"] = serviceName
	}

	seckillTotal := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace:   "seckill",
		Subsystem:   "order",
		Name:        "requests_total",
		Help:        "Total seckill requests.",
		ConstLabels: constLabels,
	}, []string{"result", "reason"})

	seckillLatency := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace:   "seckill",
		Subsystem:   "order",
		Name:        "request_duration_seconds",
		Help:        "Seckill request duration in seconds.",
		Buckets:     prometheus.DefBuckets,
		ConstLabels: constLabels,
	}, []string{"result"})

	kafkaPublishTotal := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace:   "seckill",
		Subsystem:   "kafka",
		Name:        "publish_total",
		Help:        "Total kafka publish attempts.",
		ConstLabels: constLabels,
	}, []string{"topic", "result"})

	kafkaConsumeTotal := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace:   "seckill",
		Subsystem:   "kafka",
		Name:        "consume_total",
		Help:        "Total kafka consume results.",
		ConstLabels: constLabels,
	}, []string{"topic", "result"})

	kafkaConsumeLatency := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace:   "seckill",
		Subsystem:   "kafka",
		Name:        "consume_duration_seconds",
		Help:        "Kafka consume handling duration in seconds.",
		Buckets:     prometheus.DefBuckets,
		ConstLabels: constLabels,
	}, []string{"topic", "result"})

	retryTotal := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace:   "seckill",
		Subsystem:   "kafka",
		Name:        "retry_total",
		Help:        "Total retry or DLQ events.",
		ConstLabels: constLabels,
	}, []string{"phase"})

	registry.MustRegister(seckillTotal, seckillLatency, kafkaPublishTotal, kafkaConsumeTotal, kafkaConsumeLatency, retryTotal)

	return &SeckillMetrics{
		seckillTotal:        seckillTotal,
		seckillLatency:      seckillLatency,
		kafkaPublishTotal:   kafkaPublishTotal,
		kafkaConsumeTotal:   kafkaConsumeTotal,
		kafkaConsumeLatency: kafkaConsumeLatency,
		retryTotal:          retryTotal,
	}
}
// ObserveSeckill 记录一次秒杀请求的结果与耗时
func (m *SeckillMetrics) ObserveSeckill(result, reason string, duration time.Duration) {
	if m == nil {
		return
	}
	if reason == "" {
		reason = "unknown"
	}
	m.seckillTotal.WithLabelValues(result, reason).Inc()
	m.seckillLatency.WithLabelValues(result).Observe(duration.Seconds())
}
// ObserveKafkaPublish 记录一次 Kafka 消息发布的结果
func (m *SeckillMetrics) ObserveKafkaPublish(topic, result string) {
	if m == nil {
		return
	}
	m.kafkaPublishTotal.WithLabelValues(topic, result).Inc()
}
// ObserveKafkaConsume 记录一次 Kafka 消息消费的结果与耗时
func (m *SeckillMetrics) ObserveKafkaConsume(topic, result string, duration time.Duration) {
	if m == nil {
		return
	}
	m.kafkaConsumeTotal.WithLabelValues(topic, result).Inc()
	m.kafkaConsumeLatency.WithLabelValues(topic, result).Observe(duration.Seconds())
}
// ObserveRetry 记录一次重试或死信处理事件
func (m *SeckillMetrics) ObserveRetry(phase string) {
	if m == nil {
		return
	}
	m.retryTotal.WithLabelValues(phase).Inc()
}
