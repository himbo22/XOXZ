package util

import (
	"context"

	"github.com/himbo22/xoxz/common-service/monitoring/telemetry"
	xoxz "github.com/himbo22/xoxz/common-service/xoxz/logger"
	"go.uber.org/zap"
)

func TraceFields(ctx context.Context, l xoxz.XoxzLogger) xoxz.XoxzLogger {
	return l.With(
		zap.String("trace_id", telemetry.TraceID(ctx)),
		zap.String("span_id", telemetry.SpanID(ctx)),
	)
}
