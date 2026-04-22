package user

import (
	"github.com/nghiapd-andpad/todo-project-intern/pkg/auth"
	userv1 "github.com/nghiapd-andpad/todo-project-intern/proto/user/v1"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/config"
	usecase "github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/usecase/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
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

func NewGRPCServer(cfg *config.Config, handler *UserHandler) *grpc.Server {
	s := grpc.NewServer(
		grpc.UnaryInterceptor(auth.UnaryServerInterceptor()),
	)
	userv1.RegisterUserServiceServer(s, handler)
	reflection.Register(s)
	return s
}
