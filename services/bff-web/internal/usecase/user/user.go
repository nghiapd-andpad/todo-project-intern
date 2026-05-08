// Package user provides use cases related to user management, such as retrieving user information by ID, username, or email.
package user

import (
	"context"
	"fmt"

	"github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/domain/entity"
	"github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/domain/gateway"
)

type UserGetter struct {
	userGateway gateway.UserGateway
}

func NewUserGetter(userGateway gateway.UserGateway) *UserGetter {
	return &UserGetter{userGateway: userGateway}
}

func (u *UserGetter) GetByID(ctx context.Context, id string) (*entity.User, error) {
	user, err := u.userGateway.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("UserGetter.GetByID: %w", err)
	}

	return user, nil
}

func (u *UserGetter) GetByIDs(ctx context.Context, ids []string) ([]*entity.User, error) {
	users, err := u.userGateway.GetByIDs(ctx, ids)
	if err != nil {
		return nil, fmt.Errorf("UserGetter.GetByIDs: %w", err)
	}

	return users, nil
}

func (u *UserGetter) GetByUsername(ctx context.Context, username string) (*entity.User, error) {
	user, err := u.userGateway.GetByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("UserGetter.GetByUsername: %w", err)
	}

	return user, nil
}

func (u *UserGetter) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	user, err := u.userGateway.GetByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("UserGetter.GetByEmail: %w", err)
	}

	return user, nil
}
