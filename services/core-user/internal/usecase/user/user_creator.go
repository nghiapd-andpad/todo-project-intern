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

type userCreator struct {
	userRepo gateway.UserCommandsGateway
}

func NewUserCreator(repo gateway.UserCommandsGateway) UserCreator {
	return &userCreator{userRepo: repo}
}

func (u *userCreator) Register(ctx context.Context, in *input.UserRegister) (*output.UserRegister, error) {
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create Entity
	userEnt := &entity.User{
		Username:       in.Username,
		Email:          in.Email,
		HashedPassword: string(hashedPassword),
	}

	// Save to DB
	if err := u.userRepo.Create(ctx, userEnt); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Return output DTO
	return &output.UserRegister{
		User: &output.UserDTO{
			ID:       userEnt.ID.String(),
			Username: userEnt.Username,
			Email:    userEnt.Email,
		},
	}, nil
}
