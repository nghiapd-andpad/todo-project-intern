// Package user provides use cases related to user management, such as retrieving user information by ID, username, or email.
package user

import (
	"context"
	"fmt"

	"github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/domain/gateway"
	"github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/usecase/output"
)

type UserGetter struct {
	userGateway gateway.UserGateway
}

func NewUserGetter(userGateway gateway.UserGateway) *UserGetter {
	return &UserGetter{userGateway: userGateway}
}

func (u *UserGetter) GetByID(ctx context.Context, id string) (*output.UserOutput, error) {
	user, err := u.userGateway.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("UserGetter.GetByID: %w", err)
	}

	return &output.UserOutput{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}, nil
}

func (u *UserGetter) GetByIDs(ctx context.Context, ids []string) ([]*output.UserOutput, error) {
	users, err := u.userGateway.GetByIDs(ctx, ids)
	if err != nil {
		return nil, fmt.Errorf("UserGetter.GetByIDs: %w", err)
	}

	// map response
	res := make([]*output.UserOutput, 0, len(users))
	for _, user := range users {
		res = append(res, &output.UserOutput{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
		})
	}

	return res, nil
}

func (u *UserGetter) GetByUsername(ctx context.Context, username string) (*output.UserOutput, error) {
	user, err := u.userGateway.GetByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("UserGetter.GetByUsername: %w", err)
	}

	return &output.UserOutput{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}, nil
}

func (u *UserGetter) GetByEmail(ctx context.Context, email string) (*output.UserOutput, error) {
	user, err := u.userGateway.GetByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("UserGetter.GetByEmail: %w", err)
	}

	return &output.UserOutput{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}, nil
}
