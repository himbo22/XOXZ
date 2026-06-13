package telemetry

import (
	"context"
	"errors"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
)

func InitMeter(ctx context.Context, cfg Config) (shutdown func(ctx context.Context) error, err error) {
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

	opts := []otlpmetricgrpc.Option{
		otlpmetricgrpc.WithEndpoint(cfg.Endpoint),
	}
	if cfg.Insecure {
		opts = append(opts, otlpmetricgrpc.WithInsecure())
	}

	exp, err := otlpmetricgrpc.New(ctx, opts...)
	if err != nil {
		return shutdown, err
	}

	res, err := buildResource(ctx, cfg)
	if err != nil {
		return shutdown, err
	}

	mp := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(exp)),
		sdkmetric.WithResource(res),
	)
	shutdownFuncs = append(shutdownFuncs, mp.Shutdown)
	otel.SetMeterProvider(mp)

	return shutdown, nil
}

func GetMeter(name string) metric.Meter {
	return otel.GetMeterProvider().Meter(name)
}
