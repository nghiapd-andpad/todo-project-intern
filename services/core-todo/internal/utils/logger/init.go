package logger

import (
	"fmt"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/config"
)

func New(cfg *config.Config) (*zap.Logger, func(), error) {
	if cfg == nil {
		return nil, nil, fmt.Errorf("config is nil")
	}

	level := zapcore.InfoLevel

	switch strings.ToLower(cfg.LogLevel) {
	case "debug":
		level = zapcore.DebugLevel
	case "info":
		level = zapcore.InfoLevel
	case "warn":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	default:
		return nil, nil, fmt.Errorf("invalid log level: %s", cfg.LogLevel)
	}

	var zapCfg zap.Config

	switch strings.ToLower(cfg.LogFormat) {
	case "json":
		zapCfg = zap.NewProductionConfig()
	case "console":
		zapCfg = zap.NewDevelopmentConfig()
	default:
		return nil, nil, fmt.Errorf("invalid log format: %s", cfg.LogFormat)
	}

	zapCfg.Level = zap.NewAtomicLevelAt(level)

	zapLogger, err := zapCfg.Build()
	if err != nil {
		return nil, nil, fmt.Errorf("build zap logger: %w", err)
	}

	SetLogger(zapLogger)

	cleanup := func() {
		_ = zapLogger.Sync()
	}

	return zapLogger, cleanup, nil
}
