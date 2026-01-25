package observability

import (
	"context"
	"strings"

	"github.com/segmentio/kafka-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)
// kafkaHeaderCarrier 实现了 TextMapCarrier 接口，用于在 Kafka 消息头中注入和提取追踪信息
type kafkaHeaderCarrier struct {
	headers *[]kafka.Header
}
// Get 返回指定 key 的值
func (c kafkaHeaderCarrier) Get(key string) string {
	if c.headers == nil {
		return ""
	}
	// 遍历 headers 查找对应的 key
	for _, h := range *c.headers {
		// 忽略大小写比较 key
		if strings.EqualFold(h.Key, key) {
			return string(h.Value)
		}
	}
	return ""
}
// Set 设置指定 key 的值
func (c kafkaHeaderCarrier) Set(key, value string) {
	if c.headers == nil {
		return
	}
	headers := *c.headers
	// 遍历 headers 查找对应的 key 并更新值
	for i, h := range headers {
		if strings.EqualFold(h.Key, key) {
			headers[i].Value = []byte(value)
			*c.headers = headers
			return
		}
	}
	// 如果没有找到对应的 key，则添加新的 header
	headers = append(headers, kafka.Header{Key: key, Value: []byte(value)})
	*c.headers = headers
}
// Keys 返回所有的 key 列表
func (c kafkaHeaderCarrier) Keys() []string {
	if c.headers == nil {
		return nil
	}
	// 初始化字符串切片列表
	keys := make([]string, 0, len(*c.headers))
	for _, h := range *c.headers {
		keys = append(keys, h.Key)
	}
	return keys
}

// InjectKafkaHeaders 将当前跟踪上下文写入 kafka headers 中
func InjectKafkaHeaders(ctx context.Context, headers *[]kafka.Header) {
	if headers == nil {
		return
	}
	otel.GetTextMapPropagator().Inject(ctx, kafkaHeaderCarrier{headers: headers})
}

// ExtractKafkaContext 从 kafka headers 中读取 trace 上下文
func ExtractKafkaContext(ctx context.Context, headers []kafka.Header) context.Context {
	carrier := kafkaHeaderCarrier{headers: &headers}
	return otel.GetTextMapPropagator().Extract(ctx, carrier)
}

// InjectKafkaBaggage 将当前 baggage 写入 kafka headers 中
func InjectKafkaBaggage(ctx context.Context, headers *[]kafka.Header, propagator propagation.TextMapPropagator) {
	if headers == nil || propagator == nil {
		return
	}
	propagator.Inject(ctx, kafkaHeaderCarrier{headers: headers})
}
