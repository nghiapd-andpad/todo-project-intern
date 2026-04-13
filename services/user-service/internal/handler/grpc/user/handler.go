package user

import (
	userv1 "github.com/nghiapd-andpad/todo-project-intern/proto/user/v1"
	"github.com/nghiapd-andpad/todo-project-intern/services/user-service/internal/usecase/user"
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
