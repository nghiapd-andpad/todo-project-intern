package service

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/config"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/domain/entity"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/domain/gateway"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/usecase/input"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/usecase/output"
)

type UserAuthenticator struct {
	userQueriesGateway gateway.UserQueriesGateway
	tokenManager       gateway.TokenManager
	cfg                *config.Config
}

func NewUserAuthenticator(
	userQueriesGateway gateway.UserQueriesGateway,
	tokenManager gateway.TokenManager,
	cfg *config.Config,
) *UserAuthenticator {
	return &UserAuthenticator{
		userQueriesGateway: userQueriesGateway,
		tokenManager:       tokenManager,
		cfg:                cfg,
	}
}

func (u *UserAuthenticator) Login(ctx context.Context, in *input.UserLogin) (*output.UserLogin, error) {
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
	duration := time.Duration(u.cfg.JWTExpireHours) * time.Hour

	token, err := u.tokenManager.Generate(ctx, gateway.TokenPayload{
		UserID: userEnt.ID,
		Roles:  []string{"user"},
	}, duration)
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
