package util

import (
	"context"

	"go.uber.org/zap"
)

// Logger interface test mock
type Logger interface {
	Debug(msg string, fields ...zap.Field)
	Info(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
	Fatal(msg string, fields ...zap.Field)
	Debugf(template string, args ...interface{})
	Infof(template string, args ...interface{})
	Warnf(template string, args ...interface{})
	Errorf(template string, args ...interface{})
	Fatalf(template string, args ...interface{})
	WithContext(ctx context.Context) Logger
	With(fields ...zap.Field) Logger
}

// ZapLogger wrapper zap.Logger
type ZapLogger struct {
	logger *zap.Logger
}

// NewZapLogger ZapLogger
func NewZapLogger(logger *zap.Logger) Logger {
	return &ZapLogger{logger: logger}
}

// Debug log debug level
func (l *ZapLogger) Debug(msg string, fields ...zap.Field) {
	l.logger.Debug(msg, fields...)
}

// Info log info level
func (l *ZapLogger) Info(msg string, fields ...zap.Field) {
	l.logger.Info(msg, fields...)
}

// Warn log warn level
func (l *ZapLogger) Warn(msg string, fields ...zap.Field) {
	l.logger.Warn(msg, fields...)
}

// Error log error level
func (l *ZapLogger) Error(msg string, fields ...zap.Field) {
	l.logger.Error(msg, fields...)
}

// Fatal log fatal level
func (l *ZapLogger) Fatal(msg string, fields ...zap.Field) {
	l.logger.Fatal(msg, fields...)
}

// Debugf log debug level format
func (l *ZapLogger) Debugf(template string, args ...interface{}) {
	l.logger.Sugar().Debugf(template, args...)
}

// Infof log info level format
func (l *ZapLogger) Infof(template string, args ...interface{}) {
	l.logger.Sugar().Infof(template, args...)
}

// Warnf log warn level format
func (l *ZapLogger) Warnf(template string, args ...interface{}) {
	l.logger.Sugar().Warnf(template, args...)
}

// Errorf log error level format
func (l *ZapLogger) Errorf(template string, args ...interface{}) {
	l.logger.Sugar().Errorf(template, args...)
}

// Fatalf log fatal level format
func (l *ZapLogger) Fatalf(template string, args ...interface{}) {
	l.logger.Sugar().Fatalf(template, args...)
}

// WithContext context logger
func (l *ZapLogger) WithContext(ctx context.Context) Logger {
	return l
}

// With fields logger
func (l *ZapLogger) With(fields ...zap.Field) Logger {
	return &ZapLogger{logger: l.logger.With(fields...)}
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
