package interceptor

import (
	"context"
	"strings"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	logutil "github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/utils/logger"
)

// LoggingUnaryServerInterceptor injects zap logger into each request context
func LoggingUnaryServerInterceptor(baseLogger *zap.Logger) grpc.UnaryServerInterceptor {
	if baseLogger == nil {
		baseLogger = zap.NewNop()
	}

	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		start := time.Now()

		serviceName, methodName := parseFullMethod(info.FullMethod)

		l := baseLogger.With(
			zap.String("grpc.full_method", info.FullMethod),
			zap.String("grpc.service", serviceName),
			zap.String("grpc.method", methodName),
		)

		ctx = logutil.ToContext(ctx, l)

		logutil.Info(ctx, "grpc request started")

		resp, err := handler(ctx, req)

		duration := time.Since(start)

		if err != nil {
			logutil.Error(ctx, "grpc request failed",
				zap.Duration("duration", duration),
				zap.Error(err),
			)
			return resp, err
		}

		logutil.Info(ctx, "grpc request completed",
			zap.Duration("duration", duration),
		)

		return resp, nil
	}
}

// parseFullMethod splits "/package.ServiceName/MethodName" into (service, method).
func parseFullMethod(fullMethod string) (service string, method string) {
	fullMethod = strings.TrimPrefix(fullMethod, "/")

	parts := strings.Split(fullMethod, "/")
	if len(parts) != 2 {
		return "unknown", fullMethod
	}

	return parts[0], parts[1]
}
