// Package user provides use cases related to user management, such as retrieving user information by ID, username, or email.
package user

import (
	"context"
	"fmt"

	"github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/domain/entity"
	"github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/domain/gateway"
)

type UserGetter interface {
	GetByID(ctx context.Context, id string) (*entity.User, error)
	GetByIDs(ctx context.Context, ids []string) ([]*entity.User, error)
	GetByUsername(ctx context.Context, username string) (*entity.User, error)
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
}

type userUsecase struct {
	userGateway gateway.UserGateway
}

func NewUserUsecase(userGateway gateway.UserGateway) *userUsecase {
	return &userUsecase{userGateway: userGateway}
}

func (u *userUsecase) GetByID(ctx context.Context, id string) (*entity.User, error) {
	user, err := u.userGateway.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("userUsecase.GetByID: %w", err)
	}
	return user, nil
}

func (u *userUsecase) GetByIDs(ctx context.Context, ids []string) ([]*entity.User, error) {
	users, err := u.userGateway.GetByIDs(ctx, ids)
	if err != nil {
		return nil, fmt.Errorf("userUsecase.GetByIDs: %w", err)
	}
	return users, nil
}

func (u *userUsecase) GetByUsername(ctx context.Context, username string) (*entity.User, error) {
	user, err := u.userGateway.GetByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("userUsecase.GetByUsername: %w", err)
	}
	return user, nil
}

func (u *userUsecase) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	user, err := u.userGateway.GetByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("userUsecase.GetByEmail: %w", err)
	}
	return user, nil
}
