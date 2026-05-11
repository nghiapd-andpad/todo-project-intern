package service

import (
	"context"
	"fmt"

	"github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/domain/entity"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/domain/gateway"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/usecase/output"
)

type UserGetter struct {
	userQueriesGateway gateway.UserQueriesGateway
}

func NewUserGetter(userQueriesGateway gateway.UserQueriesGateway) *UserGetter {
	return &UserGetter{userQueriesGateway: userQueriesGateway}
}

func (u *UserGetter) GetByID(ctx context.Context, id entity.UserID) (*output.UserDTO, error) {
	userEnt, err := u.userQueriesGateway.GetByID(ctx, id)

	if err != nil {
		return nil, fmt.Errorf("userGetter.GetByID: %w", err)
	}
	if userEnt == nil {
		return nil, entity.NewNotFound("user not found").
			WithDetail("user_id", id.String())
	}

	return &output.UserDTO{
		ID:       userEnt.ID.String(),
		Username: userEnt.Username,
		Email:    userEnt.Email,
	}, nil
}

func (u *UserGetter) GetByIDs(ctx context.Context, ids []entity.UserID) ([]*output.UserDTO, error) {
	userEntities, err := u.userQueriesGateway.GetByIDs(ctx, ids)
	if err != nil {
		return nil, fmt.Errorf("userGetter.GetByIDs: %w", err)
	}

	dtos := make([]*output.UserDTO, 0, len(userEntities))
	for _, ent := range userEntities {
		if ent == nil {
			continue
		}
		dtos = append(dtos, &output.UserDTO{
			ID:       ent.ID.String(),
			Username: ent.Username,
			Email:    ent.Email,
		})
	}

	return dtos, nil
}

func (u *UserGetter) GetByUsername(ctx context.Context, username string) (*output.UserDTO, error) {
	userEnt, err := u.userQueriesGateway.GetByUsername(ctx, username)

	if err != nil {
		return nil, fmt.Errorf("userGetter.GetByUsername: %w", err)
	}
	if userEnt == nil {
		return nil, entity.NewNotFound("user not found").
			WithDetail("username", username)
	}

	return &output.UserDTO{
		ID:       userEnt.ID.String(),
		Username: userEnt.Username,
		Email:    userEnt.Email,
	}, nil
}

func (u *UserGetter) GetByEmail(ctx context.Context, email string) (*output.UserDTO, error) {
	userEnt, err := u.userQueriesGateway.GetByEmail(ctx, email)

	if err != nil {
		return nil, fmt.Errorf("userGetter.GetByEmail: %w", err)
	}
	if userEnt == nil {
		return nil, entity.NewNotFound("user not found").
			WithDetail("email", email)
	}

	return &output.UserDTO{
		ID:       userEnt.ID.String(),
		Username: userEnt.Username,
		Email:    userEnt.Email,
	}, nil
}
