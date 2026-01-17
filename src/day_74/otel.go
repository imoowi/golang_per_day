package main

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

// initProvider 初始化 OpenTelemetry 的 Trace 和 Metric Provider
// 返回一个清理函数，用于在程序退出时关闭 provider
func initProvider() func(context.Context) {
	ctx := context.Background()
	// 这里的地址对应 K8s Service 名
	collectorAddr := "otel-collector.monitoring.svc.cluster.local:4317"

	// 创建资源，设置服务名称为 "gin-otel-demo"
	res, _ := resource.New(ctx, resource.WithAttributes(semconv.ServiceNameKey.String("gin-otel-demo")))

	// 1. 初始化 Trace Provider
	// 创建 OTLP gRPC 导出器，将 trace 数据发送到 OpenTelemetry Collector
	traceExporter, _ := otlptracegrpc.New(ctx,
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint(collectorAddr),
	)
	// 创建 TracerProvider，配置批处理导出器和资源
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(traceExporter),
		sdktrace.WithResource(res),
	)
	// 设置全局 TracerProvider
	otel.SetTracerProvider(tp)

	// 2. 初始化 Metric Provider
	// 创建 OTLP gRPC 指标导出器，将 metrics 数据发送到 OpenTelemetry Collector
	metricExporter, _ := otlpmetricgrpc.New(ctx,
		otlpmetricgrpc.WithInsecure(),
		otlpmetricgrpc.WithEndpoint(collectorAddr),
	)
	// 创建周期性读取器，定期从指标导出器读取数据
	reader := metric.NewPeriodicReader(metricExporter)
	// 创建 MeterProvider，配置读取器和资源
	mp := metric.NewMeterProvider(
		metric.WithReader(reader),
		metric.WithResource(res),
	)
	// 设置全局 MeterProvider
	otel.SetMeterProvider(mp)

	// 返回清理函数，用于优雅关闭
	return func(ctx context.Context) {
		_ = tp.Shutdown(ctx)
		_ = mp.Shutdown(ctx)
	}
}
