package user

import (
	"context"
	"time"

	"github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/domain/entity"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/domain/gateway"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/usecase/user/input"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/usecase/user/output"
	"golang.org/x/crypto/bcrypt"
)

type userAuthenticator struct {
	userQueriesGateway gateway.UserQueriesGateway
	tokenManager       gateway.TokenManager
}

func NewUserAuthenticator(repo gateway.UserQueriesGateway, tokenGenerator gateway.TokenManager) UserAuthenticator {
	return &userAuthenticator{userQueriesGateway: repo, tokenManager: tokenGenerator}
}

func (u *userAuthenticator) Login(ctx context.Context, in *input.UserLogin) (*output.UserLogin, error) {
	// Find user by username
	userEnt, err := u.userQueriesGateway.GetByUsername(ctx, in.Username)
	if err != nil {
		return nil, entity.ErrInvalidCredentials
	}

	// check password
	err = bcrypt.CompareHashAndPassword([]byte(userEnt.HashedPassword), []byte(in.Password))
	if err != nil {
		return nil, entity.ErrInvalidCredentials
	}

	// Create JWT token
	payload := gateway.TokenPayload{
		UserID: userEnt.ID,
		Roles:  []string{"user"},
	}

	token, error := u.tokenManager.Generate(ctx, payload, 24*time.Hour)
	if error != nil {
		return nil, entity.ErrInternal
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

func (u *userAuthenticator) Verify(ctx context.Context, token string) (*output.VerifyTokenOutput, error) {
	// Gọi sang Infra (Manager)
	payload, err := u.tokenManager.Verify(ctx, token)
	if err != nil {
		return nil, err
	}

	return &output.VerifyTokenOutput{
		UserID: payload.UserID.String(),
		Roles:  payload.Roles,
	}, nil
}
