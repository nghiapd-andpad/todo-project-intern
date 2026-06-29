package context

import (
	"context"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors"
)

type grpcCallMetaKey struct{} // type for grpc call metadata in context

// WithGRPCCallMeta creates a new context with the provided gRPC call metadata.
func WithGRPCCallMeta(
	ctx context.Context,
	callMeta interceptors.CallMeta,
) context.Context {
	return context.WithValue(ctx, &grpcCallMetaKey{}, callMeta)
}

// GRPCCallMetaFromContext retrieves the gRPC call metadata from the context.
// If the metadata is not present, it returns nil.
func GRPCCallMetaFromContext(ctx context.Context) interceptors.CallMeta {
	if v := ctx.Value(&grpcCallMetaKey{}); v != nil {
		callMeta, ok := v.(interceptors.CallMeta)
		if ok {
			return callMeta
		}
	}

	return interceptors.CallMeta{}
}
