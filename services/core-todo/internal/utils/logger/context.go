// Package logger provides Zap logger initialization and context-aware logging helpers.
package logger

import (
	"context"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"go.uber.org/zap"
)

var baseLogger *zap.Logger = zap.NewNop()

func SetLogger(l *zap.Logger) {
	if l == nil {
		baseLogger = zap.NewNop()
		return
	}

	baseLogger = l
}

func FromContext(ctx context.Context) *zap.Logger {
	if ctx == nil {
		return baseLogger
	}

	if l := ctxzap.Extract(ctx); l != nil {
		return l
	}

	return baseLogger
}

func ToContext(ctx context.Context, l *zap.Logger) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}

	if l == nil {
		l = baseLogger
	}

	return ctxzap.ToContext(ctx, l)
}
