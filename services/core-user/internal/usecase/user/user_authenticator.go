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
	userRepo       gateway.UserQueriesGateway
	tokenGenerator gateway.TokenGenerator
}

func NewUserAuthenticator(repo gateway.UserQueriesGateway, tokenGenerator gateway.TokenGenerator) UserAuthenticator {
	return &userAuthenticator{userRepo: repo, tokenGenerator: tokenGenerator}
}

func (u *userAuthenticator) Login(ctx context.Context, in *input.UserLogin) (*output.UserLogin, error) {
	// Find user by username
	userEnt, err := u.userRepo.GetByUsername(ctx, in.Username)
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

	token, error := u.tokenGenerator.Generate(ctx, payload, 24*time.Hour)
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
