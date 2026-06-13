package telemetry

import (
	"context"
	"errors"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

type ExporterType string

const (
	ExporterGRPC ExporterType = "grpc"
	ExporterHTTP ExporterType = "http"
)

// Config holds the configuration for the telemetry package.
type Config struct {
	ServiceName    string
	ServiceVersion string
	Environment    string // e.g. "production", "staging", "development"
	ExporterType   ExporterType
	Endpoint       string  // OTel Collector endpoint e.g. "localhost:4317" for gRPC, "localhost:4318" for HTTP
	SampleRate     float64 // 1.0 = 100%, 0.1 = 10%
	Insecure       bool    // disable TLS (useful for local development)
}

// DefaultConfig returns a sensible default config for local development.
func DefaultConfig(serviceName string) Config {
	return Config{
		ServiceName:    serviceName,
		ServiceVersion: "1.0.0",
		Environment:    "development",
		ExporterType:   ExporterGRPC,
		Endpoint:       "localhost:4317",
		SampleRate:     1.0,
		Insecure:       true,
	}
}

type ShutdownFunc func()

// InitTracer initializes the OpenTelemetry tracer provider.
// Returns a shutdown function that should be deferred in main().
//
// Usage:
//
//	shutdown, err := telemetry.InitTracer(ctx, telemetry.DefaultConfig("order-service"))
//	if err != nil {
//	    log.Fatal(err)
//	}
//	defer shutdown()
func InitTracer(ctx context.Context, cfg Config) (shutdown func(ctx context.Context) error, err error) {
	// validate config
	var shutdownFuncs []func(context.Context) error

	shutdown = func(ctx context.Context) error {
		var shutdownErr error

		for i := len(shutdownFuncs) - 1; i >= 0; i-- {
			if e := shutdownFuncs[i](ctx); e != nil {
				shutdownErr = errors.Join(shutdownErr, e)
			}
		}

		return shutdownErr
	}

	// 1. Build resource (service metadata attached to every span)
	res, err := buildResource(ctx, cfg)
	if err != nil {
		return shutdown, err
	}

	// 2. Build exporter
	exp, err := buildExporter(ctx, cfg)
	if err != nil {
		return shutdown, err
	}
	shutdownFuncs = append(shutdownFuncs, exp.Shutdown)

	// 3. Build tracer provider
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(buildSampler(cfg.SampleRate)),
	)
	shutdownFuncs = append(shutdownFuncs, tp.Shutdown)

	// 4. Set global tracer provider
	otel.SetTracerProvider(tp)

	// 5. Set global propagator (W3C TraceContext + Baggage)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	return shutdown, nil
}

// GetTracer returns a named tracer from the global provider.
// Use the package or component name as the tracer name.
//
// Usage:
//
//	tracer := telemetry.GetTracer("order-service/payment")
//	ctx, span := tracer.Start(ctx, "ProcessPayment")
//	defer span.End()
func GetTracer(name string) trace.Tracer {
	return otel.GetTracerProvider().Tracer(name)
}

// GetPropagator returns the global text map propagator.
// Use this to inject/extract trace context from HTTP headers or Kafka headers.
//
// Inject (outgoing request):
//
//	telemetry.GetPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))
//
// Extract (incoming request):
//
//	ctx = telemetry.GetPropagator().Extract(ctx, propagation.HeaderCarrier(req.Header))
func GetPropagator() propagation.TextMapPropagator {
	return otel.GetTextMapPropagator()
}
