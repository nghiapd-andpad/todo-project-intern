package security

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/config"
	"github.com/nghiapd-andpad/todo-project-intern/services/core-user/internal/domain/gateway"
)

type jwtManager struct {
	secretKey []byte
}

func NewJWTManager(cfg *config.Config) gateway.TokenGenerator {
	return &jwtManager{
		secretKey: []byte(cfg.JWTSecret),
	}
}

type userClaims struct {
	jwt.RegisteredClaims
	UserID int64    `json:"user_id"`
	Roles  []string `json:"roles"`
}

func (j *jwtManager) Generate(ctx context.Context, payload gateway.TokenPayload, duration time.Duration) (string, error) {
	claims := userClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		UserID: int64(payload.UserID),
		Roles:  payload.Roles,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	signedToken, err := token.SignedString(j.secretKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return signedToken, nil
}
