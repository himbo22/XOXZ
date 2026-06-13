package logger

import (
	"context"

	"go.uber.org/zap"
)

// 1. Define Type Alias.
// Using '=' makes Logger.Field and zap.Field the SAME type.
// Zero conversion cost.
type Field = zap.Field

type XOXZ struct {
	Logger *zap.Logger
}

// Logger interface test mock
type XoxzLogger interface {
	// Info logs with key-value pairs. e.g.: Logger.Info("failed", "user_id", 123, "attempt", 3)
	Debug(msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(msg string, fields ...Field)
	Debugf(template string, args ...interface{})
	Infof(template string, args ...interface{})
	Warnf(template string, args ...interface{})
	Errorf(template string, args ...interface{})
	Fatalf(template string, args ...interface{})
	WithContext(ctx context.Context) *XOXZ
	With(fields ...Field) *XOXZ
	WithTrace(ctx context.Context) XoxzLogger
	WithEcho() XoxzLogger
}

// NewxoxzLogger XoxzLogger
func NewxoxzLogger(Logger *zap.Logger) XoxzLogger {
	return &XOXZ{Logger: Logger}
}

// Debug implements [Logger].
func (l *XOXZ) Debug(msg string, fields ...Field) {
	l.Logger.Debug(msg, fields...)
}

// Error implements [Logger].
func (l *XOXZ) Error(msg string, fields ...Field) {
	l.Logger.Error(msg, fields...)
}

// Info implements [Logger].
func (l *XOXZ) Info(msg string, fields ...Field) {
	l.Logger.Info(msg, fields...)
}

// Warn implements [Logger].
func (l *XOXZ) Warn(msg string, fields ...Field) {
	l.Logger.Warn(msg, fields...)
}

// Debugf log debug level format
func (l *XOXZ) Debugf(template string, args ...interface{}) {
	l.Logger.Sugar().Debugf(template, args...)
}

// Infof log info level format
func (l *XOXZ) Infof(template string, args ...interface{}) {
	l.Logger.Sugar().Infof(template, args...)
}

// Warnf log warn level format
func (l *XOXZ) Warnf(template string, args ...interface{}) {
	l.Logger.Sugar().Warnf(template, args...)
}

// Errorf log error level format
func (l *XOXZ) Errorf(template string, args ...interface{}) {
	l.Logger.Sugar().Errorf(template, args...)
}

// Fatalf log fatal level format
func (l *XOXZ) Fatalf(template string, args ...interface{}) {
	l.Logger.Sugar().Fatalf(template, args...)
}

// WithContext context Logger
func (l *XOXZ) WithContext(ctx context.Context) *XOXZ {
	// Get TraceID from middleware
	if traceID, ok := ctx.Value("X-Trace-ID").(string); ok {
		// FIX BUG: Must use zap.String()
		return &XOXZ{
			Logger: l.Logger.With(zap.String("trace_id", traceID)),
		}
	}

	// If no traceID, return as-is
	return l
}

// With fields Logger
func (l *XOXZ) With(args ...Field) *XOXZ {
	return &XOXZ{
		Logger: l.Logger.With(args...),
	}
}

func (l *XOXZ) WithTrace(ctx context.Context) XoxzLogger {
	l.Logger.With(zap.String("trace_id", ctx.Value("X-Trace-ID").(string)))
	return l
}

// WithEcho implements [XoxzLogger].
func (l *XOXZ) WithEcho() XoxzLogger {
	return &XOXZ{
		Logger: l.Logger.WithOptions(
			zap.AddCallerSkip(1),
		),
	}
}

// Helper functions fields
func String(key, val string) zap.Field {
	return zap.String(key, val)
}

func Int(key string, val int) zap.Field {
	return zap.Int(key, val)
}

func Int64(key string, val int64) zap.Field {
	return zap.Int64(key, val)
}

func Float64(key string, val float64) zap.Field {
	return zap.Float64(key, val)
}

func Bool(key string, val bool) zap.Field {
	return zap.Bool(key, val)
}

func Any(key string, val interface{}) zap.Field {
	return zap.Any(key, val)
}

func Error(err error) zap.Field {
	return zap.Error(err)
}

// RequestID field tracing
func RequestID(id string) zap.Field {
	return zap.String("request_id", id)
}

// UserID field user tracking
func UserID(id string) zap.Field {
	return zap.String("user_id", id)
}

// Method field HTTP method
func Method(method string) zap.Field {
	return zap.String("method", method)
}

// Path field HTTP path
func Path(path string) zap.Field {
	return zap.String("path", path)
}

// StatusCode field HTTP status
func StatusCode(code int) zap.Field {
	return zap.Int("status_code", code)
}

// Duration field timing
func Duration(duration interface{}) zap.Field {
	return zap.Any("duration", duration)
}
