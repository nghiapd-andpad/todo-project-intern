// Package context provides helper functions for storing and retrieving request-scoped metadata from context.Context.
package context

import (
	"context"

	"github.com/nghiapd-andpad/todo-project-intern/services/core-todo/internal/domain/entity"
)

type serviceNameKey struct{}

func WithServiceName(ctx context.Context, serviceName string) context.Context {
	return context.WithValue(ctx, serviceNameKey{}, serviceName)
}

func ExtractServiceName(ctx context.Context) (string, error) {
	if v := ctx.Value(serviceNameKey{}); v != nil {
		s, ok := v.(string)
		if ok {
			return s, nil
		}
	}
	return "", entity.NewInternal("ExtractServiceName: failed to extract ServiceName")
}

type methodNameKey struct{}

func WithMethodName(ctx context.Context, methodName string) context.Context {
	return context.WithValue(ctx, methodNameKey{}, methodName)
}

func ExtractMethodName(ctx context.Context) (string, error) {
	if v := ctx.Value(methodNameKey{}); v != nil {
		s, ok := v.(string)
		if ok {
			return s, nil
		}
	}
	return "", entity.NewInternal("ExtractMethodName: failed to method Name")
}
