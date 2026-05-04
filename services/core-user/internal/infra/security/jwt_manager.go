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

func (j *jwtManager) Generate(
	ctx context.Context,
	payload gateway.TokenPayload,
	duration time.Duration,
) (string, error) {
	claims := UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		UserID: int64(payload.UserID),
		Roles:  payload.Roles,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signed, err := token.SignedString(j.secretKey)
	if err != nil {
		return "", fmt.Errorf("sign token: %w", err)
	}

	return signed, nil
}
