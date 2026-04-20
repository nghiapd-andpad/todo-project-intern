package user

import (
	"context"

	"github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/domain/entity"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/domain/gateway"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/usecase/user/input"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/usecase/user/output"
	"golang.org/x/crypto/bcrypt"
)

type userCreator struct {
	userCommands gateway.UserCommandsGateway
	userQueries  gateway.UserQueriesGateway
}

func NewUserCreator(userCommands gateway.UserCommandsGateway, userQueries gateway.UserQueriesGateway) UserCreator {
	return &userCreator{userCommands: userCommands, userQueries: userQueries}
}

func (u *userCreator) Register(ctx context.Context, in *input.UserRegister) (*output.UserRegister, error) {
	// Check username already exists
	existingUser, _ := u.userQueries.GetByUsername(ctx, in.Username)
	if existingUser != nil {
		return nil, entity.ErrUsernameAlreadyExists
	}

	// Check email already exists
	existingUser, _ = u.userQueries.GetByEmail(ctx, in.Email)
	if existingUser != nil {
		return nil, entity.ErrEmailAlreadyExists
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, entity.ErrInternal
	}

	// Create Entity
	userEnt := &entity.User{
		Username:       in.Username,
		Email:          in.Email,
		HashedPassword: string(hashedPassword),
	}

	// Save to DB
	createdUser, err := u.userCommands.Create(ctx, userEnt)
	if err != nil {
		return nil, err
	}

	// Return output DTO
	return &output.UserRegister{
		User: &output.UserDTO{
			ID:       createdUser.ID.String(),
			Username: createdUser.Username,
			Email:    createdUser.Email,
		},
	}, nil
}
