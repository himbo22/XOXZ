package telemetry

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v5"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
)

func EchoMetricsMiddleware() echo.MiddlewareFunc {
	meter := otel.GetMeterProvider().Meter("http.server")

	requestCount, _ := meter.Int64Counter(
		"http.server.request_count",
		metric.WithDescription("Total number of HTTP requests"),
	)
	requestDuration, _ := meter.Float64Histogram(
		"http.server.request_duration_ms",
		metric.WithDescription("HTTP request duration in milliseconds"),
	)
	activeRequests, _ := meter.Int64UpDownCounter(
		"http.server.active_requests",
		metric.WithDescription("Number of requests currently in flight"),
	)

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			start := time.Now()

			attrs := metric.WithAttributes(
				attribute.String("http.method", c.Request().Method),
				attribute.String("http.route", c.Path()),
			)

			activeRequests.Add(c.Request().Context(), 1, attrs)
			defer activeRequests.Add(c.Request().Context(), -1, attrs)

			err := next(c)

			_, status := echo.ResolveResponseStatus(c.Response(), err)
			attrsWithStatus := metric.WithAttributes(
				attribute.String("http.method", c.Request().Method),
				attribute.String("http.route", c.Path()),
				attribute.Int("http.status_code", status),
			)

			requestCount.Add(c.Request().Context(), 1, attrsWithStatus)
			requestDuration.Record(c.Request().Context(),
				float64(time.Since(start).Milliseconds()),
				attrsWithStatus,
			)

			return err
		}
	}
}

func EchoTracingMiddleware() echo.MiddlewareFunc {
	propagator := otel.GetTextMapPropagator()

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			req := c.Request()

			// Extract incoming trace context from headers
			ctx := propagator.Extract(req.Context(), propagation.HeaderCarrier(req.Header))

			// Start span
			tracer := otel.GetTracerProvider().Tracer(c.Path())
			spanName := req.Method + " " + c.Path()
			ctx, span := tracer.Start(ctx, spanName,
				trace.WithSpanKind(trace.SpanKindServer),
				trace.WithAttributes(
					semconv.HTTPMethod(req.Method),
					semconv.HTTPTarget(req.URL.Path),
					semconv.NetHostName(req.Host),
					attribute.String("http.user_agent", req.UserAgent()),
					attribute.String("http.request_id", c.Response().Header().Get(echo.HeaderXRequestID)),
				),
			)
			defer span.End()

			// Pass new ctx into request
			c.SetRequest(req.WithContext(ctx))

			// Call next handler
			err := next(c)

			// Record status
			_, status := echo.ResolveResponseStatus(c.Response(), err)
			span.SetAttributes(semconv.HTTPStatusCode(status))
			if err != nil {
				span.RecordError(err)
				span.SetStatus(codes.Error, err.Error())
			} else if status >= 500 {
				span.SetStatus(codes.Error, http.StatusText(status))
			}

			return err
		}
	}
}

// HTTPMiddleware is a standard net/http middleware that:
// - Extracts incoming trace context from headers (if called by another traced service)
// - Starts a new span for the incoming request
// - Injects span attributes (method, path, status code)
//
// Usage (standard net/http):
//
//	mux := http.NewServeMux()
//	http.ListenAndServe(":8080", telemetry.HTTPMiddleware("order-service")(mux))
//
// Usage (with any router that accepts http.Handler):
//
//	router.Use(telemetry.HTTPMiddleware("order-service"))
func HTTPMiddleware(serviceName string) func(http.Handler) http.Handler {
	tracer := otel.GetTracerProvider().Tracer(serviceName)
	propagator := otel.GetTextMapPropagator()

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract trace context from incoming headers
			ctx := propagator.Extract(r.Context(), propagation.HeaderCarrier(r.Header))
			// Start span
			spanName := r.Method + " " + r.URL.Path
			ctx, span := tracer.Start(ctx, spanName,
				trace.WithSpanKind(trace.SpanKindServer),
				trace.WithAttributes(
					semconv.HTTPMethod(r.Method),
					semconv.HTTPURL(r.URL.String()),
					semconv.HTTPTarget(r.URL.Path),
					semconv.NetHostName(r.Host),
					attribute.String("http.user_agent", r.UserAgent()),
				),
			)
			defer span.End()

			// Wrap ResponseWriter to capture status code
			rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
			next.ServeHTTP(rw, r.WithContext(ctx))

			// Record status code on span
			span.SetAttributes(semconv.HTTPStatusCode(rw.statusCode))
			if rw.statusCode >= 500 {
				span.SetStatus(codes.Error, http.StatusText(rw.statusCode))
			}
		})
	}
}

// InjectHTTPHeaders injects the current trace context into outgoing HTTP request headers.
// Use this when making HTTP calls to other services.
//
// Usage:
//
//	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
//	telemetry.InjectHTTPHeaders(ctx, req)
//	client.Do(req)
func InjectHTTPHeaders(ctx interface{ Done() <-chan struct{} }, req *http.Request) {
	otel.GetTextMapPropagator().Inject(req.Context(), propagation.HeaderCarrier(req.Header))
}

// responseWriter wraps http.ResponseWriter to capture the status code.
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
