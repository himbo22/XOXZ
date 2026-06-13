package telemetry

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// StartSpan starts a new span as a child of the current span in ctx.
// Returns the new context (containing the span) and an end function to defer.
//
// Usage:
//
//	ctx, end := telemetry.StartSpan(ctx, "order-service", "ProcessPayment")
//	defer end(nil) // pass error if any
//
//	result, err := doSomething()
//	defer end(err)
func StartSpan(ctx context.Context, tracerName, spanName string, opts ...trace.SpanStartOption) (context.Context, func(error)) {
	tracer := GetTracer(tracerName)
	newCtx, span := tracer.Start(ctx, spanName, opts...)

	end := func(err error) {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
	}

	return newCtx, end
}

// SpanFromContext returns the current span from context.
// Useful when you want to add attributes to an existing span.
//
// Usage:
//
//	span := telemetry.SpanFromContext(ctx)
//	span.SetAttributes(attribute.String("user.id", userID))
func SpanFromContext(ctx context.Context) trace.Span {
	return trace.SpanFromContext(ctx)
}

// AddAttribute adds a key-value attribute to the current span.
//
// Usage:
//
//	telemetry.AddAttribute(ctx, "order.id", orderID)
//	telemetry.AddAttribute(ctx, "user.id", userID)
func AddAttribute(ctx context.Context, key string, value any) {
	span := trace.SpanFromContext(ctx)
	if !span.IsRecording() {
		return
	}

	switch v := value.(type) {
	case string:
		span.SetAttributes(attribute.String(key, v))
	case int:
		span.SetAttributes(attribute.Int(key, v))
	case int64:
		span.SetAttributes(attribute.Int64(key, v))
	case float64:
		span.SetAttributes(attribute.Float64(key, v))
	case bool:
		span.SetAttributes(attribute.Bool(key, v))
	default:
		span.SetAttributes(attribute.String(key, fmt.Sprintf("%v", v)))
	}
}

// RecordError records an error on the current span without ending it.
// Use this when you want to note an error but continue the span.
//
// Usage:
//
//	if err != nil {
//	    telemetry.RecordError(ctx, err)
//	}
func RecordError(ctx context.Context, err error) {
	if err == nil {
		return
	}
	span := trace.SpanFromContext(ctx)
	span.RecordError(err)
	span.SetStatus(codes.Error, err.Error())
}

// TraceID returns the trace ID of the current span as a string.
// Useful for attaching trace_id to log lines.
//
// Usage:
//
//	logger.Info("processing order", "trace_id", telemetry.TraceID(ctx))
func TraceID(ctx context.Context) string {
	span := trace.SpanFromContext(ctx)
	if !span.SpanContext().IsValid() {
		return ""
	}
	return span.SpanContext().TraceID().String()
}

// SpanID returns the span ID of the current span as a string.
//
// Usage:
//
//	logger.Info("processing order", "span_id", telemetry.SpanID(ctx))
func SpanID(ctx context.Context) string {
	span := trace.SpanFromContext(ctx)
	if !span.SpanContext().IsValid() {
		return ""
	}
	return span.SpanContext().SpanID().String()
}
