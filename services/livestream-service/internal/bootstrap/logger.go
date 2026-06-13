package bootstrap

import (
	xoxz "github.com/himbo22/xoxz/common-service/xoxz/logger"
	"github.com/himbo22/xoxz/livestream-service/internal/config"
)

// InitLogger logger
func InitLogger(config config.LoggerConfig) (*xoxz.XOXZ, func()) {
	xoxzLoggerConfig := xoxz.Config{
		Path:                 config.Path,
		File:                 config.File,
		Prefix:               config.Prefix,
		Level:                config.Level,
		TimeFormat:           config.TimeFormat,
		CtxKeys:              config.CtxKeys,
		Header:               config.Header,
		StSkip:               config.StSkip,
		Stdout:               config.Stdout,
		RotateSize:           config.RotateSize,
		RotateExpire:         config.RotateExpire,
		RotateBackupLimit:    config.RotateBackupLimit,
		RotateBackupExpire:   config.RotateBackupExpire,
		RotateBackupCompress: config.RotateBackupCompress,
		RotateCheckInterval:  config.RotateCheckInterval,
		StdoutColorDisabled:  config.StdoutColorDisabled,
		WriterColorEnable:    config.WriterColorEnable,
		Flags:                config.Flags,
	}
	return xoxz.InitXoXZLogger(xoxzLoggerConfig)
}
