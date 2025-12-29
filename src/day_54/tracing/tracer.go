package tracing

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// Init 初始化 OpenTelemetry 跟踪系统
// serviceName: 服务名称，用于在分布式追踪系统中标识当前服务
func Init(serviceName string) {
	// 创建带超时的上下文，防止初始化过程阻塞过久
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	// 创建 OTLP HTTP 导出器，将跟踪数据发送到 Jaeger 服务器
	exp, err := otlptracehttp.New(
		ctx,
		// 设置 Jaeger 服务器端点
		otlptracehttp.WithEndpoint("jaeger:4318"),
		// 允许使用非安全连接（生产环境建议使用安全连接）
		otlptracehttp.WithInsecure(),
	)
	if err != nil {
		// 如果导出器创建失败，直接panic
		panic(err)
	}

	// 创建跟踪提供器（TracerProvider）
	tp := sdktrace.NewTracerProvider(
		// 使用批处理器发送跟踪数据，提高性能
		sdktrace.WithBatcher(exp),
		// 设置资源属性，包括服务名称
		sdktrace.WithResource(
			resource.NewWithAttributes(
				"",
				attribute.String("service.name", serviceName),
			),
		),
	)

	// 设置全局跟踪提供器
	otel.SetTracerProvider(tp)
}
