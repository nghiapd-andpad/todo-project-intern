package security

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/config"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/domain/entity"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/domain/gateway"
)

type jwtManager struct {
	secretKey []byte
}

type UserClaims struct {
	jwt.RegisteredClaims
	UserID int64    `json:"user_id"`
	Roles  []string `json:"roles"`
}

func NewJWTManager(cfg *config.Config) gateway.TokenManager {
	return &jwtManager{
		secretKey: []byte(cfg.JWTSecret),
	}
}

func (j *jwtManager) Generate(ctx context.Context, payload gateway.TokenPayload, duration time.Duration) (string, error) {
	claims := UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		UserID: int64(payload.UserID),
		Roles:  payload.Roles,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(j.secretKey))
	if err != nil {
		return "", fmt.Errorf("%w: %v", entity.ErrTokenSigning, err)
	}

	return signedToken, nil
}

func (j *jwtManager) Verify(ctx context.Context, tokenStr string) (*gateway.TokenPayload, error) {
	claims := &UserClaims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, entity.ErrUntrustedMethod
		}
		return []byte(j.secretKey), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, entity.ErrExpiredToken
		}
		return nil, fmt.Errorf("%w: %v", entity.ErrInvalidToken, err)
	}

	if !token.Valid {
		return nil, entity.ErrInvalidToken
	}

	return &gateway.TokenPayload{
		UserID: entity.UserID(claims.UserID),
		Roles:  claims.Roles,
	}, nil
}
