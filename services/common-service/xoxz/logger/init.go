package logger

import (
	"fmt"
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Config struct {
	Path                 string   `yaml:"path"`
	File                 string   `yaml:"file"`
	Prefix               string   `yaml:"prefix"`
	Level                string   `yaml:"level"`
	TimeFormat           string   `yaml:"timeFormat"`
	CtxKeys              []string `yaml:"ctxKeys"`
	Header               bool     `yaml:"header"`
	StSkip               int      `yaml:"stSkip"`
	Stdout               bool     `yaml:"stdout"`
	RotateSize           int      `yaml:"rotateSize"`
	RotateExpire         int      `yaml:"rotateExpire"`
	RotateBackupLimit    int      `yaml:"rotateBackupLimit"`
	RotateBackupExpire   int      `yaml:"rotateBackupExpire"`
	RotateBackupCompress int      `yaml:"rotateBackupCompress"`
	RotateCheckInterval  string   `yaml:"rotateCheckInterval"`
	StdoutColorDisabled  bool     `yaml:"stdoutColorDisabled"`
	WriterColorEnable    bool     `yaml:"writerColorEnable"`
	Flags                int      `yaml:"flags"`
}

// InitLogger Zap logger
func InitXoXZLogger(config Config) (*XOXZ, func()) {
	if err := os.MkdirAll(config.Path, 0755); err != nil {
		fmt.Printf("Warning: Failed to create log directory, logging to stdout only: %v\n", err)
		config.File = "" // Disable file logging if path is invalid
	}

	// Use AtomicLevel so the level can be changed at runtime later
	atom := zap.NewAtomicLevel()
	level, err := zapcore.ParseLevel(config.Level)
	if err != nil {
		level = zapcore.InfoLevel
	}
	atom.SetLevel(level)

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	// Add console colors if config allows
	consoleEncoderConfig := encoderConfig
	if !config.StdoutColorDisabled {
		consoleEncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	var cores []zapcore.Core

	// File Core
	if config.File != "" {
		w := zapcore.AddSync(&lumberjack.Logger{
			Filename:   filepath.Join(config.Path, config.File),
			MaxSize:    config.RotateSize,
			MaxBackups: config.RotateBackupLimit,
			MaxAge:     config.RotateExpire,
			Compress:   config.RotateBackupCompress > 0,
		})
		cores = append(cores, zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), w, atom))
	}

	// Stdout Core
	if config.Stdout {
		cores = append(cores, zapcore.NewCore(
			zapcore.NewConsoleEncoder(consoleEncoderConfig),
			zapcore.AddSync(os.Stdout),
			atom,
		))
	}

	core := zapcore.NewTee(cores...)

	// Add useful built-in options
	logger := zap.New(core,
		zap.AddCaller(),
		zap.AddCallerSkip(config.StSkip),
		zap.AddStacktrace(zapcore.ErrorLevel), // Automatically print stacktrace on Error log
	)

	// Return logger and cleanup function
	cleanup := func() {
		_ = logger.Sync()
	}

	return &XOXZ{Logger: logger}, cleanup
}

// InitSugarLogger
// func InitSugarLogger(config LoggerConfig) (*XOXZ_Logger, func()) {
// 	sugar, cleanup := InitXoXZLogger(config)

// 	return &XOXZ_Logger{logger: sugar.logger.Sugar()}, cleanup
// }
