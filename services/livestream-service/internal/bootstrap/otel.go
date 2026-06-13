package bootstrap

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/himbo22/xoxz/common-service/monitoring/telemetry"
	"github.com/himbo22/xoxz/livestream-service/internal/config"
)

func InitOtelCollector(cfg config.TelemetryConfig) (shutdown func(context.Context) error, err error) {
	// 1. PREVENT SERVER HANG: Limit initialization to 10 seconds
	initCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var exportType telemetry.ExporterType
	switch cfg.ExporterType {
	case "http":
		exportType = telemetry.ExporterHTTP
	default:
		exportType = telemetry.ExporterGRPC
	}

	telemetryCfg := telemetry.Config{
		ServiceName:    cfg.ServiceName,
		ServiceVersion: cfg.ServiceVersion,
		Environment:    cfg.Environment,
		ExporterType:   exportType,
		Endpoint:       cfg.Endpoint,
		SampleRate:     cfg.SampleRate,
		Insecure:       cfg.Insecure,
	}

	var tracerShutdown func(context.Context) error
	tracerShutdown, err = telemetry.InitTracer(initCtx, telemetryCfg) // Use initCtx
	if err != nil {
		return nil, fmt.Errorf("failed to init tracer: %w", err)
	}

	var meterShutdown func(context.Context) error
	meterShutdown, err = telemetry.InitMeter(initCtx, telemetryCfg) // Use initCtx
	if err != nil {
		// Use background context here because initCtx may have timed out
		_ = tracerShutdown(context.Background())
		return nil, fmt.Errorf("failed to init meter: %w", err)
	}

	// 2. COLLECT SHUTDOWN ERRORS (Using errors.Join from Go 1.20+)
	shutdown = func(shutdownCtx context.Context) error {
		var errs []error

		if err := meterShutdown(shutdownCtx); err != nil {
			errs = append(errs, fmt.Errorf("meter shutdown error: %w", err))
		}

		if err := tracerShutdown(shutdownCtx); err != nil {
			errs = append(errs, fmt.Errorf("tracer shutdown error: %w", err))
		}

		if len(errs) > 0 {
			return errors.Join(errs...) // Join all errors and report once
		}

		log.Println("OTel resources flushed and shut down successfully.")
		return nil
	}

	log.Printf("Otel init OK: %+v", cfg)

	return shutdown, nil
}
