package user

import (
	"context"
	"fmt"

	"github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/domain/entity"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/domain/gateway"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/usecase/user/input"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/usecase/user/output"
	"golang.org/x/crypto/bcrypt"
)

type UserCreator interface {
	Register(ctx context.Context, in *input.UserRegister) (*output.UserRegister, error)
}

type userCreator struct {
	userCommandsGateway gateway.UserCommandsGateway
	userQueriesGateway  gateway.UserQueriesGateway
}

func NewUserCreator(
	userCommandsGateway gateway.UserCommandsGateway,
	userQueriesGateway gateway.UserQueriesGateway,
) UserCreator {
	return &userCreator{
		userCommandsGateway: userCommandsGateway,
		userQueriesGateway:  userQueriesGateway,
	}
}

func (u *userCreator) Register(ctx context.Context, in *input.UserRegister) (*output.UserRegister, error) {
	// Check username duplicate
	existing, err := u.userQueriesGateway.GetByUsername(ctx, in.Username)
	if err != nil {
		return nil, fmt.Errorf("userCreator.Register check username: %w", err)
	}
	if existing != nil {
		return nil, entity.NewUsernameAlreadyExists()
	}

	// Check email duplicate
	existing, err = u.userQueriesGateway.GetByEmail(ctx, in.Email)
	if err != nil {
		return nil, fmt.Errorf("userCreator.Register check email: %w", err)
	}
	if existing != nil {
		return nil, entity.NewEmailAlreadyExists()
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("userCreator.Register hash password: %w", err)
	}

	// Create entity
	userEnt := &entity.User{
		Username:       in.Username,
		Email:          in.Email,
		HashedPassword: string(hashedPassword),
	}

	created, err := u.userCommandsGateway.Create(ctx, userEnt)
	if err != nil {
		return nil, fmt.Errorf("userCreator.Register: %w", err)
	}

	return &output.UserRegister{
		User: &output.UserDTO{
			ID:       created.ID.String(),
			Username: created.Username,
			Email:    created.Email,
		},
	}, nil
}
