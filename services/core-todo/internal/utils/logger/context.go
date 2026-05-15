package logger

import (
	"context"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"go.uber.org/zap"
)

var logger *zap.Logger = zap.NewNop()

func SetLogger(l *zap.Logger) {
	if l == nil {
		logger = zap.NewNop()
		return
	}

	logger = l
}

func FromContext(ctx context.Context) *zap.Logger {
	if ctx == nil {
		return logger
	}

	l := ctxzap.Extract(ctx)
	if l != nil {
		return l
	}

	return logger
}

func ToContext(ctx context.Context, l *zap.Logger) context.Context {
	if l == nil {
		l = logger
	}

	return ctxzap.ToContext(ctx, l)
}
