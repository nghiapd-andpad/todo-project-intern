package user

import (
	userv1 "github.com/nghiapd-andpad/todo-project-intern/proto/user/v1"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/config"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/usecase/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type UserHandler struct {
	userv1.UnimplementedUserServiceServer
	UserCreator       user.UserCreator
	UserAuthenticator user.UserAuthenticator
}

func NewUserHandler(
	creator user.UserCreator,
	authenticator user.UserAuthenticator,
) *UserHandler {
	return &UserHandler{
		UserCreator:       creator,
		UserAuthenticator: authenticator,
	}
}

func NewGRPCServer(cfg *config.Config, handler *UserHandler) *grpc.Server {
	s := grpc.NewServer()
	userv1.RegisterUserServiceServer(s, handler)
	reflection.Register(s)
	return s
}
