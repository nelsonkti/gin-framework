package middleware

import (
	"bytes"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go-framework/internal/server"
	"go-framework/util/helper"
	"go-framework/util/xhttp"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"io"
	"strings"
)

// OTELMiddleware otel中间件
func OTELMiddleware(svc *server.SvcContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		span := trace.SpanFromContext(c.Request.Context())
		spanctx := c.Request.Context()
		if span == nil {
			c.Next()
			return
		}

		isRootSpan := !span.SpanContext().IsValid() || !span.SpanContext().IsRemote()

		if isRootSpan {
			span.SetAttributes(attribute.String("span.origin", "internal"))
		} else {
			span.SetAttributes(attribute.String("span.origin", "external"))
		}

		defer span.End()

		// 记录请求体
		if c.Request.Body != nil {
			reqData := getRequestData(c, svc)
			marshal, _ := helper.Marshal(reqData)
			span.SetAttributes(attribute.String("http.request_body", string(marshal)))
		}
		wireHeader(c, spanctx, span)
		rm := xhttp.NewRespJsonModifier(c)
		c.Writer = rm

		c.Next() // 执行下一个中间件和请求处理

		otelReport(c, span, rm)
	}
}

func wireHeader(c *gin.Context, spanctx context.Context, span trace.Span) {
	// 4. 应答客户端时， 在 Header 中默认添加 TraceID
	traceid := span.SpanContext().TraceID().String()
	c.Header("TraceID", traceid)

	// 6. 向后传递 Header: traceparent
	pp := propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
	)

	carrier := propagation.MapCarrier{}
	pp.Inject(spanctx, carrier)

	for k, v := range carrier {
		c.Header(k, v)
	}
}

// otelReport otel上报
func otelReport(c *gin.Context, span trace.Span, rm *xhttp.RespJsonModifier) {
	// 记录响应体
	if checkJson(rm) {
		data := make(map[string]interface{})
		rm.WriteMeta(c, data)

		body := rm.Body()
		_ = helper.UmMarshal(body, &data)
		status := "200"
		if data["code"] != 0 {
			status = fmt.Sprintf("%d", helper.ConvertToInt(data["code"]))
		}

		span.SetAttributes(attribute.String("http.status_code", status))
		span.SetAttributes(attribute.String("http.response_body", string(body)))
	}

	if cb, ok := c.Get(gin.BodyBytesKey); ok {
		if cbb, ok := cb.([]byte); ok {
			c.Request.Body = io.NopCloser(bytes.NewBuffer(cbb))
		}
	}
}

// checkJson 函数用于检查响应是否为 JSON 格式
func checkJson(rm *xhttp.RespJsonModifier) bool {
	contentType := rm.ResponseWriter.Header().Get("Content-Type")
	return strings.Contains(contentType, "application/json")
}

// getRequestData 用于获取请求数据。
func getRequestData(c *gin.Context, svc *server.SvcContext) map[string]interface{} {
	reqData := make(map[string]interface{})
	err3 := c.ShouldBindBodyWith(&reqData, binding.JSON)
	if err3 != nil {
		svc.Logger.Errorf("tracerMiddleware reqData Error: %+v", err3)
		reqData = nil
	}

	if cb, ok := c.Get(gin.BodyBytesKey); ok {
		if cbb, ok := cb.([]byte); ok {
			c.Request.Body = io.NopCloser(bytes.NewBuffer(cbb))
		}
	}
	return reqData
}
