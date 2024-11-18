package tracer

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"os"
)

const TracerKey = "otel-go-contrib-tracer"

// Span 创建 span
func Span(ctx context.Context, spanName string, opts ...trace.SpanStartOption) (spanctx context.Context, span trace.Span) {
	value := ctx.Value(TracerKey)
	tracer, ok := value.(trace.Tracer)
	if !ok {
		return ctx, nil
	}

	// gin 特殊
	if c, ok := ctx.(*gin.Context); ok {
		spanctx, span = tracer.Start(c.Request.Context(), spanName, opts...)

		spanctx = context.WithValue(spanctx, TracerKey, tracer)
	} else {
		spanctx, span = tracer.Start(ctx, spanName, opts...)
	}

	// 设置 Attr
	attrkv, ok := ctx.Value("attrkv").(map[string]string)
	if ok {
		SpanSetStringAttr(span, attrkv)
	}

	SpanSetStringAttr(span, map[string]string{
		"server.host": os.Getenv("HOSTNAME"),
	})

	return spanctx, span
}

// SpanSetStringAttr 设置 span 属性
func SpanSetStringAttr(span trace.Span, kvs map[string]string) {
	attrkv := []attribute.KeyValue{}

	for k, v := range kvs {
		attrkv = append(attrkv, attribute.KeyValue{
			Key:   attribute.Key(k),
			Value: attribute.StringValue(v),
		})
	}

	span.SetAttributes(attrkv...)
}

// SpanContextWithAttr 设置 span 属性
func SpanContextWithAttr(ctx context.Context, kv map[string]string) context.Context {

	value := ctx.Value("attrkv")
	attrkv, ok := value.(map[string]string)
	if !ok {
		attrkv = make(map[string]string, 0)
	}

	for k, v := range kv {
		attrkv[k] = v
	}

	return context.WithValue(ctx, "attrkv", attrkv)
}
