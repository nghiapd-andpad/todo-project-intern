package sentry

import (
	"context"
	"errors"
	"time"

	"github.com/DataDog/dd-trace-go/v2/ddtrace/tracer"
	"github.com/getsentry/sentry-go"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	apperrors "github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/entity"
	contextutil "github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/utils/context"
)

const sentryFlushTimeout = 10 * time.Second
const MetadataForceSentryKey = "ForceSentry"

func InitSentry(
	sentryDsn string,
	sentryEnvironment string,
	logger *zap.Logger,
) (closeFunc func()) {
	if len(sentryDsn) == 0 {
		return func() {}
	}

	if err := sentry.Init(sentry.ClientOptions{
		Dsn:              sentryDsn,
		Environment:      sentryEnvironment,
		AttachStacktrace: true,
	}); err != nil {
		logger.Error("failed to initialize sentry", zap.Error(err))
		return func() {}
	}
	return func() {
		// panicの原因エラーをSentryに送る。復旧はしない。
		causeErr := recover()
		if causeErr != nil {
			sentry.CurrentHub().Recover(causeErr)
			sentry.Flush(sentryFlushTimeout)
			panic(causeErr)
		}
		sentry.Flush(sentryFlushTimeout)
	}
}

func CaptureException(ctx context.Context, err error) {
	if err == nil {
		return
	}

	hub := HubFromContext(ctx)
	hub.WithScope(func(scope *sentry.Scope) {
		setupScopeFromContext(ctx, scope)
		scope.SetLevel(getLevelByError(err))
		hub.CaptureException(err)
	})
}

func CaptureError(ctx context.Context, err error) {
	if ShouldSkipError(err) {
		return
	}

	hub := HubFromContext(ctx)
	hub.WithScope(func(scope *sentry.Scope) {
		setupScopeFromContext(ctx, scope)
		scope.SetLevel(getLevelByError(err))
		hub.CaptureException(err)
	})
}

func HubFromContext(ctx context.Context) *sentry.Hub {
	hub := sentry.GetHubFromContext(ctx)
	if hub == nil {
		hub = sentry.CurrentHub().Clone()
	}
	return hub
}

// DefaultHubToContext injects the default Sentry Hub to the context.
func DefaultHubToContext(ctx context.Context) context.Context {
	hub := sentry.CurrentHub().Clone()
	return sentry.SetHubOnContext(ctx, hub)
}

// RecoverWithContext recovers from a panic in the provided context and
// captures the panic in Sentry.
func RecoverWithContext(ctx context.Context, r any) {
	hub := HubFromContext(ctx)
	hub.WithScope(func(scope *sentry.Scope) {
		setupScopeFromContext(ctx, scope)
		hub.RecoverWithContext(ctx, r)
	})
}

func setupScopeFromContext(ctx context.Context, scope *sentry.Scope) {
	callMeta := contextutil.GRPCCallMetaFromContext(ctx)
	if callMeta.Service != "" || callMeta.Method != "" {
		scope.SetContext("grpc", map[string]any{
			"service": callMeta.Service,
			"method":  callMeta.Method,
		})
	}

	if span, ok := tracer.SpanFromContext(ctx); ok {
		spanCtx := span.Context()
		scope.SetContext("datadog", map[string]any{
			"dd.trace_id": spanCtx.TraceID(),
			"dd.span_id":  spanCtx.SpanID(),
		})
	}
}

func ShouldSkipError(err error) bool {
	// Skip nil errors.
	if err == nil {
		return true
	}

	// Return true for canceled errors to skip them for Sentry reporting.
	if errors.Is(err, context.Canceled) {
		return true
	}
	if status.Code(err) == codes.Canceled {
		return true
	}

	var appErr *apperrors.AppError
	if errors.As(err, &appErr) {
		switch appErr.Code {
		case apperrors.ErrInternal:
			return false
		case apperrors.ErrNotFound,
			apperrors.ErrInvalidParameter,
			apperrors.ErrAuthN,
			apperrors.ErrAuthZ:
			return true
		default:
			return false
		}
	}

	return false
}

func getLevelByError(err error) sentry.Level {
	var vErrors validator.ValidationErrors
	if errors.As(err, &vErrors) {
		return sentry.LevelWarning
	}

	var appError *apperrors.AppError
	ok := errors.As(err, &appError)
	if !ok {
		return sentry.LevelError
	}

	switch appError.Code {
	case apperrors.ErrInternal:
		return sentry.LevelError
	case apperrors.ErrNotFound,
		apperrors.ErrInvalidParameter,
		apperrors.ErrAuthN,
		apperrors.ErrAuthZ:
		return sentry.LevelWarning
	default:
		return sentry.LevelError
	}
}
