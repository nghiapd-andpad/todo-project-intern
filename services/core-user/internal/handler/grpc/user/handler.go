// Package user provides gRPC handlers for user-related operations.
package user

import (
	"context"
	"fmt"

	"buf.build/go/protovalidate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"

	"github.com/nghiapd-andpad/todo-project-intern/pkg/auth"
	userv1 "github.com/nghiapd-andpad/todo-project-intern/proto/user/v1"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/config"
	usecase "github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/usecase/user"
)

type UserHandler struct {
	userv1.UnimplementedUserServiceServer

	userCreator       usecase.UserCreator
	userAuthenticator usecase.UserAuthenticator
	userGetter        usecase.UserGetter
}

func NewUserHandler(
	userCreator usecase.UserCreator,
	userAuthenticator usecase.UserAuthenticator,
	userGetter usecase.UserGetter,
) *UserHandler {
	return &UserHandler{
		userCreator:       userCreator,
		userAuthenticator: userAuthenticator,
		userGetter:        userGetter,
	}
}

func NewGRPCServer(cfg *config.Config, handler *UserHandler) (*grpc.Server, error) {
	validator, err := protovalidate.New()
	if err != nil {
		return nil, fmt.Errorf("create validator: %w", err)
	}

	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			auth.UnaryServerInterceptor(),
			func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
				if msg, ok := req.(proto.Message); ok {
					if err := validator.Validate(msg); err != nil {
						return nil, status.Error(codes.InvalidArgument, err.Error())
					}
				}
				return handler(ctx, req)
			},
		),
	)

	userv1.RegisterUserServiceServer(s, handler)
	reflection.Register(s)
	return s, nil
}

func ProvideGRPCServer(cfg *config.Config, handler *UserHandler) (*grpc.Server, func(), error) {
	s, err := NewGRPCServer(cfg, handler)
	if err != nil {
		return nil, nil, err
	}

	return s, func() {}, nil
}
