// Package jwt provides utilities for handling JWT in the BFF service, including token verification and payload extraction.
package jwt

import (
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v5"

	"github.com/nghiapd-andpad/todo-project-intern/services/bff-web/internal/config"
)

type Claims struct {
	jwt.RegisteredClaims
	UserID int64    `json:"user_id"`
	Roles  []string `json:"roles"`
}

type TokenPayload struct {
	UserID string
	Roles  []string
}

type JWTManager struct {
	secretKey []byte
}

func NewJwtManager(cfg *config.Config) *JWTManager {
	return &JWTManager{secretKey: []byte(cfg.JWTSecret)}
}

func (m *JWTManager) Verify(tokenStr string) (*TokenPayload, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(
		tokenStr,
		claims,
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return m.secretKey, nil
		},
	)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, fmt.Errorf("token expired")
		}
		return nil, fmt.Errorf("invalid token")
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return &TokenPayload{
		UserID: fmt.Sprintf("%d", claims.UserID),
		Roles:  claims.Roles,
	}, nil
}
