package usecase

import (
	"context"

	"github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/domain/entity"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/usecase/input"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/usecase/output"
)

type UserAuthenticator interface {
	Login(ctx context.Context, in *input.UserLogin) (*output.UserLogin, error)
}

type UserCreator interface {
	Register(ctx context.Context, in *input.UserRegister) (*output.UserRegister, error)
}

type UserGetter interface {
	GetByID(ctx context.Context, id entity.UserID) (*output.UserDTO, error)
	GetByIDs(ctx context.Context, ids []entity.UserID) ([]*output.UserDTO, error)
	GetByUsername(ctx context.Context, username string) (*output.UserDTO, error)
	GetByEmail(ctx context.Context, email string) (*output.UserDTO, error)
}
