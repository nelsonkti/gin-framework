package tracer

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.15.0"
	"log"
	"os"
	"time"
)

type Tracer struct {
	closerFunc func()
}

const (
	SERVICE_VERSION = "1.0"
)

func NewOpentelemetry(serviceName, env, endpoint, urlPath string) *Tracer {
	ctx := context.Background()

	var traceExporter *otlptrace.Exporter
	var batchSpanProcessor sdktrace.SpanProcessor

	traceExporter, batchSpanProcessor = newHTTPExporterAndSpanProcessor(ctx, endpoint, urlPath)

	traceProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(newResource(ctx, serviceName, env)),
		sdktrace.WithSpanProcessor(batchSpanProcessor),
	)

	otel.SetTracerProvider(traceProvider)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	return &Tracer{
		closerFunc: func() {
			cxt, cancel := context.WithTimeout(ctx, time.Second)
			defer cancel()
			if err := traceExporter.Shutdown(cxt); err != nil {
				otel.Handle(err)
			}
		},
	}
}

// 设置应用资源
func newResource(ctx context.Context, serviceName, env string) *resource.Resource {
	hostName, _ := os.Hostname()

	r, err := resource.New(
		ctx,
		resource.WithFromEnv(),
		resource.WithProcess(),
		resource.WithTelemetrySDK(),
		resource.WithHost(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String(serviceName),        // 应用名
			semconv.ServiceVersionKey.String(SERVICE_VERSION), // 应用版本
			semconv.DeploymentEnvironmentKey.String(env),      // 部署环境
			semconv.HostNameKey.String(hostName),              // 主机名
		),
	)

	if err != nil {
		log.Fatalf("%s: %v", "Failed to create OpenTelemetry resource", err)
	}
	return r
}

func newHTTPExporterAndSpanProcessor(ctx context.Context, endpoint, urlPath string) (*otlptrace.Exporter, sdktrace.SpanProcessor) {
	traceExporter, err := otlptrace.New(ctx, otlptracehttp.NewClient(
		otlptracehttp.WithEndpoint(endpoint),
		otlptracehttp.WithURLPath(urlPath),
		otlptracehttp.WithInsecure(),
		otlptracehttp.WithCompression(1)))

	if err != nil {
		log.Fatalf("%s: %v", "Failed to create the OpenTelemetry trace exporter", err)
	}

	batchSpanProcessor := sdktrace.NewBatchSpanProcessor(traceExporter)

	return traceExporter, batchSpanProcessor
}

func (t *Tracer) Close() {
	t.closerFunc()
}

// InitOpenTelemetry OpenTelemetry 初始化方法
//func InitOpenTelemetry() func() {
//	ctx := context.Background()
//
//	var traceExporter *otlptrace.Exporter
//	var batchSpanProcessor sdktrace.SpanProcessor
//
//	traceExporter, batchSpanProcessor = newHTTPExporterAndSpanProcessor(ctx)
//
//	traceProvider := sdktrace.NewTracerProvider(
//		sdktrace.WithSampler(sdktrace.AlwaysSample()),
//		sdktrace.WithResource(newResource(ctx)),
//		sdktrace.WithSpanProcessor(batchSpanProcessor))
//
//	otel.SetTracerProvider(traceProvider)
//	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
//
//	return func() {
//		cxt, cancel := context.WithTimeout(ctx, time.Second)
//		defer cancel()
//		if err := traceExporter.Shutdown(cxt); err != nil {
//			otel.Handle(err)
//		}
//	}
//}
