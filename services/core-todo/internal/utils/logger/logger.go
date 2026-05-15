package logger

import (
	"context"
	"fmt"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"go.uber.org/zap"

	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/utils/sentry"
)

func CaptureError(ctx context.Context, err error, message string) {
	if err == nil {
		return
	}

	wErr := fmt.Errorf("%s: %w", message, err)

	sentry.CaptureException(ctx, wErr)

	FromContext(ctx).Error(message, zap.Error(err))
}

func Debug(ctx context.Context, msg string, fields ...zap.Field) {
	FromContext(ctx).Debug(msg, fields...)
}

func Info(ctx context.Context, msg string, fields ...zap.Field) {
	FromContext(ctx).Info(msg, fields...)
}

func Warn(ctx context.Context, msg string, fields ...zap.Field) {
	FromContext(ctx).Warn(msg, fields...)
}

func Error(ctx context.Context, msg string, fields ...zap.Field) {
	FromContext(ctx).Error(msg, fields...)
}

func WithFields(ctx context.Context, fields ...zap.Field) context.Context {
	l := FromContext(ctx).With(fields...)
	return ctxzap.ToContext(ctx, l)
}
