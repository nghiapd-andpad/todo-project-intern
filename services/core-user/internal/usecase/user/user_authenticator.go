package user

import (
	"context"
	"fmt"
	"time"

	"github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/domain/entity"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/domain/gateway"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/usecase/user/input"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/usecase/user/output"
	"golang.org/x/crypto/bcrypt"
)

type UserAuthenticator interface {
	Login(ctx context.Context, in *input.UserLogin) (*output.UserLogin, error)
}

type userAuthenticator struct {
	userQueriesGateway gateway.UserQueriesGateway
	tokenManager       gateway.TokenManager
}

func NewUserAuthenticator(
	userQueriesGateway gateway.UserQueriesGateway,
	tokenManager gateway.TokenManager,
) UserAuthenticator {
	return &userAuthenticator{
		userQueriesGateway: userQueriesGateway,
		tokenManager:       tokenManager,
	}
}

func (u *userAuthenticator) Login(ctx context.Context, in *input.UserLogin) (*output.UserLogin, error) {
	// Get user by username
	userEnt, err := u.userQueriesGateway.GetByUsername(ctx, in.Username)
	if err != nil {
		return nil, fmt.Errorf("userAuthenticator.Login: %w", err)
	}

	if userEnt == nil {
		return nil, entity.NewInvalidCredentials()
	}

	if err := bcrypt.CompareHashAndPassword([]byte(userEnt.HashedPassword), []byte(in.Password)); err != nil {
		return nil, entity.NewInvalidCredentials()
	}

	// Generate JWT
	token, err := u.tokenManager.Generate(ctx, gateway.TokenPayload{
		UserID: userEnt.ID,
		Roles:  []string{"user"},
	}, 24*time.Hour)
	if err != nil {
		return nil, fmt.Errorf("userAuthenticator.Login generate token: %w", err)
	}

	return &output.UserLogin{
		AccessToken: token,
		User: &output.UserDTO{
			ID:       userEnt.ID.String(),
			Username: userEnt.Username,
			Email:    userEnt.Email,
		},
	}, nil
}
