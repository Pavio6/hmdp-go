package observability

import (
	"context"
	"errors"
	"math"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

type TracingConfig struct {
	Enabled          bool
	OTLPGrpcEndpoint string
	Insecure         bool
	SampleRate       float64
}

type ResourceConfig struct {
	ServiceName string
	Environment string
}
// SetupTracing 初始化 OpenTelemetry 追踪系统
func SetupTracing(ctx context.Context, tracing TracingConfig, resourceCfg ResourceConfig) (func(context.Context) error, error) {
	if !tracing.Enabled {
		return func(context.Context) error { return nil }, nil
	}
	if tracing.OTLPGrpcEndpoint == "" {
		return nil, errors.New("otlp grpc endpoint is required when tracing is enabled")
	}
	// 创建 OTLP gRPC 导出器
	opts := []otlptracegrpc.Option{otlptracegrpc.WithEndpoint(tracing.OTLPGrpcEndpoint)}
	if tracing.Insecure {
		opts = append(opts, otlptracegrpc.WithInsecure())
	}
	// 创建 OTLP 追踪导出器
	exporter, err := otlptracegrpc.New(ctx, opts...)
	if err != nil {
		return nil, err
	}

	sampleRate := tracing.SampleRate
	if math.IsNaN(sampleRate) || sampleRate <= 0 || sampleRate > 1 {
		sampleRate = 1
	}
	// 创建资源信息
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(resourceCfg.ServiceName),
			attribute.String("deployment.environment", resourceCfg.Environment),
		),
	)
	if err != nil {
		return nil, err
	}
	// 创建 TracerProvider
	provider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithSampler(sdktrace.ParentBased(sdktrace.TraceIDRatioBased(sampleRate))),
		sdktrace.WithResource(res),
	)
	// 设置全局 TracerProvider 和 传播器
	otel.SetTracerProvider(provider)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	return provider.Shutdown, nil
}
