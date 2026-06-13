package bootstrap

import (
	"context"
	"fmt"
	"log"

	"github.com/himbo22/xoxz/artist-service/internal/config"
	"github.com/himbo22/xoxz/common-service/monitoring/telemetry"
)

func InitOtelCollector(cfg config.TelemetryConfig) (func(), error) {
	ctx := context.Background()

	var exportType telemetry.ExporterType
	switch cfg.ExporterType {
	case "http":
		exportType = telemetry.ExporterHTTP
	default:
		exportType = telemetry.ExporterGRPC
	}

	shutdown, err := telemetry.InitTracer(ctx, telemetry.Config{
		ServiceName:    cfg.ServiceName,
		ServiceVersion: cfg.ServiceVersion,
		Environment:    cfg.Environment,
		ExporterType:   exportType,
		Endpoint:       cfg.Endpoint,
		SampleRate:     cfg.SampleRate,
		Insecure:       cfg.Insecure,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to otel-collector collector: %w", err)
	}

	log.Printf("Otel init OK: %v", cfg)
	return shutdown, nil
}
